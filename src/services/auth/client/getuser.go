package client

import (
	"fmt"

	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/services/auth/client/models"
)

func (c *Client) GetUser(userUUID string) (*models.GetUserResponse, *api.Error) {
	return api.SendRequest[any, models.GetUserResponse](c.Client, c.AuthHeader, "GET", fmt.Sprintf("auth/user/adminget/%s", userUUID), nil)
}
