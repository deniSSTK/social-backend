package middleware

import (
	"social-backend/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		JWTMiddleware(authService.JwtService)(c)
		if c.IsAborted() {
			return
		}

		UserMiddleware(authService.UserService)(c)
		if c.IsAborted() {
			return
		}
	}
}
