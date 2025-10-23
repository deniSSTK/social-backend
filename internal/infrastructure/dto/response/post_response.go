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

type GetPostCountsById struct {
	LikesCount    int `json:"likesCount"`
	CommentsCount int `json:"commentsCount"`
}

type GetFeedPostByUserId struct {
	Id                 uuid.UUID   `json:"id"`
	Description        string      `json:"description"`
	ImagesUrls         []string    `json:"imagesUrls"`
	ImagesCount        int         `json:"imagesCount"`
	LikesCount         int         `json:"likesCount"`
	CommentsCount      int         `json:"commentsCount"`
	Author             UserPreview `json:"author"`
	CreatedAt          time.Time   `json:"createdAt"`
	IfCurrentUserLiked bool        `json:"ifCurrentUserLiked"`
}
