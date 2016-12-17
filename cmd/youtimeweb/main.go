package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"github.com/kardianos/osext"
	"github.com/spf13/viper"
)

type youtimeConfig struct {
	Port   string
	Config map[string]string
}

// App in main app
type App struct {
	router  *Router
	gp      globalPresenter
	logr    appLogger
	youtime youtimeConfig
}

// globalPresenter contains the fields neccessary for presenting in all templates
type globalPresenter struct {
	SiteName    string
	Description string
	SiteURL     string
}

// TODO localPresenter if we have using template
func SetupApp(r *Router, logger appLogger, templateDirectoryPath string) *App {
	gp := globalPresenter{
		SiteName:    "youtime",
		Description: "Api for native app",
		SiteURL:     "api.floatingcube.com",
	}
	youtimeConfig := youtimeConfig{
		Port:   viper.GetString("port"),
		Config: viper.GetStringMapString("youtime"),
	}
	fmt.Println(viper.GetString("port"))
	return &App{
		router:  r,
		gp:      gp,
		logr:    logger,
		youtime: youtimeConfig,
	}
}

func main() {
	pwd, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatalf("cannot retrieve present working directory: %i", 0600, nil)
	}
	err = LoadConfiguration(pwd)
	if err != nil && os.Getenv("PORT") == "" {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	r := NewRouter()
	logr := newLogger()
	a := SetupApp(r, logr, "")
	port := os.Getenv("PORT")
	if port == "" {
		port = a.youtime.Port
		fmt.Println("port from file", port)
	}

	common := alice.New(context.ClearHandler, a.loggingHandler, a.recoverHandler)
	r.Get("/video/youtube/:id", common.Then(a.Wrap(a.GetYoutubeHandler())))
	r.Post("/video/youtube/:id", common.Then(a.Wrap(a.PostYoutubeHandler())))
	fmt.Println("port from var config", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Errorf("error on serve server %s", err)
	}
}

func LoadConfiguration(pwd string) error {
	viper.SetConfigName("youtime-config")
	viper.AddConfigPath(pwd)
	devPath := pwd[:len(pwd)-3] + "src/youtime/cmd/youtimeweb/"
	_, file, _, _ := runtime.Caller(1)
	configPath := path.Dir(file)
	viper.AddConfigPath(devPath)
	viper.AddConfigPath(configPath)
	return viper.ReadInConfig() // Find and read the config file
}
