package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/savio04/youtube-video-summarizer/internal/controllers"
	"github.com/savio04/youtube-video-summarizer/internal/database"
	"go.uber.org/zap"
)

type application struct {
	logger *zap.SugaredLogger
}

func (app *application) startHttpServer() http.Handler {
	server := chi.NewRouter()
	server.Use(middleware.RequestID)
	server.Use(middleware.RealIP)
	server.Use(middleware.Logger)
	server.Use(middleware.Recoverer)

	err := database.Init()
	if err != nil {
		app.logger.Fatal("Failed to connect postgres", zap.Error(err))
	}

	server.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
			data := map[string]string{
				"status":  "ok",
				"version": "v1",
			}

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)

			json.NewEncoder(writer).Encode(data)
		})

		createVideoController := controllers.NewCreateVideoController()
		r.Post("/videos", createVideoController.Handler)
	})

	app.logger.Info("Starting server on port 8080...")

	if err := http.ListenAndServe(":8080", server); err != nil {
		app.logger.Fatal("Failed to start server", zap.Error(err))
	}

	return server
}
