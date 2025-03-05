package gin

import (
	"compass-backend/pkg/api/api_errors"
	"compass-backend/pkg/common/models"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

func SetUser(c *gin.Context, user *models.User) {
	c.Set("user", user)
}

func GetUser(c *gin.Context) (*models.User, error) {
	user, ok := c.Get("user")
	if !ok {
		return nil, api_errors.ErrorUnauthorized
	}
	return user.(*models.User), nil
}

func SetTeam(c *gin.Context, team *models.Team) {
	c.Set("team", team)
}

func GetTeam(c *gin.Context) (*models.Team, error) {
	team, ok := c.Get("team")
	if !ok {
		return nil, api_errors.ErrorTeamNotFound
	}
	return team.(*models.Team), nil
}

func SetFile(c *gin.Context, file multipart.File) {
	c.Set("file", file)
}

func GetFile(c *gin.Context) (multipart.File, error) {
	file, ok := c.Get("file")
	if !ok {
		return nil, api_errors.ErrorFileNotFound
	}

	return file.(multipart.File), nil
}

func SetFiles(c *gin.Context, files []*multipart.FileHeader) {
	c.Set("files", files)
}

func GetFiles(c *gin.Context) ([]*multipart.FileHeader, error) {
	files, ok := c.Get("files")
	if !ok {
		return nil, api_errors.ErrorFilesNotFound
	}

	return files.([]*multipart.FileHeader), nil
}
