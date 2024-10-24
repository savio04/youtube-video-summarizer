package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/savio04/youtube-video-summarizer/internal/domains/video/entities"
)

type DbVideoRepository struct {
	db *pgxpool.Pool
}

func NewDbVideoRepository(db *pgxpool.Pool) *DbVideoRepository {
	return &DbVideoRepository{
		db: db,
	}
}

func (dbVideo *DbVideoRepository) Create(video *entities.Video) (*entities.Video, error) {
	return nil, nil
}
