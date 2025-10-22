package usecase

import (
	"context"
	"social-backend/internal/infrastructure/db/repository"
	"social-backend/internal/infrastructure/tx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FollowUsecase struct {
	baseRepo   *repository.BaseRepo
	followRepo *repository.FollowRepository
}

func NewFollowUsecase(baseRepo *repository.BaseRepo, followRepo *repository.FollowRepository) *FollowUsecase {
	return &FollowUsecase{baseRepo, followRepo}
}

func (uc *FollowUsecase) Follow(ctx context.Context, userId, followToId uuid.UUID) error {
	return tx.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (err error) {
		if err = uc.followRepo.FollowTx(ctx, exec, userId, followToId); err != nil {
			return err
		}

		if err = uc.followRepo.AddFollowerTx(ctx, exec, followToId); err != nil {
			return err
		}

		if err = uc.followRepo.AddFollowingTx(ctx, exec, userId); err != nil {
			return err
		}

		return nil
	})
}
