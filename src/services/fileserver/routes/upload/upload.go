package upload

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lumi4x/news/src/common/api"
	"github.com/lumi4x/news/src/services/fileserver/client"
	"github.com/lumi4x/news/src/services/fileserver/config"
	"github.com/lumi4x/news/src/services/fileserver/utils/filestorage"
)

func Upload(w *api.Writer, r *http.Request, _ httprouter.Params, jwt *api.JWT, _ any) (*client.UploadResponse, *api.Error) {
	return upload(w, r, jwt, config.DefaultDirectoryName)
}

func UploadPublic(w *api.Writer, r *http.Request, _ httprouter.Params, jwt *api.JWT, _ any) (*client.UploadResponse, *api.Error) {
	return upload(w, r, jwt, config.PublicDirectoryName)
}

func upload(w *api.Writer, r *http.Request, jwt *api.JWT, dir string) (*client.UploadResponse, *api.Error) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(15 << 20) // limit the file size to 15MB, if changing also change in nginx
	if err != nil {
		return nil, api.NewBadRequestError(err, "limited to 15 mb of file upload")
	}

	// Get the file from the request
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, api.NewBadRequestError(err, "failed getting form file")
	}
	defer file.Close()

	storageKey, err := filestorage.Default.Upload(dir, file)
	if err != nil {
		return nil, api.NewInternalError(err, "failed uploading file to file storage")
	}

	url := client.GetURL(storageKey)

	return &client.UploadResponse{StorageKey: storageKey, URL: url}, nil
}
