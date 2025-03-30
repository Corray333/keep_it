package app

import (
	"os"

	"github.com/Corray333/keep_it/parsers/web/internal/config"
	"github.com/Corray333/keep_it/parsers/web/internal/file"
	"github.com/Corray333/keep_it/parsers/web/internal/repository"
	"github.com/Corray333/keep_it/parsers/web/internal/service"
	"github.com/Corray333/keep_it/parsers/web/internal/transport"
)

type App struct {
	store     *repository.Storage
	service   *service.Service
	transport *transport.Transport
}

func New() *App {
	config.MustInit(os.Args[1])

	storage := repository.New()

	fileManager := file.NewFileManager()
	service := service.New(storage, fileManager)
	transport := transport.New(service)

	app := &App{
		store:     storage,
		service:   service,
		transport: transport,
	}

	return app
}

func (app *App) Run() {
	app.transport.Run()
}
