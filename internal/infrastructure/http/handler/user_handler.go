package handler

import (
	"net/http"
	"social-backend/internal/infrastructure/auth"
	"social-backend/internal/infrastructure/http/api_dto"
	"social-backend/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUC     *usecase.UserUsecase
	jwtService auth.JWTService
}

func NewUserHandler(userUC *usecase.UserUsecase, jwtService auth.JWTService) *UserHandler {
	return &UserHandler{userUC, jwtService}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/users")

	group.POST("", h.createUser)
	group.POST("/log-in", h.login)

	//protected := group.Group("/", middleware.JWTMiddleware(h.jwtService))
}

func (h *UserHandler) createUser(c *gin.Context) {
	var dto api_dto.PostUsersJSONBody

	if err := c.ShouldBindJSON(&dto); err != nil {
		HandleError(c, http.StatusBadRequest, err)
		return
	}

	userId, err := h.userUC.Create(c, dto)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	token, err := h.jwtService.GenerateToken(userId, time.Duration(auth.OneMonth))
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	auth.SetCookie(c, auth.JWTTokenCookie, token, auth.OneMonth)

	c.Status(http.StatusCreated)
}

func (h *UserHandler) login(c *gin.Context) {
	var dto api_dto.PostUsersLogInJSONRequestBody

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
		HandleError(c, http.StatusUnauthorized, nil)
		return
	}

	token, err := h.jwtService.GenerateToken(userId, time.Duration(auth.OneMonth))
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err)
		return
	}

	auth.SetCookie(c, auth.JWTTokenCookie, token, auth.OneMonth)

	c.Status(http.StatusOK)
}
