package repositories

import "github.com/savio04/first-api/domains/auth/entities"

type VideoRepository interface {
	Create(video *entities.Video) (*entities.Video, error)
}
