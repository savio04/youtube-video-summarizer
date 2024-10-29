package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/savio04/youtube-video-summarizer/domains/video/entities"
	"github.com/savio04/youtube-video-summarizer/domains/video/repositories"
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
	newVideo := &entities.Video{}

	err := dbVideo.db.QueryRow(context.Background(), `INSERT INTO videos (external_id, url) VALUES($1, $2) RETURNING id, url, external_id`, video.ExternalId, video.Url).Scan(&newVideo.Id, &newVideo.Url, &newVideo.ExternalId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return newVideo, nil
}

func (dbVideo *DbVideoRepository) FindOne(params *repositories.FindOneVideoParams) (*entities.Video, error) {
	newVideo := &entities.Video{}

	err := dbVideo.db.QueryRow(
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
