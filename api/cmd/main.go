package main

import (
	"os"

	"github.com/Corray333/keep_it/internal/app"
	"github.com/Corray333/keep_it/internal/config"
)

func main() {
	config.MustInit(os.Args[1])
	app.New().Run()
}
