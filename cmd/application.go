package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/savio04/youtube-video-summarizer/internal/controllers"
	"github.com/savio04/youtube-video-summarizer/internal/database"
	"github.com/savio04/youtube-video-summarizer/internal/env"
	"go.uber.org/zap"
)

type application struct {
	logger *zap.SugaredLogger
}

func (app *application) startHttpServer() http.Handler {
	if err := env.LoadEnvs(); err != nil {
		app.logger.Fatal("Coudn't load .env file")
	}

	server := chi.NewRouter()
	server.Use(middleware.RequestID)
	server.Use(middleware.RealIP)
	server.Use(middleware.Logger)
	server.Use(middleware.Recoverer)
	server.Use(middleware.AllowContentType("application/json"))

	if err := database.Init(); err != nil {
		app.logger.Fatal("Failed to connect postgres", zap.Error(err))
	}

	server.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
			data := map[string]string{
				"status":  "ok",
				"version": "v1",
			}

			writer.WriteHeader(http.StatusOK)

			json.NewEncoder(writer).Encode(data)
		})

		createVideoController := controllers.NewCreateVideoController()
		r.Post("/videos", createVideoController.Handler)
	})

	port := env.GetEnvOrDie("HTTP_PORT")

	app.logger.Info("Starting server on port " + port + "...")

	if err := http.ListenAndServe(":"+port, server); err != nil {
		app.logger.Fatal("Failed to start server", zap.Error(err))
	}

	return server
}
