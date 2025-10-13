package tx

import (
	"context"
	"social-backend/internal/infrastructure/db/repository"

	"github.com/jackc/pgx/v5"
)

func WithTx[T any](ctx context.Context, base *repository.BaseRepo, fn func(ctx context.Context, exec pgx.Tx) (T, error)) (res T, err error) {
	tx, err := base.BeginTx(ctx)
	if err != nil {
		return res, err
	}

	defer func() {
		_ = txReturn(tx, ctx, err)
	}()

	return fn(ctx, tx)
}

func WithTxVoid(ctx context.Context, base *repository.BaseRepo, fn func(ctx context.Context, exec pgx.Tx) error) error {
	tx, err := base.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = txReturn(tx, ctx, err)
	}()

	return fn(ctx, tx)
}

func txReturn(tx pgx.Tx, ctx context.Context, err error) error {
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
		}
		return err
	} else {
		return tx.Commit(ctx)
	}
}
