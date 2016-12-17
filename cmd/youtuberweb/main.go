package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/kardianos/osext"
	"github.com/spf13/viper"
)

const isDevelopment = true

type youtuberConfig struct {
	expireTime   int
	cookieSecret string
	isProduction bool
	config       map[string]string
}

// App in main app
type App struct {
	router   *Router
	store    *sessions.CookieStore
	gp       globalPresenter
	logr     appLogger
	youtuber youtuberConfig
}

func (a *App) GetStore() *sessions.CookieStore {
	return a.store
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
		SiteName:    "youtuber",
		Description: "Api for native app",
		SiteURL:     "api.floatingcube.com",
	}
	youtuberConfig := youtuberConfig{
		config: viper.GetStringMapString("youtuber"),
	}
	return &App{
		router:   r,
		gp:       gp,
		store:    sessions.NewCookieStore([]byte(youtuberConfig.cookieSecret)),
		logr:     logger,
		youtuber: youtuberConfig,
	}
}

func main() {
	pwd, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatalf("cannot retrieve present working directory: %i", 0600, nil)
	}

	err = LoadConfiguration(pwd)
	if err != nil && viper.GetBool("isProduction") {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	//TODO config static file path and template file path

	r := NewRouter()
	logr := newLogger()
	a := SetupApp(r, logr, "")

	common := alice.New(context.ClearHandler, a.loggingHandler, a.recoverHandler)
	r.Get("/video/youtube/:id", common.Then(a.Wrap(a.GetYoutubeHandler())))
	r.Post("/video/youtube/:id", common.Then(a.Wrap(a.PostYoutubeHandler())))

	//r.ServeFiles("/static/*filepath", http.Dir(staticFilePath))

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Errorf("error on serve server %s", err)
	}
}

func LoadConfiguration(pwd string) error {
	viper.SetConfigName("youtuber-config")
	viper.AddConfigPath(pwd)
	devPath := pwd[:len(pwd)-3] + "/src/youtuber/cmd/youtuberweb/"
	_, file, _, _ := runtime.Caller(1)
	configPath := path.Dir(file)
	viper.AddConfigPath(devPath)
	viper.AddConfigPath(configPath)
	viper.SetDefault("path", devPath)
	return viper.ReadInConfig() // Find and read the config file
}
