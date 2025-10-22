package handler

import (
	"net/http"
	"social-backend/internal/infrastructure/auth"
	"social-backend/internal/infrastructure/http/context"
	"social-backend/internal/infrastructure/http/middleware"
	"social-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type HashtagHandler struct {
	hashtagUC   *usecase.HashtagUsecase
	authService *auth.AuthService
}

func NewHashtagHandler(hashtagUC *usecase.HashtagUsecase, authService *auth.AuthService) *HashtagHandler {
	return &HashtagHandler{hashtagUC, authService}
}

func (h *HashtagHandler) RegisterRoutes(group *gin.RouterGroup) {
	protected := group.Group("/hashtag", middleware.AuthMiddleware(h.authService))

	protected.GET("/by-name/:"+string(context.ContextParamValue), h.getByName)
}

func (h *HashtagHandler) getByName(c *gin.Context) {
	value := c.Param(string(context.ContextParamValue))
	if value == "" {
		HandleError(c, http.StatusBadRequest, nil)
		return
	}

	hashtags, err := h.hashtagUC.GetByName(c.Request.Context(), value)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, hashtags)
}
