package imgbb

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// ImgRequest ...
type ImgRequest struct {
	client   http.Client
	key      string
	endpoint string
}

// HttpResponse ...
type HttpResponse struct {
	Data       FileData `json:"data"`
	StatusCode int      `json:"status"`
	Success    bool     `json:"success"`
}

// Image ...
type Image struct {
	Name       string `json:"name"`
	Size       int    `json:"size"`
	Expiration string `json:"expiration"`
	File       []byte `json:"file"`
}

// FileData ...
type FileData struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	UrlViewer  string   `json:"url_viewer"`
	Url        string   `json:"url"`
	DisplayUrl string   `json:"display_url"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Size       int      `json:"size"`
	Time       string   `json:"time"`
	Expiration string   `json:"expiration"`
	Image      FileInfo `json:"image"`
	Thumb      FileInfo `json:"thumb"`
	Medium     FileInfo `json:"medium"`
	DeleteUrl  string   `json:"delete_url"`
}

// FileInfo ...
type FileInfo struct {
	Filename  string `json:"filename"`
	Name      string `json:"name"`
	Mime      string `json:"mime"`
	Extension string `json:"extension"`
	Url       string `json:"url"`
}

// New ...
func New(newClient http.Client, newKey string, newOpts ...Option) *ImgRequest {
	req := &ImgRequest{
		client:   newClient,
		key:      newKey,
		endpoint: "https://api.imgbb.com/1/upload",
	}

	for _, opt := range newOpts {
		opt(req)
	}

	return req
}

func (imgReq *ImgRequest) Upload(img *Image) (*HttpResponse, error) {

	switch size := img.Size; {
	case size <= 0:
		return nil, ErrResp{
			StatusCode: http.StatusBadRequest,
			StatusText: http.StatusText(http.StatusBadRequest),
			ErrMsg:     errors.New("image file is empty"),
		}
	case size > 33554432:
		return nil, ErrResp{
			StatusCode: http.StatusBadRequest,
			StatusText: http.StatusText(http.StatusBadRequest),
			ErrMsg:     errors.New("image is more than 32mb"),
		}
	default:
		ioReader, ioWriter := io.Pipe()
		mWriter := multipart.NewWriter(ioWriter)

		go FileUploadedPart(ioReader, ioWriter, mWriter, imgReq, img)

		req, err := http.NewRequest(http.MethodPost, imgReq.endpoint, ioReader)
		if err != nil {
			return nil, ErrResp{
				StatusCode: http.StatusInternalServerError,
				StatusText: http.StatusText(http.StatusInternalServerError),
				ErrMsg:     errors.New(fmt.Sprintf("new request: %v", err)),
			}
		}

		req.Header.Add("Content-Type", mWriter.FormDataContentType())
		req.Header.Add("Host", "imgbb.com")
		req.Header.Add("Origin", "https://imgbb.com")
		req.Header.Add("Referer", "https://imgbb.com")

		resp, err := imgReq.client.Do(req)
		if err != nil {
			return nil, ErrResp{
				StatusCode: http.StatusInternalServerError,
				StatusText: http.StatusText(http.StatusInternalServerError),
				ErrMsg:     errors.New(fmt.Sprintf("http client request do: %v", err)),
			}
		}
		defer resp.Body.Close()

		return ResponseParse(resp)
	}

}
