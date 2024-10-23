package repositories

import "github.com/savio04/youtube-video-summarizer/internal/domains/video/entities"

type VideoRepository interface {
	Create(video *entities.Video) (*entities.Video, error)
}
