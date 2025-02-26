package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuccessEmpty(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func SuccessJson(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func SuccessMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func SuccessJsonWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{"data": data, "message": message})
}
