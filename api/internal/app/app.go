package app

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Corray333/keep_it/internal/config"
	"github.com/Corray333/keep_it/internal/domains/category"
	"github.com/Corray333/keep_it/internal/domains/note"
	"github.com/Corray333/keep_it/internal/domains/user"
	"github.com/Corray333/keep_it/internal/storage"
	"github.com/spf13/viper"
)

type controller interface {
	Build()
	Run()
}

type App struct {
	server      *http.Server
	controllers []controller
}

func (app *App) AddController(c controller) {
	app.controllers = append(app.controllers, c)
}

func New() *App {
	config.MustInit(os.Args[1])

	app := &App{}

	router := newRouter()

	// TODO: add timeouts
	server := &http.Server{
		Addr:    "0.0.0.0:" + viper.GetString("server.port"),
		Handler: router,
	}

	app.server = server

	store, err := storage.New()
	if err != nil {
		panic(err)
	}

	// fileManager := files.New()

	userController := user.NewUserController(router, store)
	app.AddController(userController)

	noteController := note.NewNoteController(router, store, userController.GetService())
	app.AddController(noteController)

	categoryController := category.NewCategoryController(router, store, userController.GetService())
	app.AddController(categoryController)

	return app
}

func (app *App) Init() *App {
	for _, c := range app.controllers {
		c.Build()
	}
	return app
}

func (app *App) Run() {
	slog.Info("Server started at " + app.server.Addr)
	for _, c := range app.controllers {
		go c.Run()
	}
	slog.Error(app.server.ListenAndServe().Error())
}
