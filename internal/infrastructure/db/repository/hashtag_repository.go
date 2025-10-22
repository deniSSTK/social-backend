package repository

import (
	"context"
	"social-backend/internal/domain/hashtag"
	"social-backend/internal/infrastructure/execer"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HashtagRepository struct {
	conn *pgxpool.Pool
}

func NewHashtagRepository(conn *pgxpool.Pool) *HashtagRepository {
	return &HashtagRepository{conn}
}

func (r *HashtagRepository) InsertTx(ctx context.Context, exec execer.Execer, text string) (uuid.UUID, error) {
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

func (r *HashtagRepository) GetByName(ctx context.Context, name string) ([]hashtag.Hashtag, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, name
		FROM hashtags
		WHERE text ILIKE $1
		LIMIT 5
	`, name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hashtags, err := pgx.CollectRows(rows, pgx.RowToStructByName[hashtag.Hashtag])
	if err != nil {
		return nil, err
	}

	return hashtags, nil
}
