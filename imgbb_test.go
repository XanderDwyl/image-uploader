package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	imgbb "github.com/XanderDwyl/image-uploader/imgbb"
	"github.com/stretchr/testify/assert"
)

var testImg = []byte{120, 97, 110, 100, 101, 114, 100, 119, 121, 108, 46, 80, 78, 71}

func Test_SuccessUpload(t *testing.T) {
	resp := `{
		"data": {
			"id": "2ndCYJK",
			"title": "c1f64245afb2",
			"url_viewer": "https://ibb.co/2ndCYJK",
			"url": "https://i.ibb.co/w04Prt6/c1f64245afb2.gif",
			"display_url": "https://i.ibb.co/98W13PY/c1f64245afb2.gif",
			"width": 1,
			"height": 1,
			"size": 42,
			"time": "1552042565",
			"expiration":"0",
			"image": {
				"filename": "c1f64245afb2.gif",
				"name": "c1f64245afb2",
				"mime": "image/gif",
				"extension": "gif",
				"url": "https://i.ibb.co/w04Prt6/c1f64245afb2.gif"
			},
			"thumb": {
				"filename": "c1f64245afb2.gif",
				"name": "c1f64245afb2",
				"mime": "image/gif",
				"extension": "gif",
				"url": "https://i.ibb.co/2ndCYJK/c1f64245afb2.gif"
			},
			"medium": {
				"filename": "c1f64245afb2.gif",
				"name": "c1f64245afb2",
				"mime": "image/gif",
				"extension": "gif",
				"url": "https://i.ibb.co/98W13PY/c1f64245afb2.gif"
			},
			"delete_url": "https://ibb.co/2ndCYJK/670a7e48ddcb85ac340c717a41047e5c"
		},
		"success": true,
		"status": 200
	}`

	img := imgbb.NewImage("name", "", testImg)

	ts := httptest.NewTLSServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				fmt.Fprintln(w, resp)
			},
		),
	)
	defer ts.Close()

	apiClient := imgbb.New(*ts.Client(), "secret-key", imgbb.WithEndpoint(ts.URL))

	expect := &imgbb.HttpResponse{
		Data: imgbb.FileData{
			ID:         "2ndCYJK",
			Title:      "c1f64245afb2",
			UrlViewer:  "https://ibb.co/2ndCYJK",
			Url:        "https://i.ibb.co/w04Prt6/c1f64245afb2.gif",
			DisplayUrl: "https://i.ibb.co/98W13PY/c1f64245afb2.gif",
			Width:      1,
			Height:     1,
			Size:       42,
			Time:       "1552042565",
			Expiration: "0",
			Image: imgbb.FileInfo{
				Filename:  "c1f64245afb2.gif",
				Name:      "c1f64245afb2",
				Mime:      "image/gif",
				Extension: "gif",
				Url:       "https://i.ibb.co/w04Prt6/c1f64245afb2.gif",
			},
			Thumb: imgbb.FileInfo{
				Filename:  "c1f64245afb2.gif",
				Name:      "c1f64245afb2",
				Mime:      "image/gif",
				Extension: "gif",
				Url:       "https://i.ibb.co/2ndCYJK/c1f64245afb2.gif",
			},
			Medium: imgbb.FileInfo{
				Filename:  "c1f64245afb2.gif",
				Name:      "c1f64245afb2",
				Mime:      "image/gif",
				Extension: "gif",
				Url:       "https://i.ibb.co/98W13PY/c1f64245afb2.gif",
			},
			DeleteUrl: "https://ibb.co/2ndCYJK/670a7e48ddcb85ac340c717a41047e5c",
		},
		Success:    true,
		StatusCode: http.StatusOK,
	}

	actual, err := apiClient.Upload(img)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func Test_OverSizeImage(t *testing.T) {
	img := &imgbb.Image{
		Name:       "name",
		Size:       len(testImg) * 10000000,
		Expiration: "",
		File:       testImg,
	}

	apiClient := imgbb.New(http.Client{}, "secret-key")

	_, err := apiClient.Upload(img)

	assert.ErrorIs(t, err, imgbb.ErrResp{
		StatusCode: http.StatusBadRequest,
		StatusText: http.StatusText(http.StatusBadRequest),
	})
}
