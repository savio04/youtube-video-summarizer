package controllers

import (
	"encoding/json"
	"net/http"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (healthController *HealthController) Handler(writer http.ResponseWriter, _ *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": "v1",
	}

	writer.WriteHeader(http.StatusOK)

	json.NewEncoder(writer).Encode(data)
}
