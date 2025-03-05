package middlewares

import (
	"compass-backend/pkg/api/api/responses"
	"compass-backend/pkg/api/lib"
	gin_utils "compass-backend/pkg/api/utils/gin"
	fx_utils "compass-backend/pkg/common/fx"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"mime/multipart"
)

type IUploadMiddleware interface {
	Handle() gin.HandlerFunc
}

type uploadMiddlewareParams struct {
	fx.In

	Env lib.IEnv
}

type uploadMiddleware struct {
	maxFileSize int64
}

func FxUploadMiddleware() fx.Option {
	return fx_utils.AsProvider(newUploadMiddleware, new(IUploadMiddleware))
}

func newUploadMiddleware(params uploadMiddlewareParams) IUploadMiddleware {
	return &uploadMiddleware{
		maxFileSize: params.Env.GetMaxFileSize(),
	}
}

func (m uploadMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Request.ParseMultipartForm(m.maxFileSize)
		if err != nil {
			responses.BadRequest(c, "Unable to parse multipart form")
			return
		}

		form := c.Request.MultipartForm
		if form == nil {
			c.Next()
			return
		}

		files := form.File["file"]
		if files == nil {
			c.Next()
			return
		}

		var uploadedFiles []*multipart.FileHeader
		for _, fileHeader := range files {
			uploadedFiles = append(uploadedFiles, fileHeader)
		}

		gin_utils.SetFiles(c, uploadedFiles)

		c.Next()
	}
}
