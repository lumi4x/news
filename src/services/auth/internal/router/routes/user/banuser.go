package user

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/services/auth/client"
	"github.com/lumi4x/news/src/services/auth/client/models"
	"github.com/lumi4x/news/src/services/auth/internal/ent"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build"
	"github.com/lumi4x/news/src/services/auth/internal/ent/build/user"
	"github.com/lumi4x/news/src/services/auth/internal/ent/mixins"
)

func BanUser(_ *api.Writer, r *http.Request, params httprouter.Params, _ *api.JWT, _ models.BanUserRequest) (*models.BanUserResponse, *api.Error) {
	userUUID := params.ByName("uuid")
	parsedUserUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return nil, api.NewBadRequestError(err, "bad uuid provided to ban user")
	}

	u, err := ent.Database.User.Query().Where(user.UUIDEQ(parsedUserUUID)).Only(mixins.SkipBanFilter(r.Context()))
	if err != nil {
		if _, ok := err.(*build.NotFoundError); ok {
			return nil, api.NewBadRequestError(err, "user not found").WithCode(client.UserNotFound)
		}
		return nil, api.NewInternalError(err, "failed querying db for banning user")
	}

	if !u.BanTime.IsZero() {
		return nil, api.NewBadRequestError(nil, "user is already banned").WithCode(client.UserAlreadyBanned)
	}

	err = ent.Database.User.Update().SetBanTime(time.Now()).Where(user.UUIDEQ(parsedUserUUID)).Exec(r.Context())
	if err != nil {
		return nil, api.NewInternalError(err, "failed querying db for ban")
	}

	return &models.BanUserResponse{
		BannedUserDiscordId: u.DiscordUserID,
	}, nil
}

func UnbanUser(_ *api.Writer, r *http.Request, params httprouter.Params, _ *api.JWT, _ models.UnbanUserRequest) (*models.UnbanUserResponse, *api.Error) {
	userUUID := params.ByName("uuid")
	parsedUserUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return nil, api.NewBadRequestError(err, "bad uuid provided to unban user")
	}

	u, err := ent.Database.User.Query().Where(user.UUIDEQ(parsedUserUUID)).Only(mixins.SkipBanFilter(r.Context()))
	if err != nil {
		if _, ok := err.(*build.NotFoundError); ok {
			return nil, api.NewBadRequestError(err, "user not found").WithCode(client.UserNotFound)
		}
		return nil, api.NewInternalError(err, "failed querying db for banning user")
	}

	if u.BanTime.IsZero() {
		return nil, api.NewBadRequestError(nil, "user already unbanned").WithCode(client.UserAlreadyUnbanned)
	}

	err = ent.Database.User.Update().ClearBanTime().Where(user.UUID(parsedUserUUID)).Exec(mixins.SkipBanFilter(r.Context()))
	if err != nil {
		return nil, api.NewInternalError(err, "failed querying db for unban")
	}

	return &models.UnbanUserResponse{
		UserDiscordId: u.DiscordUserID,
	}, nil
}
