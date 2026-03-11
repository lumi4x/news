package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/services/app/client/models"
	"github.com/lumi4x/news/src/services/app/internal/modules/organization"
)

func GetOrganizations(w *api.Writer, r *http.Request, p httprouter.Params, j *api.JWT) (*models.GetOrganizationsResponse, *api.Error) {
	orgs, err := organization.GetAll(r.Context())
	if err != nil {
		return nil, api.NewInternalError(err, "failed getting organizations")
	}

	return &models.GetOrganizationsResponse{
		Organizations: organization.TranslateOrganizations(orgs),
	}, nil
}
