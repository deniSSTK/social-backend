package handler

import (
	"net/http"
	"social-backend/internal/infrastructure/auth"
	"social-backend/internal/infrastructure/dto/request"
	"social-backend/internal/infrastructure/errors"
	"social-backend/internal/infrastructure/http/context"
	"social-backend/internal/infrastructure/http/middleware"
	"social-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	postUC      *usecase.PostUsecase
	authService *auth.AuthService
}

func NewPostHandler(postUC *usecase.PostUsecase, authService *auth.AuthService) *PostHandler {
	return &PostHandler{postUC, authService}
}

func (h *PostHandler) RegisterRoutes(router *gin.RouterGroup) {
	protected := router.Group("/posts", middleware.AuthMiddleware(h.authService))

	protected.GET("/:"+string(context.ContextParamPostId), h.getById)
	protected.GET("/user/:"+string(context.ContextParamUserId)+"/:"+string(context.ContextParamOffset), h.getByUserId)

	protected.POST("/", h.createPost)
}

func (h *PostHandler) createPost(c *gin.Context) {
	var dto request.InsertPost
	if err := c.ShouldBindJSON(&dto); err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	userId := context.GetContextUserId(c)
	if userId == uuid.Nil {
		HandleError(c, http.StatusUnauthorized, errors.ContextUserIdEmpty)
		return
	}

	dto.TargetPost.AuthorId = userId

	if err := h.postUC.Insert(c.Request.Context(), dto); err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *PostHandler) getById(c *gin.Context) {
	postId, err := context.GetContextParamUUID(c, context.ContextParamPostId)
	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	post, err := h.postUC.GetById(c.Request.Context(), postId)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) getByUserId(c *gin.Context) {
	userId, err := context.GetContextParamUUID(c, context.ContextParamUserId)
	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	offset, err := context.GetContextParamInt(c, context.ContextParamOffset)
	if err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	posts, err := h.postUC.GetUserPosts(c.Request.Context(), userId, offset)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}
