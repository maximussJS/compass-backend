package responses

import (
	"compass-backend/pkg/api/api_errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var err = api_errors.ErrorInternalServerError

func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
