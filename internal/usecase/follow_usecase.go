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

func (uc *FollowUsecase) FollowManipulation(ctx context.Context, userId, followToId uuid.UUID, count int) error {
	return tx.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (err error) {
		if count == 1 {
			if err = uc.followRepo.InsertTx(ctx, exec, userId, followToId); err != nil {
				return err
			}
		} else if count == -1 {
			if err = uc.followRepo.DeleteTx(ctx, exec, userId, followToId); err != nil {
				return err
			}
		}

		if err = uc.followRepo.UpdateFollowerCountTx(ctx, exec, followToId, count); err != nil {
			return err
		}

		if err = uc.followRepo.UpdateFollowingCountTx(ctx, exec, userId, count); err != nil {
			return err
		}

		return nil
	})
}
