package repository

import (
	"context"
	"social-backend/internal/domain/user"
	"social-backend/internal/infrastructure/dto/request"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{conn}
}

func (r *UserRepository) Insert(ctx context.Context, dto request.CreateUser, passwordHash string, userId uuid.UUID) error {
	_, err := r.conn.Exec(ctx, `
		INSERT INTO users (id, username, email, password_hash) 
		VALUES ($1, $2, $3, $4)
	`, userId, dto.Username, dto.Email, passwordHash)
	return err
}

func (r *UserRepository) GetPasswordHashByEmailOrUsername(ctx context.Context, dto request.LogIn) (user.User, error) {
	var targetUser user.User
	if err := r.conn.QueryRow(ctx, `
		SELECT id, password_hash
		FROM users 
		WHERE email = $1 OR username = $1
	`, dto.EmailOrUsername).Scan(&targetUser.Id, &targetUser.PasswordHash); err != nil {
		return user.User{}, err
	}
	return targetUser, nil
}

func (r *UserRepository) GetUsernameById(ctx context.Context, userId uuid.UUID) (string, error) {
	var username string
	if err := r.conn.QueryRow(ctx, `
		SELECT username
		FROM users
		WHERE id = $1
	`, userId).Scan(&username); err != nil {
		return "", err
	}

	return username, nil
}
