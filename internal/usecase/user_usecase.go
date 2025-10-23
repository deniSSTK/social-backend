package usecase

import (
	"context"
	"social-backend/internal/infrastructure/db/repository"
	"social-backend/internal/infrastructure/dto/request"
	"social-backend/internal/infrastructure/dto/response"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo *repository.UserRepository
}

func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo}
}

func (uc *UserUsecase) Insert(ctx context.Context, dto request.CreateUser) (uuid.UUID, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	userId := uuid.New()

	if err = uc.userRepo.Insert(ctx, dto, string(hashBytes), userId); err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

func (uc *UserUsecase) Login(ctx context.Context, dto request.LogIn) (uuid.UUID, error) {
	user, err := uc.userRepo.GetPasswordHashByEmailOrUsername(ctx, dto)
	if err != nil {
		return uuid.Nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return uuid.Nil, err
	}

	return user.Id, nil
}

func (uc *UserUsecase) GetUsernameById(ctx context.Context, userId uuid.UUID) (string, error) {
	return uc.userRepo.GetUsernameById(ctx, userId)
}

func (uc *UserUsecase) GetUserInfoByName(ctx context.Context, username string, currentUserId uuid.UUID) (response.GetUserInfo, error) {
	return uc.userRepo.GetUserInfoByName(ctx, username, currentUserId)
}
