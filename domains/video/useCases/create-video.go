package usecases

import (
	"fmt"
	"net/url"

	"github.com/savio04/youtube-video-summarizer/domains/video/entities"
	"github.com/savio04/youtube-video-summarizer/domains/video/repositories"
)

type CreateVideoUseCase struct {
	videoRepository repositories.VideoRepository
}

func NewCreateVideoUseCase(repo repositories.VideoRepository) *CreateVideoUseCase {
	return &CreateVideoUseCase{
		videoRepository: repo,
	}
}

func (useCase *CreateVideoUseCase) Execute(videoUrl string) (*entities.Video, error) {
	parsedURL, err := url.Parse(videoUrl)
	if err != nil {
		return nil, err
	}

	queryParams := parsedURL.Query()
	videoID := queryParams.Get("v")

	if videoID == "" {
		return nil, fmt.Errorf("video ID not found in the URL")
	}

	newVideoData := &entities.Video{
		Id:         nil,
		ExternalId: &videoID,
		Url:        videoUrl,
		Summary:    nil,
	}

	params := &repositories.FindOneVideoParams{
		ExternalId: videoID,
		Url:        nil,
	}

	videoAlreadyExists, err := useCase.videoRepository.FindOne(params)
	if err != nil {
		return nil, err
	}

	if videoAlreadyExists != nil {
		return videoAlreadyExists, nil
	}

	newVideo, err := useCase.videoRepository.Create(newVideoData)
	if err != nil {
		return nil, err
	}

	return newVideo, nil
}
