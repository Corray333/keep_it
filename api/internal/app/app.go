package app

import (
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/internal/domains/user"
	"github.com/Corray333/keep_it/internal/global_storage"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

type App struct {
	Server *http.Server
	Store  global_storage.Storage
}

type Controller interface {
	Init(router *chi.Mux, store global_storage.Storage)
}

func New() *App {
	router := chi.NewMux()
	// TODO: add timeouts
	server := http.Server{
		Addr:    "0.0.0.0:" + viper.GetString("port"),
		Handler: router,
	}

	app := &App{
		Server: &server,
		Store:  global_storage.New(),
	}

	app.Use(user.NewController())

	return app
}

func (app *App) Use(c Controller) {
	c.Init(app.Server.Handler.(*chi.Mux), app.Store)
}

func (app *App) Run() {
	slog.Info("Server is starting...")
	slog.Info(app.Server.ListenAndServe().Error())
}
