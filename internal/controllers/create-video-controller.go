package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/savio04/youtube-video-summarizer/internal/database"
	"github.com/savio04/youtube-video-summarizer/internal/database/repositories"
	usecases "github.com/savio04/youtube-video-summarizer/internal/domains/video/useCases"
)

type CreateVideoController struct {
	createVideoUseCase *usecases.CreateVideoUseCase
}

func NewCreateVideoController() *CreateVideoController {
	videoRepository := repositories.NewDbVideoRepository(database.Db)
	createVideoUseCase := usecases.NewCreateVideoUseCase(videoRepository)

	return &CreateVideoController{
		createVideoUseCase,
	}
}

func (controller *CreateVideoController) Handler(writer http.ResponseWriter, request *http.Request) {
	controller.createVideoUseCase.Execute()

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	json.NewEncoder(writer).Encode(struct {
		Message string `json:"message"`
	}{Message: "Video created"})
}
