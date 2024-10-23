package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type application struct {
	logger *zap.SugaredLogger
}

func (app *application) newApp() http.Handler {
	server := chi.NewRouter()
	server.Use(middleware.RequestID)
	server.Use(middleware.RealIP)
	server.Use(middleware.Logger)
	server.Use(middleware.Recoverer)

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
	})

	app.logger.Info("Starting server on port 8080...")

	if err := http.ListenAndServe(":8080", server); err != nil {
		app.logger.Fatal("Failed to start server", zap.Error(err))
	}

	return server
}
