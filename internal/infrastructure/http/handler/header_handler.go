package handler

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func HandleETag(c *gin.Context, obj interface{}) {
	jsonBytes, _ := json.Marshal(obj)

	hash := md5.Sum(jsonBytes)
	etag := hex.EncodeToString(hash[:])

	c.Header("Access-Control-Expose-Headers", "ETag")
	c.Header("ETag", etag)

	if match := c.GetHeader("If-None-Match"); match != "" && match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	c.Data(http.StatusOK, "application/json", jsonBytes)
}
