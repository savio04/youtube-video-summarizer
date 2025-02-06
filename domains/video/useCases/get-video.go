package usecases

import (
	"github.com/savio04/youtube-video-summarizer/domains/video/entities"
	"github.com/savio04/youtube-video-summarizer/domains/video/repositories"
)

type GetVideoUseCase struct {
	videoRepository repositories.VideoRepository
}

func NewGetVideoUseCase(repo repositories.VideoRepository) *GetVideoUseCase {
	return &GetVideoUseCase{
		videoRepository: repo,
	}
}

func (useCase *GetVideoUseCase) Execute(slug *string) (*entities.Video, error) {
	params := &repositories.FindOneVideoParams{
		ExternalId: slug,
		Url:        nil,
	}

	data, err := useCase.videoRepository.FindOne(params)
	if err != nil {
		return nil, err
	}

	return data, nil
}
