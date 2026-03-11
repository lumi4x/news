package stitch

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/lumi4x/news/src/common/api"
	common_config "github.com/lumi4x/news/src/common/config"
	"github.com/lumi4x/news/src/services/fileserver/client"
	"github.com/lumi4x/news/src/services/fileserver/config"
	"github.com/lumi4x/news/src/services/fileserver/utils/filestorage"
	"github.com/sirupsen/logrus"
)

func Stitch(w *api.Writer, r *http.Request, _ httprouter.Params, jwt *api.JWT, request client.StitchRequest) (*client.StitchResponse, *api.Error) {
	var images []image.Image

	logrus.WithField("request.StorageKeys", request.StorageKeys).Info("Stitching")
	for _, storageKey := range request.StorageKeys {
		file, err := os.CreateTemp("", "imagetostitch-*.png")
		if err != nil {
			return nil, api.NewInternalError(err, "failed creating temp file")
		}
		defer os.Remove(file.Name())

		exists, err := filestorage.Default.Download(storageKey, file, r.Context())
		if err != nil {
			return nil, api.NewInternalError(err, "failed downloading file")
		}

		if !exists {
			return nil, api.NewBadRequestError(err, "file not found for stitching")
		}

		img, _, err := image.Decode(file)
		if err != nil {
			return nil, api.NewInternalError(err, "failed decoding image")
		}

		images = append(images, img)
	}

	// Calculate the dimensions of the stitched image
	padding := 10
	var width, height int
	for _, img := range images {
		width += img.Bounds().Dx() + padding
		if img.Bounds().Dy() > height {
			height = img.Bounds().Dy()
		}
	}
	width -= padding // Remove padding from the last image

	// Create a new image with the desired dimensions and fill it with transparent pixels
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	// Copy each image onto the new image with padding
	offset := 0
	for _, img := range images {
		logrus.WithField("offset", offset).Info("Image")
		draw.Draw(newImg, image.Rect(offset+padding, 0, offset+img.Bounds().Dx()+padding, img.Bounds().Dy()), img, image.Point{}, draw.Src)
		offset += img.Bounds().Dx() + padding
	}

	file, err := os.CreateTemp("", "stitchedimage-*.png")
	if err != nil {
		return nil, api.NewInternalError(err, "failed creating temp file for stitch")
	}
	defer os.Remove(file.Name())

	err = png.Encode(file, newImg)
	if err != nil {
		return nil, api.NewInternalError(err, "failed encoding stitched image")
	}

	storageKey, err := filestorage.Default.Upload(config.PublicDirectoryName, file)
	if err != nil {
		return nil, api.NewInternalError(err, "failed uploading stitched image")
	}

	url := fmt.Sprintf("%s/api/v1/fileserver/file/%s", common_config.SystemExternalURL, storageKey)

	return &client.StitchResponse{StorageKey: storageKey, URL: url}, nil
}
