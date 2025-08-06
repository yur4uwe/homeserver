package actions

import (
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/utils"
	"net/http"
	"strings"
)

func OpenFileHandler(w http.ResponseWriter, r *http.Request, connection *db.Database, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName string) {
	file := file_schema.NewUserFileBuilder().
		WithFilename(strings.Split(fileName, ".")[0]).
		WithPath(strings.Join(relativePath, "/")).
		WithOwnerID(reqCont.Requester.ID)

	file, err := connection.GetFile(r.Context(), file, reqCont.Requester)
	if err != nil {
		utils.ErrorResponse(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, file.Path)
}
