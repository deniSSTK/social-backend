package cookie

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, cookieName CookieName, CookieValue string, duration time.Duration) {
	isProd := os.Getenv("RUN_MODE") == "prod"

	if isProd {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	c.SetCookie(
		string(cookieName),
		CookieValue,
		int(duration),
		"/",
		"",
		isProd,
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
