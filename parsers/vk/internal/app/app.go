package app

import (
	"os"

	"github.com/Corray333/keep_it/parsers/vk/internal/config"
	"github.com/Corray333/keep_it/parsers/vk/internal/external"
	"github.com/Corray333/keep_it/parsers/vk/internal/file"
	"github.com/Corray333/keep_it/parsers/vk/internal/repository"
	"github.com/Corray333/keep_it/parsers/vk/internal/service"
	"github.com/Corray333/keep_it/parsers/vk/internal/transport"
	"github.com/Corray333/keep_it/parsers/vk/internal/vk"
)

type App struct {
	store     *repository.Storage
	service   *service.Service
	transport *transport.Transport
}

func New() *App {
	config.MustInit(os.Args[1])

	storage := repository.New()
	vkClient := vk.New()
	external := external.New(vkClient)
	fileManager := file.New()
	service := service.New(storage, external, fileManager)

	transport := transport.New(service, vkClient)

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
