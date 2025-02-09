package controllers

import (
	"encoding/json"
	"net/http"

	usecases "github.com/savio04/youtube-video-summarizer/domains/video/useCases"
	"github.com/savio04/youtube-video-summarizer/internal/database/repositories"
)

type CreateVideoController struct {
	createVideoUseCase *usecases.CreateVideoUseCase
}

func NewCreateVideoController() *CreateVideoController {
	videoRepository := repositories.NewDbVideoRepository()
	createVideoUseCase := usecases.NewCreateVideoUseCase(videoRepository)

	return &CreateVideoController{
		createVideoUseCase,
	}
}

type CreateVideoPayload struct {
	Url        string `json:"url" validate:"required,url"`
	ExternalId string `json:"externalId" validate:"required"`
}

func (controller *CreateVideoController) Handler(writer http.ResponseWriter, request *http.Request) {
	var body CreateVideoPayload

	decoder := json.NewDecoder(request.Body)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		http.Error(writer, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	data, err := controller.createVideoUseCase.Execute(body.Url, body.ExternalId)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(writer).Encode(struct {
			Message string `json:"message"`
		}{Message: err.Error()})

		return
	}

	writer.WriteHeader(http.StatusOK)

	json.NewEncoder(writer).Encode(struct {
		Payload any `json:"payload"`
	}{Payload: data})
}
