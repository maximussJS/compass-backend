package gin

import (
	"compass-backend/pkg/api/api/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	idParam     = "id"
	teamIdParam = "teamId"
	tokenParam  = "token"
)

func GetUintIdParam(c *gin.Context) (uint, bool) {
	idStr := c.Param(idParam)
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		responses.BadRequest(c, fmt.Sprintf("Invalid %s parameter", idParam))
		return 0, false
	}

	return uint(id), true
}

func GetStringIdParam(c *gin.Context) (string, bool) {
	id := c.Param(idParam)

	if id == "" {
		responses.BadRequest(c, fmt.Sprintf("Invalid %s parameter", idParam))
		return "", false
	}

	return id, true
}

func GetStringTeamIdParam(c *gin.Context) (string, bool) {
	id := c.Param(teamIdParam)

	if id == "" {
		responses.BadRequest(c, fmt.Sprintf("Invalid %s parameter", teamIdParam))
		return "", false
	}

	return id, true
}

func GetTokenParam(c *gin.Context) (string, bool) {
	token := c.Param(tokenParam)
	if token == "" {
		responses.BadRequest(c, fmt.Sprintf("Invalid %s parameter", tokenParam))
		return "", false
	}

	return token, true
}
