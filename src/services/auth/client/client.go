package client

import (
	"net/http"

	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/common/config"
	"github.com/lumi4x/news/src/services/auth/client/models"
)

type Client struct {
	Client     *http.Client
	AuthHeader string
}

func NewClientWithLogin(usernameOrEmail, password string) (*Client, *api.Error) {
	c := Client{
		Client: api.NewClient(),
	}

	response, aErr := api.SendRequest[models.LoginRequest, models.LoginResponse](c.Client, c.AuthHeader, "POST", "auth/login", &models.LoginRequest{
		UsernameOrEmail: usernameOrEmail,
		Password:        password,
	})
	if aErr != nil {
		return nil, aErr
	}

	c.AuthHeader = "Bearer " + response.Token

	return &c, nil
}

func NewClientWithDiscordUser(discordUserId string) (*Client, *api.Error) {
	c := Client{
		Client: api.NewClient(),
	}

	response, aErr := api.SendRequest[models.DiscordLoginRequest, models.DiscordLoginResponse](c.Client, c.AuthHeader, "POST", "auth/discord_login", &models.DiscordLoginRequest{
		InternalLoginAuthSecret: config.InternalLoginAuthSecret,
		DiscordUserId:           discordUserId,
	})
	if aErr != nil {
		return nil, aErr
	}

	c.AuthHeader = "Bearer " + response.Token

	return &c, nil
}

func NewClientWithInternalLogin() (*Client, *api.Error) {
	c := Client{
		Client: api.NewClient(),
	}

	response, aErr := api.SendRequest[models.InternalLoginRequest, models.InternalLoginResponse](c.Client, c.AuthHeader, "POST", "auth/internal_login", &models.InternalLoginRequest{
		InternalLoginAuthSecret: config.InternalLoginAuthSecret,
	})
	if aErr != nil {
		return nil, aErr
	}

	c.AuthHeader = "Bearer " + response.Token

	return &c, nil
}
