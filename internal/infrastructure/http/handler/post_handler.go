package handler

import (
	"net/http"
	"social-backend/internal/infrastructure/auth"
	"social-backend/internal/infrastructure/dto/request"
	"social-backend/internal/infrastructure/http/middleware"
	"social-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postUC     *usecase.PostUsecase
	jwtService auth.JWTService
}

func NewPostHandler(postUC *usecase.PostUsecase, jwtService auth.JWTService) *PostHandler {
	return &PostHandler{postUC, jwtService}
}

func (h *PostHandler) RegisterRoutes(router *gin.RouterGroup) {
	protected := router.Group("/posts", middleware.JWTMiddleware(h.jwtService))

	protected.POST("/", h.createPost)
}

func (h *PostHandler) createPost(c *gin.Context) {
	var dto request.InsertPost
	if err := c.ShouldBindJSON(&dto); err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.postUC.Insert(c.Request.Context(), dto); err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}
