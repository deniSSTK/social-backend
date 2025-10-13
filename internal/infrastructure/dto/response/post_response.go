package response

import (
	"time"

	"github.com/google/uuid"
)

type GetFeedPostByUserId struct {
	Id           uuid.UUID   `json:"id"`
	Description  string      `json:"description"`
	ImagesUrls   []string    `json:"images_urls"`
	LikeCount    int         `json:"like_count"`
	CommentCount int         `json:"comment_count"`
	Author       GetUserInfo `json:"author"`
	CreatedAt    time.Time   `json:"created_at"`
}
