package repository

import (
	"context"
	"social-backend/internal/infrastructure/execer"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FollowRepository struct {
	conn *pgxpool.Pool
}

func NewFollowRepository(conn *pgxpool.Pool) *FollowRepository {
	return &FollowRepository{conn}
}

func (r *FollowRepository) InsertTx(ctx context.Context, exec execer.Execer, followerId, followToId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		INSERT INTO followings
		(follower_id, follow_to_id)
		VALUES ($1, $2)
	`, followerId, followToId)
	return err
}

func (r *FollowRepository) UpdateFollowerCountTx(ctx context.Context, exec execer.Execer, userId uuid.UUID, count int) error {
	_, err := exec.Exec(ctx, `
		UPDATE users
		SET followers = followers + $1
		WHERE id = $2
	`, count, userId)
	return err
}

func (r *FollowRepository) UpdateFollowingCountTx(ctx context.Context, exec execer.Execer, userId uuid.UUID, count int) error {
	_, err := exec.Exec(ctx, `
		UPDATE users
		SET following = following + $1
		WHERE id = $2
	`, count, userId)
	return err
}

func (r *FollowRepository) DeleteTx(ctx context.Context, exec execer.Execer, followerId, followToId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		DELETE FROM followings
		WHERE follower_id = $1 AND follow_to_id = $2
	`, followerId, followToId)
	return err
}
