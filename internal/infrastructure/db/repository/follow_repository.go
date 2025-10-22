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

func (r *FollowRepository) FollowTx(ctx context.Context, exec execer.Execer, userId, followToId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		INSERT INTO followings
		(follower_id, follow_to_id)
		VALUES ($1, $2)
	`, userId, followToId)
	return err
}

func (r *FollowRepository) AddFollowerTx(ctx context.Context, exec execer.Execer, userId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE users
		SET followers = followers + 1
		WHERE id = $1
	`, userId)
	return err
}

func (r *FollowRepository) AddFollowingTx(ctx context.Context, exec execer.Execer, userId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE users
		SET followings = followings + 1
		WHERE id = $1
	`, userId)
	return err
}
