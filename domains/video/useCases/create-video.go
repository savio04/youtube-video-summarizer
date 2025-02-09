package usecases

import (
	"fmt"

	"github.com/savio04/youtube-video-summarizer/domains/video/entities"
	"github.com/savio04/youtube-video-summarizer/domains/video/repositories"
	"github.com/savio04/youtube-video-summarizer/internal/queue"
	"github.com/savio04/youtube-video-summarizer/internal/utils"
)

type CreateVideoUseCase struct {
	videoRepository repositories.VideoRepository
}

func NewCreateVideoUseCase(repo repositories.VideoRepository) *CreateVideoUseCase {
	return &CreateVideoUseCase{
		videoRepository: repo,
	}
}

func (useCase *CreateVideoUseCase) Execute(videoUrl string, externalId string) (*entities.Video, error) {
	fmt.Println("externalId", externalId)

	videoAlreadyExists, err := useCase.findVideo(&externalId)
	if err != nil {
		return nil, err
	}

	if videoAlreadyExists != nil {
		if *videoAlreadyExists.Status == "FAILED" {
			newStatus := "PROCESSING"

			videoAlreadyExists.Status = &newStatus

			useCase.videoRepository.UpdateByExternalId(*videoAlreadyExists.ExternalId, &repositories.UpdateParams{
				Status: &newStatus,
			})

			errQueue := queue.InsertIntoQueue(utils.QueueVideoProcessing, externalId)
			if errQueue != nil {
				return nil, errQueue
			}
		}

		return videoAlreadyExists, nil
	}

	newVideoData := &entities.Video{
		Id:         nil,
		ExternalId: &externalId,
		Url:        videoUrl,
		Summary:    nil,
	}

	newVideo, err := useCase.videoRepository.Create(newVideoData)
	if err != nil {
		return nil, err
	}

	errQueue := queue.InsertIntoQueue(utils.QueueVideoProcessing, externalId)
	if errQueue != nil {
		return nil, errQueue
	}

	return newVideo, nil
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
