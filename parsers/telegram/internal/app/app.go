package app

import (
	"os"

	"github.com/Corray333/keep_it/parsers/telegram/internal/config"
	"github.com/Corray333/keep_it/parsers/telegram/internal/external"
	"github.com/Corray333/keep_it/parsers/telegram/internal/file"
	"github.com/Corray333/keep_it/parsers/telegram/internal/repository"
	"github.com/Corray333/keep_it/parsers/telegram/internal/service"
	"github.com/Corray333/keep_it/parsers/telegram/internal/telegram"
	"github.com/Corray333/keep_it/parsers/telegram/internal/transport"
)

type App struct {
	store     *repository.Storage
	service   *service.Service
	transport *transport.Transport
}

func New() *App {
	config.MustInit(os.Args[1])

	storage := repository.New()
	tgClient := telegram.New()
	external := external.New(tgClient)
	fileManager := file.NewFileManager()
	service := service.New(storage, external, fileManager)

	transport := transport.New(service, tgClient)

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
