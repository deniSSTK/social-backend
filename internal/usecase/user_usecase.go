package usecase

import (
	"context"
	"social-backend/internal/infrastructure/db/repository"
	"social-backend/internal/infrastructure/http/api_dto"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo *repository.UserRepository
}

func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo}
}

func (uc *UserUsecase) Create(ctx context.Context, dto api_dto.PostUsersJSONBody) (uuid.UUID, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	userId := uuid.New()

	if err = uc.userRepo.Create(ctx, dto, string(hashBytes), userId); err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

func (uc *UserUsecase) Login(ctx context.Context, dto api_dto.PostUsersLogInJSONRequestBody) (uuid.UUID, error) {
	user, err := uc.userRepo.GetPasswordHashByEmailOrUsername(ctx, dto)
	if err != nil {
		return uuid.Nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return uuid.Nil, err
	}

	return user.Id, nil
}
