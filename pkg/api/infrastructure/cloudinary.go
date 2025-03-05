package infrastructure

import (
	"compass-backend/pkg/api/lib"
	fx_utils "compass-backend/pkg/common/fx"
	common_lib "compass-backend/pkg/common/lib"
	"context"
	"fmt"
	go_cld "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"go.uber.org/fx"
	"mime/multipart"
)

type ICloudinary interface {
	UploadImage(ctx context.Context, file multipart.File, filePath string) (publicId, secureUrl string, err error)
}

type cloudinaryParams struct {
	fx.In

	Logger common_lib.ILogger
	Env    lib.IEnv
}

type cloudinary struct {
	cld *go_cld.Cloudinary
}

func FxCloudinary() fx.Option {
	return fx_utils.AsProvider(newCloudinary, new(ICloudinary))
}

func newCloudinary(params cloudinaryParams) ICloudinary {
	cld, err := go_cld.NewFromURL(params.Env.GetCloudinaryUrl())

	if err != nil {
		params.Logger.Error(fmt.Sprintf("Failed to create cloudinary client: %s", err))
		panic(err)
	}

	cld.Config.URL.Secure = true

	return &cloudinary{
		cld: cld,
	}
}

func (s *cloudinary) UploadImage(ctx context.Context, file multipart.File, filePath string) (publicId, secureUrl string, err error) {
	uploadParams := uploader.UploadParams{
		PublicID: filePath,
	}

	result, err := s.cld.Upload.Upload(ctx, file, uploadParams)

	if err != nil {
		return "", "", fmt.Errorf("failed to upload image: %s", err)
	}

	return result.PublicID, result.SecureURL, nil
}
