package hashtag

import "github.com/google/uuid"

type Hashtag struct {
	Id   uuid.UUID `json:"id"`
	Text string    `json:"text" binding:"max=100"`
}
