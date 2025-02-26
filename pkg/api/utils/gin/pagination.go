package gin

import (
	"compass-backend/pkg/api/api/responses"
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getIntQueryParam(c *gin.Context, param string, defaultValue, maxValue int) int {
	valueStr := c.Query(param)

	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil || value < 0 {
		responses.BadRequestWithAbort(c, fmt.Sprintf("Invalid %s query parameter. Should be a positive integer", param))
		return -1 // This value will be ignored by the caller
	}

	if value > maxValue {
		responses.BadRequestWithAbort(c, fmt.Sprintf("Invalid %s query parameter. Should be less than %d", param, maxValue))
		return -1 // This value will be ignored by the caller
	}

	return value
}

func GetLimit(c *gin.Context, defaultLimit, maxLimit int) int {
	return getIntQueryParam(c, "limit", defaultLimit, maxLimit)
}

func GetOffset(c *gin.Context) int {
	return getIntQueryParam(c, "offset", 0, math.MaxInt)
}
