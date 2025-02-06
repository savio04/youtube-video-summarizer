package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	usecases "github.com/savio04/youtube-video-summarizer/domains/video/useCases"
	"github.com/savio04/youtube-video-summarizer/internal/database/repositories"
)

type GetVideoController struct {
	getVideoUseCase *usecases.GetVideoUseCase
}

func NewGetVideoController() *GetVideoController {
	videoRepository := repositories.NewDbVideoRepository()
	getVideoUseCase := usecases.NewGetVideoUseCase(videoRepository)

	return &GetVideoController{
		getVideoUseCase: getVideoUseCase,
	}
}

func (controller *GetVideoController) Handler(writer http.ResponseWriter, request *http.Request) {
	userID := chi.URLParam(request, "videoId")

	fmt.Printf("here %s", userID)

	data, err := controller.getVideoUseCase.Execute(&userID)
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
