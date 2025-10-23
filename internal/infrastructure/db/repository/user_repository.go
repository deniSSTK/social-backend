package repository

import (
	"context"
	"social-backend/internal/domain/user"
	"social-backend/internal/infrastructure/dto/request"
	"social-backend/internal/infrastructure/dto/response"

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

func (r *UserRepository) GetUserInfoByName(ctx context.Context, username string, currentUserId uuid.UUID) (response.GetUserInfo, error) {
	var res response.GetUserInfo
	if err := r.conn.QueryRow(ctx, `
		SELECT 
		    u.id, 
		    icon_url, 
		    description, 
		    followers, 
		    following, 
		    post_count,
		    EXISTS (
		        SELECT 1
		        FROM followings f
		        WHERE f.follow_to_id = u.id AND f.follower_id = $2
  		    ) AS if_current_user_followed
		FROM users u
		WHERE username = $1
	`, username, currentUserId).Scan(
		&res.Id,
		&res.IconUrl,
		&res.Description,
		&res.Followers,
		&res.Following,
		&res.PostCount,
		&res.IfCurrentUserFollowed,
	); err != nil {
		return response.GetUserInfo{}, err
	}

	return res, nil
}

func (r *UserRepository) CheckIfExistsById(ctx context.Context, userId uuid.UUID) (bool, error) {
	var exists bool
	if err := r.conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE id = $1
		)
	`, userId).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
