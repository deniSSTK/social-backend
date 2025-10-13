package repository

import (
	"context"
	"social-backend/internal/infrastructure/execer"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HashtagRepository struct {
	conn *pgxpool.Pool
}

func NewHashtagRepository(conn *pgxpool.Pool) *HashtagRepository {
	return &HashtagRepository{conn}
}

func (repo *HashtagRepository) InsertTx(ctx context.Context, exec execer.Execer, text string) (uuid.UUID, error) {
	id := uuid.New()
	_, err := exec.Exec(ctx, `
		INSERT INTO hashtags
		(id, text)
		VALUES ($1, $2)
	`, id, text)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
