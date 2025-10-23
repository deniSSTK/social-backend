package handler

import (
	"net/http"
	"social-backend/internal/infrastructure/auth"
	"social-backend/internal/infrastructure/auth/cookie"
	"social-backend/internal/infrastructure/dto/request"
	"social-backend/internal/infrastructure/errors"
	"social-backend/internal/infrastructure/http/context"
	"social-backend/internal/infrastructure/http/middleware"
	"social-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUC      *usecase.UserUsecase
	authService *auth.AuthService
}

type AuthUser struct {
	UserId uuid.UUID `json:"userId"`
}

func NewUserHandler(userUC *usecase.UserUsecase, authService *auth.AuthService) *UserHandler {
	return &UserHandler{userUC, authService}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/users")

	group.POST("", h.createUser)
	group.POST("/log-in", h.login)

	protected := group.Group("/", middleware.AuthMiddleware(h.authService))

	protected.GET("/auth", h.authCheck)
	protected.GET("/id/username", h.getUsernameById)
	protected.GET("/info-by/username/:"+string(context.ContextParamUsername), h.getUserInfoByName)
}

func (h *UserHandler) createUser(c *gin.Context) {
	var dto request.CreateUser

	if err := c.ShouldBindJSON(&dto); err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	userId, err := h.userUC.Insert(c, dto)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	token, err := h.authService.JwtService.GenerateToken(userId, auth.OneMonth)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	cookie.SetCookie(c, cookie.JWTTokenCookie, token, auth.OneMonth)

	c.Status(http.StatusCreated)
}

func (h *UserHandler) login(c *gin.Context) {
	var dto request.LogIn

	if err := c.ShouldBindJSON(&dto); err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	userId, err := h.userUC.Login(c, dto)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	if userId == uuid.Nil {
		HandleError(c, http.StatusUnauthorized, errors.ContextUserIdEmpty)
		return
	}

	token, err := h.authService.JwtService.GenerateToken(userId, auth.OneMonth)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	cookie.SetCookie(c, cookie.JWTTokenCookie, token, auth.OneMonth)

	c.Status(http.StatusOK)
}

// add etag to check if user data was modified

func (h *UserHandler) authCheck(c *gin.Context) {
	userId := context.GetContextUserId(c)

	if userId == uuid.Nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	user := AuthUser{
		UserId: userId,
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) getUsernameById(c *gin.Context) {
	userId := context.GetContextUserId(c)

	username, err := h.userUC.GetUsernameById(c, userId)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, username)
}

func (h *UserHandler) getUserInfoByName(c *gin.Context) {
	username, exists := c.Params.Get(string(context.ContextParamUsername))
	if !exists {
		HandleError(c, http.StatusBadRequest, errors.ContextParamNotFound)
		return
	}

	userId := context.GetContextUserId(c)
	if userId == uuid.Nil {
		HandleError(c, http.StatusUnauthorized, errors.ContextUserIdEmpty)
		return
	}

	info, err := h.userUC.GetUserInfoByName(c.Request.Context(), username, userId)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	HandleETag(c, info)
}
