package login

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/services/auth/client"
	"github.com/lumi4x/news/src/services/auth/client/models"
	"github.com/lumi4x/news/src/services/auth/internal/config"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build"
	"github.com/lumi4x/news/src/services/auth/internal/modules/role"
	usermodule "github.com/lumi4x/news/src/services/auth/internal/modules/user"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
)

func Login(_ *api.Writer, r *http.Request, _ httprouter.Params, _ *api.JWT, request models.LoginRequest) (*models.LoginResponse, *api.Error) {
	// Minimal brute force protection
	time.Sleep(1 * time.Second)

	// Determine if provided username or email for later querying
	logrus.WithField("request.UsernameOrEmail", request.UsernameOrEmail).Info("Login request")

	user, err := usermodule.GetIncludeBanned(r.Context(), request.UsernameOrEmail)
	if err != nil {
		if _, ok := err.(*build.NotFoundError); ok {
			return nil, api.NewUnauthorizedError(err, "user not found")
		}
		return nil, api.NewInternalError(err, "failed querying db for login")
	}

	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return nil, api.NewInternalError(err, "failed decoding salt")
	}

	// Hash the password using PBKDF2
	hashedPassword := pbkdf2.Key([]byte(request.Password), salt, config.Iterations, config.PasswordHashLength, sha256.New)
	hashedPasswordString := base64.StdEncoding.EncodeToString(hashedPassword)

	if user.Password != hashedPasswordString {
		return nil, api.NewUnauthorizedError(nil, "password is incorrect")
	}

	if !user.BanTime.IsZero() {
		return nil, api.NewUnauthorizedError(nil, "user is banned").WithCode(client.YouAreBanned)
	}

	permissions, err := role.GetPermissions(r.Context(), user.ID)
	if err != nil {
		return nil, api.NewInternalError(err, "failed getting permissions during user login")
	}

	token, err := api.MarshalJWT(user.ID.String(), api.User.AccessLevel, permissions)
	if err != nil {
		return nil, api.NewInternalError(err, "failed creating JWT token")
	}

	// Return the signed JWT
	return &models.LoginResponse{Token: token}, nil
}
