package repository

import (
	"context"
	"social-backend/internal/domain/image"
	"social-backend/internal/infrastructure/execer"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ImageRepository struct {
	conn *pgxpool.Pool
}

func NewImageRepository(conn *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{conn}
}

func (r *ImageRepository) InsertTx(ctx context.Context, exec execer.Execer, image image.Image) error {
	_, err := exec.Exec(ctx, `
		INSERT INTO images 
		(url, position, post_id, delete_url)
		VALUES ($1, $2, $3, $4)
	`, image.Url, image.Position, image.PostId, image.DeleteUrl)
	return err
}
