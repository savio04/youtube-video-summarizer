package main

import (
	"github.com/savio04/youtube-video-summarizer/internal/logger"
)

func main() {
	// Logger
	logger.InitAppLogger()

	defer logger.AppLogger.Sync()

	// Http server
	app := &application{}

	app.startHttpServer()
}
