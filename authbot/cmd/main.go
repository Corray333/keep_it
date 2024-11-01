package main

import (
	"os"

	"github.com/corray333/keep_it_authbot/internal/app"
	"github.com/corray333/keep_it_authbot/internal/config"
)

func main() {
	config.MustInit(os.Args[1])
	app.New().Run()
}
