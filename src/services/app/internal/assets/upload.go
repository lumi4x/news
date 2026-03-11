package assets

import (
	"embed"
	"fmt"
	"path/filepath"

	auth_client "github.com/lumi4x/news/src/services/auth/client"
	"github.com/lumi4x/news/src/services/fileserver/client"
	"github.com/sirupsen/logrus"
)

//go:embed public/*
var publicAssets embed.FS

var uploadedPublicAssets = map[string]*client.UploadResponse{}

func GetPublicAsset(filename string) (*client.UploadResponse, error) {
	asset, ok := uploadedPublicAssets[filename]
	if !ok {
		return nil, fmt.Errorf("failed getting public asset, filename likely wrong")
	}

	return asset, nil
}

func init() {
	authClient, aErr := auth_client.NewClientWithInternalLogin()
	if aErr != nil {
		logrus.WithError(aErr.Err).WithField("code", aErr.Code).Panicf("failed creating new auth client for app")
	}

	client := client.NewClient(authClient)

	assetFiles, err := publicAssets.ReadDir("public")
	if err != nil {
		logrus.WithError(err).Panicf("failed reading assets dir")
	}

	for _, assetFile := range assetFiles {
		if assetFile.IsDir() {
			logrus.Panicf("found directory in assets folder")
		}

		fileName := assetFile.Name()
		filePath := filepath.Join("public", fileName)

		fileContent, err := publicAssets.ReadFile(filePath)
		if err != nil {
			logrus.WithError(err).WithField("filePath", filePath).Panicf("Failed to read file")
		}

		logrus.WithField("fileName", fileName).Info("Uploading asset")
		uploadResponse, aErr := client.UploadPublic(fileContent)
		if aErr != nil {
			logrus.WithError(aErr.Err).Panicf("failed uploading rya avatar")
		}

		uploadedPublicAssets[fileName] = uploadResponse
	}
}
