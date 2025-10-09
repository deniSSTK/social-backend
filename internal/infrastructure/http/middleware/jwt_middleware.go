package middleware

import (
	"net/http"
	"social-backend/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			auth.ClearCookie(c, auth.JWTTokenCookie)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(*auth.JWTClaims)
		c.Set("userId", claims.UserId)

		c.Next()
	}
}
