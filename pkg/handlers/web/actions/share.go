package actions

import (
	"fmt"
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

func ShareFileHandler(w http.ResponseWriter, r *http.Request, repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext, fileId int) {
	isFolder := r.URL.Query().Get("is_folder") == "true"

	file_srv_host := os.Getenv(constants.FileSrvHost)
	file_srv_port := os.Getenv(constants.FileSrvPort)

	url := "http://" + file_srv_host + ":" + file_srv_port

	path := ""
	if isFolder {
		path += "/folders/" + strconv.Itoa(fileId)
	} else {
		path += "/files/" + strconv.Itoa(fileId)
	}

	exp := strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10)
	data := path + "|" + exp

	url_hash := crypto.HashSeed(data)

	query := "?sid=" + url_hash + "&exp=" + exp

	if isFolder {
		fvp := folschema.New().WithTrustedSource([]byte(url_hash))
		folder := folschema.New().WithID(fileId).WithOwnerID(reqCont.Requester.ID)

		newfol, err := repos.FolderRepo.UpdateFolder(r.Context(), folder, fvp, reqCont.Requester)
		if err != nil {
			ErrorResponse(w, "Error updating folder", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Folder %s updated successfully\n", newfol.TrustedSource)
	} else {
		fvp := fschema.New().WithTrustedSource([]byte(url_hash))
		file := fschema.New().WithID(fileId).WithOwnerID(reqCont.Requester.ID)

		newfile, err := repos.FileRepo.UpdateFileEntry(r.Context(), file, fvp, reqCont.Requester)
		if err != nil {
			ErrorResponse(w, "Error updating file", http.StatusInternalServerError)
			return
		}

		fmt.Printf("File %s updated successfully\n", newfile.TrustedSource)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("HX-Trigger", formatters.AlertEventBuilder("File Shared", "File has been shared successfully", constants.Success))
	w.Write([]byte(url + path + query))
}
