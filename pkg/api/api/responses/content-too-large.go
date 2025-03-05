package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ContentTooLarge(c *gin.Context, message string) {
	c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": message})
}
