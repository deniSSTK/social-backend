package request

import (
	"io"
	"social-backend/internal/domain/post"

	"github.com/google/uuid"
)

type InsertPostHashtag struct {
	Text     string     `json:"text"`
	Id       *uuid.UUID `json:"id,omitempty"`
	Position int        `json:"position"`
}

type InsertPost struct {
	TargetPost post.Post            `json:"targetPost"`
	Images     []io.Reader          `json:"images"`
	Hashtags   *[]InsertPostHashtag `json:"hashtags"`
}
