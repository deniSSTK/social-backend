package image

import "github.com/google/uuid"

type Image struct {
	Id        uuid.UUID  `json:"id"`
	Url       string     `json:"url"`
	Position  *int       `json:"position"`
	PostId    *uuid.UUID `json:"postId"`
	DeleteUrl string     `json:"deleteUrl"`
}
