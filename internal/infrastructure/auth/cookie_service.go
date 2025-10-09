package auth

import (
	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, cookieName CookieName, CookieValue string, duration int) {
	c.SetCookie(
		string(cookieName),
		CookieValue,
		duration,
		"/",
		"",
		false,
		true,
	)
}

func ClearCookie(c *gin.Context, cookieName CookieName) {
	c.SetCookie(string(cookieName), "", -1, "/", "", false, true)
}

type CookieName string

const (
	JWTTokenCookie CookieName = "jwt_token"
)
