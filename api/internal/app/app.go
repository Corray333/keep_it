package app

import (
	"log/slog"
	"net/http"

	_ "github.com/Corray333/keep_it/docs"
	"github.com/Corray333/keep_it/internal/domains/user"
	"github.com/Corray333/keep_it/internal/global_storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
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
	router.Use(middleware.Logger)

	// TODO: get allowed origins, headers and methods from cfg
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Set-Cookie", "Refresh", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Максимальное время кеширования предзапроса (в секундах)
	}))

	router.Get("/swagger/*", httpSwagger.WrapHandler)

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
