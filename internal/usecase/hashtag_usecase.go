package usecase

import (
	"context"
	"social-backend/internal/domain/hashtag"
	"social-backend/internal/infrastructure/db/repository"
)

type HashtagUsecase struct {
	hashtagRepo *repository.HashtagRepository
}

func NewHashtagUsecase(hashtagRepo *repository.HashtagRepository) *HashtagUsecase {
	return &HashtagUsecase{hashtagRepo}
}

func (uc *HashtagUsecase) GetByName(ctx context.Context, name string) ([]hashtag.Hashtag, error) {
	return uc.hashtagRepo.GetByName(ctx, name)
}
