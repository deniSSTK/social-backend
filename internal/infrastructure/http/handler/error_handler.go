package handler

import (
	"fmt"
	"social-backend/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleError(c *gin.Context, status int, err error, dto ...any) {
	if err == nil {
		err = fmt.Errorf("unknown error")
	}

	if dto != nil {
		logger.Get().Error(err.Error(),
			zap.Int("status", status),
			zap.Any("dto", dto))
	} else {
		logger.Get().Error(err.Error(), zap.Int("status", status))
	}

	c.JSON(status, gin.H{"error": err.Error()})
}
