package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"error": message})
}

func NotFoundWithAbort(c *gin.Context, message string) {
	NotFound(c, message)
	c.Abort()
}
