package usecases

import (
	"github.com/savio04/youtube-video-summarizer/internal/domains/video/entities"
	"github.com/savio04/youtube-video-summarizer/internal/domains/video/repositories"
)

type CreateVideoUseCase struct {
	videoRepository repositories.VideoRepository
}

func NewCreateVideoUseCase(repo repositories.VideoRepository) *CreateVideoUseCase {
	return &CreateVideoUseCase{
		videoRepository: repo,
	}
}

func (useCase *CreateVideoUseCase) Execute() {
	newVideo := &entities.Video{
		Id:         "123",
		ExternalId: "123",
		Url:        "Teste",
		Summary:    "Hello",
	}

	useCase.videoRepository.Create(newVideo)
}
