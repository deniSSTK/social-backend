package handler

import (
	"net/http"
	"social-backend/internal/infrastructure/auth"
	"social-backend/internal/infrastructure/errors"
	"social-backend/internal/infrastructure/http/context"
	"social-backend/internal/infrastructure/http/middleware"
	"social-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FollowHandler struct {
	followUC    *usecase.FollowUsecase
	authService *auth.AuthService
}

func NewFollowHandler(followUC *usecase.FollowUsecase, authService *auth.AuthService) *FollowHandler {
	return &FollowHandler{followUC, authService}
}

func (h *FollowHandler) RegisterRoutes(router *gin.RouterGroup) {
	protected := router.Group("/follow", middleware.AuthMiddleware(h.authService))

	protected.POST("/:"+string(context.ContextParamUserId), h.followManipulation)

	protected.DELETE("/:"+string(context.ContextParamUserId), h.followManipulation)
}

func (h *FollowHandler) followManipulation(c *gin.Context) {
	userId := context.GetContextUserId(c)
	if userId == uuid.Nil {
		HandleError(c, http.StatusUnauthorized, errors.ContextUserIdEmpty)
		return
	}

	followToId, err := context.GetContextParamUUID(c, context.ContextParamUserId)
	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	var count int
	switch c.Request.Method {
	case http.MethodPost:
		count = 1
	case http.MethodDelete:
		count = -1
	}

	if err = h.followUC.FollowManipulation(c, userId, followToId, count); err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
