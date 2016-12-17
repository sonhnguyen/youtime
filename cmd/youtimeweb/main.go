package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"youtime"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"github.com/kardianos/osext"
	"github.com/spf13/viper"
)

type youtimeConfig struct {
	Port          string
	URI           string
	Dbname        string
	Collection    string
	IsDevelopment string
}

// App in main app
type App struct {
	router  *Router
	gp      globalPresenter
	logr    appLogger
	mongodb youtime.Mongodb
	config  youtimeConfig
}

// globalPresenter contains the fields neccessary for presenting in all templates
type globalPresenter struct {
	SiteName    string
	Description string
	SiteURL     string
}

// TODO localPresenter if we have using template
func SetupApp(r *Router, logger appLogger, templateDirectoryPath string) *App {
	var config youtimeConfig
	if viper.GetBool("isDevelopment") {
		config = youtimeConfig{
			IsDevelopment: viper.GetString("isDevelopment"),
			Port:          viper.GetString("port"),
			URI:           viper.GetString("uri"),
			Dbname:        viper.GetString("dbname"),
			Collection:    viper.GetString("collection"),
		}
	} else {
		config = youtimeConfig{
			IsDevelopment: os.Getenv("isDevelopment"),
			Port:          os.Getenv("port"),
			URI:           os.Getenv("uri"),
			Dbname:        os.Getenv("dbname"),
			Collection:    os.Getenv("collection"),
		}
	}

	mongo := youtime.Mongodb{URI: config.URI, Dbname: config.Dbname, Collection: config.Collection}

	gp := globalPresenter{
		SiteName:    "youtime",
		Description: "Api for native app",
		SiteURL:     "api.floatingcube.com",
	}

	return &App{
		router:  r,
		gp:      gp,
		logr:    logger,
		config:  config,
		mongodb: mongo,
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

	common := alice.New(context.ClearHandler, a.loggingHandler, a.recoverHandler)
	r.Get("/video/link", common.Then(a.Wrap(a.GetVideoByLinkHandler())))
	r.Get("/video/id/:id", common.Then(a.Wrap(a.GetVideoByIdHandler())))
	r.Post("/video/:id", common.Then(a.Wrap(a.PostCommentByIdHandler())))

	err = http.ListenAndServe(":"+a.config.Port, r)
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
