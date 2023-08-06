package imgbb

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type Option func(*ImgRequest)

func WithEndpoint(endpoint string) Option {
	return func(reqFile *ImgRequest) {
		reqFile.endpoint = endpoint
	}
}

// NewImage creates a new Image
func NewImage(name string, expiration string, file []byte) *Image {
	return &Image{
		Name:       name,
		Size:       len(file),
		Expiration: expiration,
		File:       file,
	}
}

// FileUploadedPart ...
func FileUploadedPart(
	pipeReader *io.PipeReader,
	pipeWriter *io.PipeWriter,
	reqWriter *multipart.Writer,
	req *ImgRequest,
	img *Image,
) error {
	defer pipeWriter.Close()
	defer reqWriter.Close()

	err := reqWriter.WriteField("key", req.key)
	if err != nil {
		return err
	}

	err = reqWriter.WriteField("type", "file")
	if err != nil {
		return err
	}

	err = reqWriter.WriteField("action", "upload")
	if err != nil {
		return err
	}

	if len(img.Expiration) > 0 {
		err = reqWriter.WriteField("expiration", img.Expiration)
		if err != nil {
			return err
		}
	}

	part, err := reqWriter.CreateFormFile("image", img.Name)
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, bytes.NewReader(img.File)); err != nil {
		return err
	}

	return nil
}

// ResponseParse ...
func ResponseParse(resp *http.Response) (*HttpResponse, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrResp{
			StatusCode: http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError),
			ErrMsg:     errors.New(fmt.Sprintf("read response body: %v", err)),
		}
	}

	if resp.StatusCode == http.StatusOK {
		var res HttpResponse
		if err := json.Unmarshal(data, &res); err != nil {
			return nil, ErrResp{
				StatusCode: http.StatusInternalServerError,
				StatusText: http.StatusText(http.StatusInternalServerError),
				ErrMsg:     errors.New(fmt.Sprintf("json unmarshal: %v", err)),
			}
		}

		return &res, nil
	}

	var errRes ErrResp
	if err := json.Unmarshal(data, &errRes); err != nil {
		return nil, ErrResp{
			StatusCode: http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError),
			ErrMsg:     errors.New(fmt.Sprintf("json unmarshal: %v", err)),
		}
	}

	return nil, errRes
}

// HashSum ...
func HashSum(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}
