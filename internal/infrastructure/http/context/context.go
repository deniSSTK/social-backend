package context

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContextValues string

const (
	ContextUserId ContextValues = "userId"
)

func GetContextUserId(c *gin.Context) uuid.UUID {
	userId, exists := c.Get(ContextUserId)
	if !exists {
		return uuid.Nil
	}

	id, ok := userId.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return id
}
