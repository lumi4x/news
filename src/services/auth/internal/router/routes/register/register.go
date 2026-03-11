package register

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/services/auth/client/models"
	"github.com/lumi4x/news/src/services/auth/internal/config"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build"
	"github.com/lumi4x/news/src/services/auth/internal/modules/role"
	"github.com/lumi4x/news/src/services/auth/internal/modules/user"
)

func Register(_ *api.Writer, r *http.Request, _ httprouter.Params, _ *api.JWT, request models.RegisterRequest) (*models.RegisterResponse, *api.Error) {
	// Minimal brute force protection
	time.Sleep(1 * time.Second)

	// Validate password as it cannot be validated with Ent because it is not saved plaintext
	if len(request.Password) < config.PasswordMinLength {
		return nil, api.NewBadRequestError(nil, "password too short")
	}

	newUser, err := user.Register(r.Context(), request.Email, request.Password)
	if err != nil {
		if _, ok := err.(*build.ConstraintError); ok {
			return nil, api.NewBadRequestError(err, "user might already exist")
		}
		if _, ok := err.(*build.ValidationError); ok {
			return nil, api.NewBadRequestError(err, "invalid registration input")
		}
		return nil, api.NewInternalError(err, "failed creating user in database")
	}

	permissions, err := role.GetPermissions(r.Context(), newUser.ID)
	if err != nil {
		return nil, api.NewInternalError(err, "failed getting permissions when registering user")
	}

	token, err := api.MarshalJWT(newUser.ID.String(), api.User.AccessLevel, permissions)
	if err != nil {
		return nil, api.NewInternalError(err, "failed creating JWT token")
	}

	return &models.RegisterResponse{Token: token}, nil
}
