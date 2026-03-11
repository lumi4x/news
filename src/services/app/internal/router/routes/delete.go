package routes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/services/app/client/models"
	"github.com/lumi4x/news/src/services/app/internal/modules/newsitem"
	"github.com/lumi4x/news/src/services/auth/client"
)

func Delete(w *api.Writer, r *http.Request, p httprouter.Params, j *api.JWT, request models.DeleteNewsItemRequest) (*any, *api.Error) {
	if _, ok := j.Permissions[client.KeyNewsItemDelete]; !ok {
		return nil, api.NewForbiddenError(nil, "no permission to toggle blur")
	}

	err := newsitem.SoftDelete(r.Context(), uuid.MustParse(request.RssItemId))
	if err != nil {
		return nil, api.NewInternalError(err, "failed toggling blur")
	}

	return nil, nil
}
