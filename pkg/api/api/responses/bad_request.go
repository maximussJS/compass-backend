package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func BadRequestWithAbort(c *gin.Context, message string) {
	BadRequest(c, message)
	c.Abort()
}
