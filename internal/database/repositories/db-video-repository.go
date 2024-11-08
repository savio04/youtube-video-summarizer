package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/savio04/youtube-video-summarizer/domains/video/entities"
	"github.com/savio04/youtube-video-summarizer/domains/video/repositories"
	"github.com/savio04/youtube-video-summarizer/internal/database"
)

type DbVideoRepository struct{}

func NewDbVideoRepository() *DbVideoRepository {
	return &DbVideoRepository{}
}

func (dbVideo *DbVideoRepository) Create(video *entities.Video) (*entities.Video, error) {
	newVideo := &entities.Video{}

	err := database.Db.QueryRow(context.Background(), `INSERT INTO videos (external_id, url) VALUES($1, $2) RETURNING id, url, external_id`, video.ExternalId, video.Url).Scan(&newVideo.Id, &newVideo.Url, &newVideo.ExternalId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return newVideo, nil
}

func (dbVideo *DbVideoRepository) FindOne(params *repositories.FindOneVideoParams) (*entities.Video, error) {
	newVideo := &entities.Video{}

	err := database.Db.QueryRow(
		context.Background(),
		`SELECT id, external_id, url, summary FROM videos WHERE external_id=$1`,
		params.ExternalId,
	).Scan(&newVideo.Id, &newVideo.ExternalId, &newVideo.Url, &newVideo.Summary)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return newVideo, nil
}
