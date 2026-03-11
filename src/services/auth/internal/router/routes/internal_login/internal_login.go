package internal_login

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/common/config"
	"github.com/lumi4x/news/src/services/auth/client/models"
)

func InternalLogin(_ *api.Writer, r *http.Request, _ httprouter.Params, _ *api.JWT, request models.InternalLoginRequest) (*models.InternalLoginResponse, *api.Error) {
	if request.InternalLoginAuthSecret != config.InternalLoginAuthSecret {
		return nil, api.NewUnauthorizedError(nil, "bad internal auth secret")
	}

	token, err := api.MarshalJWT("", api.Internal.AccessLevel, map[string]bool{})
	if err != nil {
		return nil, api.NewInternalError(err, "failed creating JWT token")
	}

	return &models.InternalLoginResponse{Token: token}, nil
}
