package usecase

import (
	"context"
	"social-backend/internal/domain/image"
	"social-backend/internal/domain/post"
	"social-backend/internal/infrastructure/db/repository"
	"social-backend/internal/infrastructure/dto/request"
	"social-backend/internal/infrastructure/dto/response"
	"social-backend/internal/infrastructure/errors"
	"social-backend/internal/infrastructure/execer"
	"social-backend/internal/infrastructure/imgbb"
	"social-backend/internal/infrastructure/logger"
	"social-backend/internal/infrastructure/tx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostUsecase struct {
	baseRepo     *repository.BaseRepo
	postRepo     *repository.PostRepository
	imageRepo    *repository.ImageRepository
	hashtagRepo  *repository.HashtagRepository
	imgBBService *imgbb.ImgBBService
}

func NewPostUsecase(
	baseRepo *repository.BaseRepo,
	postRepo *repository.PostRepository,
	imageRepo *repository.ImageRepository,
	hashtagRepo *repository.HashtagRepository,
	imgBBService *imgbb.ImgBBService,
) *PostUsecase {
	return &PostUsecase{
		baseRepo,
		postRepo,
		imageRepo,
		hashtagRepo,
		imgBBService,
	}
}

func (uc *PostUsecase) Insert(ctx context.Context, dto request.InsertPost) error {
	return tx.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (err error) {
		var uploadedImages []image.Image

		if len(dto.Images) == 1 {
			img, uploadErr := uc.imgBBService.Upload(dto.Images[0])
			if uploadErr != nil {
				return uploadErr
			}
			uploadedImages = append(uploadedImages, img)
		} else if len(dto.Images) > 1 {
			images, uploadErr := uc.imgBBService.UploadImages(ctx, dto.Images)
			if uploadErr != nil {
				return uploadErr
			}
			uploadedImages = images
		}

		defer func() {
			if err != nil && len(uploadedImages) > 0 {
				for _, img := range uploadedImages {
					success, delErr := uc.imgBBService.DeleteImage(img.DeleteUrl)
					if delErr != nil {
						logger.Get().Error(delErr.Error())
					}

					if !success {
						logger.Get().Error(errors.ImgBBUDeletingError.Error())
					}
				}
			}
		}()

		dto.TargetPost.Id, err = uc.postRepo.InsertTx(ctx, exec, dto.TargetPost)
		if err != nil {
			return err
		}

		for i := range uploadedImages {
			pos := new(int)
			*pos = i
			uploadedImages[i].PostId = &dto.TargetPost.Id
			uploadedImages[i].Position = pos
		}

		for _, img := range uploadedImages {
			if err = uc.imageRepo.InsertTx(ctx, exec, img); err != nil {
				return err
			}
		}

		if dto.Hashtags != nil {
			if err = uc.uploadHashtags(ctx, exec, *dto.Hashtags, dto.TargetPost.Id); err != nil {
				return err
			}
		}

		return nil
	})
}

func (uc *PostUsecase) uploadHashtags(ctx context.Context, exec execer.Execer, hashtags []request.InsertPostHashtag, postId uuid.UUID) error {
	for _, h := range hashtags {
		if h.Id == nil {
			id, err := uc.hashtagRepo.InsertTx(ctx, exec, h.Text)
			if err != nil {
				return err
			}
			h.Id = &id
		}

		if err := uc.postRepo.InsertHashtagTx(ctx, exec, post.Hashtag{
			HashtagId: *h.Id,
			PostId:    postId,
			Position:  h.Position,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (uc *PostUsecase) GetById(ctx context.Context, postId uuid.UUID) (post.Post, error) {
	return uc.postRepo.GetById(ctx, postId)
}

func (uc *PostUsecase) GetUserPosts(ctx context.Context, userId uuid.UUID, offset int) ([]response.GetPostByUserId, error) {
	return uc.postRepo.GetUserPosts(ctx, userId, offset)
}

func (uc *PostUsecase) GetPostCountsById(ctx context.Context, postId uuid.UUID) (response.GetPostCountsById, error) {
	return uc.postRepo.GetPostCountsById(ctx, postId)
}

func (uc *PostUsecase) GetFeedPosts(ctx context.Context, userId uuid.UUID, offset int) ([]response.GetFeedPostByUserId, error) {
	return uc.postRepo.GetFeedPosts(ctx, userId, offset)
}

func (uc *PostUsecase) LikePostManipulation(ctx context.Context, postId, userId uuid.UUID, count int) error {
	return tx.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (err error) {
		if count == 1 {
			if err = uc.postRepo.LikePostTx(ctx, exec, postId, userId); err != nil {
				return err
			}
		} else if count == -1 {
			if err = uc.postRepo.RemoveLikePostTx(ctx, exec, postId, userId); err != nil {
				return err
			}
		}

		if err = uc.postRepo.UpdatePostLikesCountTx(ctx, exec, count, postId); err != nil {
			return err
		}

		return nil
	})
}
