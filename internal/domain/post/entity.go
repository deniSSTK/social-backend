package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id            uuid.UUID `json:"id"`
	Description   string    `json:"description" binding:"max=1000"`
	AuthorId      uuid.UUID `json:"authorId"`
	CreatedAt     time.Time `json:"createdAt"`
	LikesCount    int       `json:"likesCount"`
	CommentsCount int       `json:"commentsCount"`
	CloseFriends  *bool     `json:"closeFriends"`
	Pinned        *bool     `json:"pinned,omitempty"`
	FirstImage    string    `json:"firstImage"`
}

type Comment struct {
	Id        uuid.UUID `json:"id"`
	PostId    uuid.UUID `json:"postId"`
	AuthorId  uuid.UUID `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	Text      string    `json:"text" binding:"max=500"`
}

type Like struct {
	PostId    uuid.UUID `json:"postId"`
	AuthorId  uuid.UUID `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
}

type Hashtag struct {
	PostId    uuid.UUID `json:"postId"`
	HashtagId uuid.UUID `json:"hashtagId"`
	Position  int       `json:"position"`
}
