package handler

import (
	"net/http"
	"social-backend/internal/infrastructure/http/api_dto"
	"social-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUC *usecase.UserUsecase
}

func NewUserHandler(userUC *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUC}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	//TODO auth.NewJWTService()
	//TODO protected := router.Group("/users", middleware.JWTMiddleware(authService))

	group := router.Group("/user")

	group.POST("/", h.createUser)
}

func (h *UserHandler) createUser(c *gin.Context) {
	var dto api_dto.PostUsersJSONBody
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userUC.Create(c, dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
