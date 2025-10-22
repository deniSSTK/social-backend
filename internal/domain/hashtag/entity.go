package hashtag

import "github.com/google/uuid"

type Hashtag struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name" binding:"max=100"`
}
