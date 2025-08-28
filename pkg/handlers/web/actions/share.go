package actions

import (
	"homeserver/constants"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/pkg/crypto"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/utils/formatters"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ShareFileHandler(w http.ResponseWriter, r *http.Request, repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext) {
	fileID := r.URL.Query().Get("id")
	isFolder := r.URL.Query().Get("is_folder") == "true"

	file_srv_host := os.Getenv(constants.FileSrvHost)
	file_srv_port := os.Getenv(constants.FileSrvPort)

	url := "http://" + file_srv_host + ":" + file_srv_port

	path := ""
	if isFolder {
		path += "/folders/" + fileID
	} else {
		path += "/files/" + fileID
	}

	exp := strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10)
	data := path + "|" + exp

	url_hash := crypto.HashSeed(data)

	query := "?sid=" + url_hash + "&exp=" + exp

	entry_id, err := strconv.Atoi(fileID)
	if err != nil {
		ErrorResponse(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	if isFolder {
		fvp := folschema.NewUserFolderBuilder().WithTrustedSource([]byte(url_hash))
		folder := folschema.NewUserFolderBuilder().WithID(entry_id).WithOwnerID(reqCont.Requester.ID)

		repos.FolderRepo.UpdateFolder(r.Context(), folder, fvp, reqCont.Requester)
	} else {
		fvp := fschema.NewUserFileBuilder().WithTrustedSource([]byte(url_hash))
		file := fschema.NewUserFileBuilder().WithID(entry_id).WithOwnerID(reqCont.Requester.ID)

		repos.FileRepo.UpdateFileEntry(r.Context(), file, fvp, reqCont.Requester)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("HX-Trigger", formatters.AlertEventBuilder("File Shared", "File has been shared successfully", constants.Success))
	w.Write([]byte(url + path + query))
}
