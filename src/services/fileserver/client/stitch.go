package client

import "github.com/lumi4x/news/src/common/api"

type StitchRequest struct {
	StorageKeys []string
}

type StitchResponse struct {
	StorageKey string
	URL        string
}

func (c *Client) Stitch(storagePaths []string) (*StitchResponse, *api.Error) {
	return api.SendRequest[StitchRequest, StitchResponse](c.Client.Client, c.AuthHeader, "POST", "fileserver/stitch", &StitchRequest{StorageKeys: storagePaths})
}
