package repositories

import "github.com/savio04/youtube-video-summarizer/domains/video/entities"

type FindOneVideoParams struct {
	Url        *string
	ExternalId *string
}

type UpdateParams struct {
	Summary *string
	Status  *string
}

type VideoRepository interface {
	Create(video *entities.Video) (*entities.Video, error)
	FindOne(params *FindOneVideoParams) (*entities.Video, error)
	UpdateByExternalId(externalId string, params *UpdateParams) error
}
