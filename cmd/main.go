package main

import (
	"go.uber.org/zap"
)

func main() {
	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Http server
	app := &application{
		logger: logger,
	}

	app.startHttpServer()
}
