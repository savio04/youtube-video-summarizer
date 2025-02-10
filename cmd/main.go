package main

import (
	"github.com/savio04/youtube-video-summarizer/internal/database"
	"github.com/savio04/youtube-video-summarizer/internal/env"
	"github.com/savio04/youtube-video-summarizer/internal/logger"
	"github.com/savio04/youtube-video-summarizer/internal/queue"
	"go.uber.org/zap"
)

func main() {
	// Logger
	logger.InitAppLogger()

	defer logger.AppLogger.Sync()

	// Load envs
	if err := env.LoadEnvs(); err != nil {
		logger.AppLogger.Fatal("Coudn't load .env file")
	}

	// Database connection
	if err := database.Init(); err != nil {
		logger.AppLogger.Fatal("Failed to connect postgres", zap.Error(err))
	}

	// Queue
	if err := queue.Init(); err != nil {
		logger.AppLogger.Fatal("Failed to connect redis", zap.Error(err))
	}

	// Http server
	app := &application{}

	app.startHttpServer()
}
