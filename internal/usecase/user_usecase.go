package usecase

import (
	"context"
	"social-backend/internal/infrastructure/db/repository"
	"social-backend/internal/infrastructure/http/api_dto"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo *repository.UserRepository
}

func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo}
}

func (uc *UserUsecase) Create(ctx context.Context, dto api_dto.PostUsersJSONBody) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err = uc.userRepo.Create(ctx, dto, string(hashBytes)); err != nil {
		return err
	}

	return nil
}
