package main

import (
	"os"
	"os/signal"
	"syscall"

	"genshin-quiz/config"
	"genshin-quiz/internal/cron"
)

func main() {
	app := config.NewApp()
	// Initialize cron scheduler
	scheduler := cron.NewScheduler(app)

	// Start the scheduler
	scheduler.Start()

	scheduler.logger.Info("Cronjob service started")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	scheduler.logger.Info("Cronjob service shutting down...")

	// Stop the scheduler
	scheduler.Stop()
}
