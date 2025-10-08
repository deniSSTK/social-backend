package repository

import (
	"context"
	"social-backend/internal/infrastructure/http/api_dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{conn}
}

func (r UserRepository) Create(ctx context.Context, dto api_dto.PostUsersJSONBody, passwordHash string) error {
	_, err := r.conn.Exec(ctx, `
		INSERT INTO users (username, email, password_hash) 
		VALUES ($1, $2, $3)
	`, dto.Username, dto.Email, passwordHash)
	return err
}
