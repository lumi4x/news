package client

import (
	"fmt"

	"github.com/lumi4x/news/src/common/api"
	common_config "github.com/lumi4x/news/src/common/config"
	"github.com/lumi4x/news/src/services/auth/client"
)

type Client struct {
	*client.Client
}

func NewClient(client *client.Client) *Client {
	return &Client{
		Client: client,
	}
}

type UploadResponse struct {
	StorageKey string
	URL        string
}

func GetURL(storageKey string) string {
	return fmt.Sprintf("%s/api/v1/fileserver/file/%s", common_config.SystemExternalURL, storageKey)
}

func (c *Client) Download(filename string) ([]byte, *api.Error) {
	return api.SendDownloadRequest[any](c.Client.Client, c.AuthHeader, "GET", fmt.Sprintf("fileserver/file/%s", filename), nil)
}

func (c *Client) Upload(data []byte) (*UploadResponse, *api.Error) {
	return api.SendMultipartFile[UploadResponse](c.Client.Client, c.AuthHeader, "POST", "fileserver/file/default", data)
}

func (c *Client) UploadPublic(data []byte) (*UploadResponse, *api.Error) {
	return api.SendMultipartFile[UploadResponse](c.Client.Client, c.AuthHeader, "POST", "fileserver/file/public", data)
}
