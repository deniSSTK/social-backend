package context

import (
	"errors"
	e "social-backend/internal/infrastructure/errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContextValues string

const (
	ContextUserId ContextValues = "userId"
)

type ContextParams string

const (
	ContextParamUserId   ContextParams = "userId"
	ContextParamUsername ContextParams = "username"

	ContextParamValue ContextParams = "value"

	ContextParamPostId ContextParams = "postId"
	ContextParamOffset ContextParams = "offset"
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

func GetContextParamInt(c *gin.Context, key ContextParams) (int, error) {
	valueStr := c.Param(string(key))
	if valueStr == "" {
		return 0, errors.New(e.ContextParamNotFound.Error() + string(key))
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func GetContextParamUUID(c *gin.Context, key ContextParams) (uuid.UUID, error) {
	idStr := c.Param(string(key))
	if idStr == "" {
		return uuid.Nil, errors.New(e.ContextParamNotFound.Error() + string(key))
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
