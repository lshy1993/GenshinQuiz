package main

import (
	"genshin-quiz/config"
	"genshin-quiz/internal/worker"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	app := config.NewApp()
	w := worker.New(worker.Dependencies{
		Config: app.Config,
		Logger: app.Logger,
		DB:     app.DB,
	})
	defer w.Queue().Client.Close()
}
