package middleware

import (
	"net/http"
	"social-backend/internal/infrastructure/auth"
	"social-backend/internal/infrastructure/auth/cookie"
	"social-backend/internal/infrastructure/errors"
	"social-backend/internal/infrastructure/http/context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserMiddleware(userService *auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := context.GetContextUserId(c)
		if userId == uuid.Nil {
			clearAndAbort(c, http.StatusUnauthorized, gin.H{"error": errors.ContextUserIdEmpty.Error()})
			return
		}

		exists, err := userService.UserRepo.CheckIfExistsById(c.Request.Context(), userId)
		if err != nil {
			clearAndAbort(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !exists {
			clearAndAbort(c, http.StatusUnauthorized, gin.H{"error": errors.UserIdDoesNotExists.Error()})
			return
		}

		c.Next()
	}
}

func clearAndAbort(c *gin.Context, status int, msg any) {
	cookie.ClearCookie(c, cookie.JWTTokenCookie)
	if msg != nil {
		c.AbortWithStatusJSON(status, gin.H{"error": msg})
	} else {
		c.AbortWithStatus(status)
	}
}
