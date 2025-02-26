package gin

import (
	"compass-backend/pkg/api/api_errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func BindData(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBind(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			var invalidArgs []invalidArgument

			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":  api_errors.ErrorInvalidRequestParameter,
				"fields": invalidArgs,
			})
			return false
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": api_errors.ErrorInternalServerError})
		return false
	}

	return true
}
