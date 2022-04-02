package utilits

import (
	"bytes"
	"context"
	"fmt"
	"github.com/conku/webp"
	"github.com/pkg/errors"
	"glide/internal/app"
	repoFiles "glide/internal/microservices/files/files/repository/files"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"
)

//go:generate mockgen -destination=mocks/mock_convert_utils.go -package=mock_utils . ImageConverter

type ImageConverter interface {
	// Convert Errors:
	// 	app.GeneralError:
	//		utils.ConvertErr
	// 		utils.UnknownExtOfFileName
	Convert(context.Context, io.Reader, repoFiles.FileName) (io.Reader, repoFiles.FileName, error)
}

type ConverterToWebp struct {
}

// Convert Errors:
// 	app.GeneralError:
//		utils.ConvertErr
// 		utils.UnknownExtOfFileName
func (cv *ConverterToWebp) Convert(_ context.Context,
	file io.Reader, name repoFiles.FileName) (io.Reader, repoFiles.FileName, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, "", app.GeneralError{
			Err:         errors.Wrap(ConvertErr, "error in get image"),
			ExternalErr: err,
		}
	}

	buf, err := webp.EncodeExactLosslessRGBA(img)
	if err != nil {
		return nil, "", app.GeneralError{
			Err:         errors.Wrap(ConvertErr, "error in webp convertor"),
			ExternalErr: err,
		}
	}

	pos := strings.LastIndex(string(name), ".")
	if pos == -1 {
		return nil, "", errors.Wrap(UnknownExtOfFileName, fmt.Sprintf("error with %s: ", name))
	}

	name = name[:pos] + ".webp"

	res := bytes.NewReader(buf)
	return res, name, nil
}
