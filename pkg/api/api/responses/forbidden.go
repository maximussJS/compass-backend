package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
}

func ForbiddenWithAbort(c *gin.Context) {
	Forbidden(c)
	c.Abort()
}
