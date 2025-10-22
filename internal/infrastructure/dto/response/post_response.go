package response

import (
	"time"

	"github.com/google/uuid"
)

type GetPostByUserId struct {
	Id           uuid.UUID `json:"id"`
	FirstImage   string    `json:"firstImage"`
	Pinned       *bool     `json:"pinned"`
	CloseFriends *bool     `json:"closeFriends"`
}

//json names are not camel

type GetFeedPostByUserId struct {
	Id           uuid.UUID   `json:"id"`
	Description  string      `json:"description"`
	ImagesUrls   []string    `json:"images_urls"`
	LikeCount    int         `json:"like_count"`
	CommentCount int         `json:"comment_count"`
	Author       GetUserInfo `json:"author"`
	CreatedAt    time.Time   `json:"created_at"`
}
