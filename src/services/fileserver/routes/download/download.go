package download

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/services/fileserver/config"
	"github.com/lumi4x/news/src/services/fileserver/utils/filestorage"
)

var FilenameRegex = regexp.MustCompile("^[A-Za-z0-9]+")

func Download(w *api.Writer, r *http.Request, params httprouter.Params, jwt *api.JWT) (*any, *api.Error) {
	return download(w, r, params, jwt, config.DefaultDirectoryName)
}

func DownloadPublic(w *api.Writer, r *http.Request, params httprouter.Params, jwt *api.JWT) (*any, *api.Error) {
	return download(w, r, params, jwt, config.PublicDirectoryName)
}

func download(w *api.Writer, r *http.Request, params httprouter.Params, jwt *api.JWT, dir string) (*any, *api.Error) {
	filename := params.ByName("filename")
	if filename == "" {
		return nil, api.NewBadRequestError(nil, "missing filename")
	}

	// Remove the extension if provided (this is meant for discord support)
	filename, _, _ = strings.Cut(filename, ".")

	if !FilenameRegex.MatchString(filename) {
		return nil, api.NewBadRequestError(nil, "bad filename")
	}

	found, err := filestorage.Default.Download(fmt.Sprintf("%s/%s", dir, filename), w, r.Context())
	if err != nil {
		return nil, api.NewInternalError(err, "failed downloading file")
	}

	if !found {
		return nil, api.NewNotFoundError(err, "file does not exist in filestorage")
	}

	return nil, nil
}
