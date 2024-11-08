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
	videoId, err := useCase.getVideoIdByUrl(videoUrl)
	if err != nil {
		return nil, err
	}

	videoAlreadyExists, err := useCase.findVideo(videoId)
	if err != nil {
		return nil, err
	}

	if videoAlreadyExists != nil {
		return videoAlreadyExists, nil
	}

	newVideoData := &entities.Video{
		Id:         nil,
		ExternalId: videoId,
		Url:        videoUrl,
		Summary:    nil,
	}

	newVideo, err := useCase.videoRepository.Create(newVideoData)
	if err != nil {
		return nil, err
	}

	// TODO: Send video to queue

	return newVideo, nil
}

func (useCase *CreateVideoUseCase) getVideoIdByUrl(videoUrl string) (*string, error) {
	parsedURL, err := url.Parse(videoUrl)
	if err != nil {
		return nil, err
	}

	queryParams := parsedURL.Query()
	videoID := queryParams.Get("v")

	if videoID == "" {
		return nil, fmt.Errorf("video ID not found in the URL")
	}

	return &videoID, nil
}

func (useCase *CreateVideoUseCase) findVideo(videoId *string) (*entities.Video, error) {
	params := &repositories.FindOneVideoParams{
		ExternalId: videoId,
		Url:        nil,
	}

	videoAlreadyExists, err := useCase.videoRepository.FindOne(params)
	if err != nil {
		return nil, err
	}

	if videoAlreadyExists != nil {
		return videoAlreadyExists, nil
	}

	return nil, nil
}
