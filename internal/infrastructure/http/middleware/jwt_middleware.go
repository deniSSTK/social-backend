package middleware

import (
	"net/http"
	"social-backend/internal/infrastructure/auth"
	"social-backend/internal/infrastructure/auth/cookie"
	"social-backend/internal/infrastructure/http/context"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(string(cookie.JWTTokenCookie))
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			cookie.ClearCookie(c, cookie.JWTTokenCookie)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(*auth.JWTClaims)
		c.Set(context.ContextUserId, claims.UserId)

		c.Next()
	}
}
