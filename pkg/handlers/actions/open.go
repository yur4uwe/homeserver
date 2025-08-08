package actions

import (
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas/file_schema"
	"net/http"
	"strings"
)

var (
	ServeFile = http.ServeFile
)

func OpenFileHandler(w http.ResponseWriter, r *http.Request, connection interfaces.DatabaseInterface, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName string) {
	file := file_schema.NewUserFileBuilder().
		WithFilename(strings.Split(fileName, ".")[0]).
		WithExtension(strings.Split(fileName, ".")[1]).
		WithPath(strings.Join(relativePath, "/")).
		WithOwnerID(reqCont.Requester.ID)

	file, err := connection.GetFile(r.Context(), file, reqCont.Requester)
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
		return
	}

	ServeFile(w, r, file.Path)
}
