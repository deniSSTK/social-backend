package usecase

import (
	"context"
	"social-backend/internal/domain/image"
	"social-backend/internal/domain/post"
	"social-backend/internal/infrastructure/db/repository"
	"social-backend/internal/infrastructure/dto/request"
	"social-backend/internal/infrastructure/execer"
	"social-backend/internal/infrastructure/imgbb"
	"social-backend/internal/infrastructure/tx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
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
	return tx.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) error {
		var uploadedImages []image.Image

		if len(dto.Images) == 1 {
			img, err := uc.imgBBService.Upload(dto.Images[0])
			if err != nil {
				return err
			}

			uploadedImages = append(uploadedImages, img)
		} else {
			var err error

			uploadedImages, err = uc.imgBBService.UploadImages(ctx, dto.Images)
			if err != nil {
				return err
			}
		}

		//TODO defer deleteImages

		var err error
		dto.TargetPost.Id, err = uc.postRepo.InsertTx(ctx, exec, dto.TargetPost)
		if err != nil {
			return err
		}

		for i := range len(uploadedImages) {
			uploadedImages[i-1].PostId = &dto.TargetPost.Id
			uploadedImages[i-1].Position = &i
		}

		if len(uploadedImages) == 1 {
			if err = uc.imageRepo.InsertTx(ctx, exec, uploadedImages[0]); err != nil {
				return err
			}
		}

		if err = uc.uploadHashtags(ctx, exec, dto.Hashtags, dto.TargetPost.Id); err != nil {
			return err
		}

		return nil
	})
}

func (uc *PostUsecase) uploadHashtags(ctx context.Context, exec execer.Execer, hashtags []request.InsertPostHashtag, postId uuid.UUID) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(len(hashtags))

	for _, h := range hashtags {
		h := h
		g.Go(func() error {
			if h.Id == nil {
				var err error
				*h.Id, err = uc.hashtagRepo.InsertTx(ctx, exec, h.Text)
				if err != nil {
					return err
				}
			}

			if err := uc.postRepo.InsertHashtagTx(ctx, exec, post.Hashtag{HashtagId: *h.Id, PostId: postId, Position: h.Position}); err != nil {
				return err
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

//func (uc *PostUsecase) GetById(ctx context.Context, postId uuid.UUID) (post.Post, error) {}

//func (uc *PostUsecase) GetUserPosts(ctx context.Context, postId uuid.UUID) (post.Post, error) {}
