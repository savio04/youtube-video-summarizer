package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/savio04/youtube-video-summarizer/internal/controllers"
	"github.com/savio04/youtube-video-summarizer/internal/env"
	"github.com/savio04/youtube-video-summarizer/internal/logger"
	"go.uber.org/zap"
)

type application struct{}

func (app *application) startHttpServer() http.Handler {
	server := chi.NewRouter()

	server.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	server.Use(middleware.RequestID)
	server.Use(middleware.RealIP)
	server.Use(middleware.Logger)
	server.Use(middleware.Recoverer)
	server.Use(middleware.AllowContentType("application/json"))

	server.Route("/v1", app.routesV1)

	port := env.GetEnvOrDie("HTTP_PORT")

	logger.AppLogger.Info("Starting server on port " + port + "...")

	if err := http.ListenAndServe(":"+port, server); err != nil {
		logger.AppLogger.Fatal("Failed to start server", zap.Error(err))
	}

	return server
}

func (app *application) routesV1(r chi.Router) {
	healthController := controllers.NewHealthController()
	r.Get("/health", healthController.Handler)

	createVideoController := controllers.NewCreateVideoController()
	r.Post("/videos", createVideoController.Handler)

	getVideoController := controllers.NewGetVideoController()
	r.Get("/videos/{videoId}", getVideoController.Handler)
}
