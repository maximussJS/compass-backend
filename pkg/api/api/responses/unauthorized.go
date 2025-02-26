package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
}

func UnauthorizedWithAbort(c *gin.Context) {
	Unauthorized(c)
	c.Abort()
}
