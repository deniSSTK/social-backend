package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID  `json:"id"`
	Username     string     `json:"username" binding:"max=50"`
	Email        string     `json:"email" binding:"max=100,email"`
	PasswordHash string     `json:"passwordHash"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
