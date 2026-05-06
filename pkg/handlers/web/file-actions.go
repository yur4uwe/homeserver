package handlers

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db"
	"homeserver/pkg/bucket"
	"homeserver/pkg/crypto"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/parser"
	"homeserver/pkg/utils"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

type FileActionHandler struct {
	Repos  *db.Database
	Bucket bucket.Bucket
	logger *slog.Logger
}

// Helpers

func (f *FileActionHandler) getID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		middleware.ErrorResponse(w, "Invalid ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

// Upload Logic

func (f *FileActionHandler) Upload(w http.ResponseWriter, r *http.Request) {}

// Delete Logic

func (f *FileActionHandler) deleteFolder(w http.ResponseWriter, r *http.Request, entryId int) string {
	reqCont, ok := middleware.GetContext(w, r)
	if !ok {
		return ""
	}

	folder := &db.UserFolder{ID: entryId, OwnerID: reqCont.Requester.ID}

	fol, err := f.Repos.Folders.Delete(r.Context(), folder, reqCont.Requester)
	if err != nil {
		middleware.ErrorResponse(w, fmt.Sprintf("Error deleting folder: %v", err), http.StatusInternalServerError)
		return ""
	}

	if err := f.Bucket.Delete(fol.RelativePath()); err != nil {
		middleware.ErrorResponse(w, fmt.Sprintf("Error removing folder from storage: %v", err), http.StatusInternalServerError)
		return ""
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("delete-success-folder-%d", fol.ID))

	return parser.NewAlert("Folder Deleted", "check", "Folder has been deleted successfully.", constants.Success)
}

func (f *FileActionHandler) deleteFile(w http.ResponseWriter, r *http.Request, entryId int) string {
	reqCont, ok := middleware.GetContext(w, r)
	if !ok {
		return ""
	}

	file := &db.UserFile{ID: entryId, OwnerID: reqCont.Requester.ID}

	uf, err := f.Repos.Files.Delete(r.Context(), file, reqCont.Requester)
	if err != nil {
		middleware.ErrorResponse(w, fmt.Sprintf("Error deleting file: %v", err), http.StatusInternalServerError)
		return ""
	}

	if err := f.Bucket.Delete(uf.RelativePath()); err != nil {
		middleware.ErrorResponse(w, fmt.Sprintf("Error removing file from storage: %v", err), http.StatusInternalServerError)
		return ""
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("delete-success-file-%d", uf.ID))

	return parser.NewAlert("File Deleted", "check", "File has been deleted successfully.", constants.Success)
}

func (f *FileActionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	entryId, ok := f.getID(w, r)
	if !ok {
		return
	}

	isFolder := r.URL.Query().Get("is_folder") == "true"

	var alert string
	if isFolder {
		alert = f.deleteFolder(w, r, entryId)
	} else {
		alert = f.deleteFile(w, r, entryId)
	}

	if alert == "" {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(alert))
}

// Rename Logic

func (f *FileActionHandler) renameFile(w http.ResponseWriter, r *http.Request, fileId int, desiredName string) string {
	reqCont, ok := middleware.GetContext(w, r)
	if !ok {
		return ""
	}

	oldFile, err := f.Repos.Files.Get(r.Context(), &db.UserFile{ID: fileId, OwnerID: reqCont.Requester.ID}, reqCont.Requester)
	if err != nil {
		middleware.ErrorResponse(w, "File not found", http.StatusNotFound)
		return ""
	}

	updates := &db.UserFile{Filename: desiredName}
	nf, err := f.Repos.Files.Update(r.Context(), oldFile, updates, reqCont.Requester)
	if err != nil {
		middleware.ErrorResponse(w, fmt.Sprintf("Error updating file in database: %v", err), http.StatusInternalServerError)
		return ""
	}

	if err := f.Bucket.Move(oldFile.RelativePath(), nf.RelativePath()); err != nil {
		middleware.ErrorResponse(w, fmt.Sprintf("Error renaming file in storage: %v", err), http.StatusInternalServerError)
		return ""
	}

	w.Header().Set("HX-Trigger", utils.AlertEvent("Success!", "File successfully renamed", constants.Success))
	return nf.Filename + "." + nf.Extension
}

func (f *FileActionHandler) renameFolder(w http.ResponseWriter, r *http.Request, folderId int, desiredName string) string {
	reqCont, ok := middleware.GetContext(w, r)
	if !ok {
		return ""
	}

	oldFolder, err := f.Repos.Folders.Get(r.Context(), &db.UserFolder{ID: folderId, OwnerID: reqCont.Requester.ID}, reqCont.Requester)
	if err != nil {
		middleware.ErrorResponse(w, "Folder not found", http.StatusNotFound)
		return ""
	}

	updates := &db.UserFolder{Name: desiredName}
	nf, err := f.Repos.Folders.Update(r.Context(), oldFolder, updates, reqCont.Requester)
	if err != nil {
		middleware.ErrorResponse(w, fmt.Sprintf("Error updating folder in database: %v", err), http.StatusInternalServerError)
		return ""
	}

	if err := f.Bucket.Move(oldFolder.RelativePath(), nf.RelativePath()); err != nil {
		middleware.ErrorResponse(w, fmt.Sprintf("Error renaming folder in storage: %v", err), http.StatusInternalServerError)
		return ""
	}

	w.Header().Set("HX-Trigger", utils.AlertEvent("Success!", "Folder successfully renamed", constants.Success))
	return nf.Name
}

func (f *FileActionHandler) Rename(w http.ResponseWriter, r *http.Request) {
	id, ok := f.getID(w, r)
	if !ok {
		return
	}

	isFolder := r.URL.Query().Get("is_folder") == "true"
	desiredName := r.URL.Query().Get("desired_name")

	if desiredName == "" {
		middleware.ErrorResponse(w, "Desired name is required", http.StatusBadRequest)
		return
	}

	var response string
	if isFolder {
		response = f.renameFolder(w, r, id, desiredName)
	} else {
		response = f.renameFile(w, r, id, desiredName)
	}

	if response != "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}

// Share Logic

func (f *FileActionHandler) performShare(w http.ResponseWriter, r *http.Request, id int, isFolder bool, reqCont *middleware.AdditionalRequestContext) string {
	file_srv_host := os.Getenv(constants.FileSrvHost)
	file_srv_port := os.Getenv(constants.FileSrvPort)
	baseUrl := "http://" + file_srv_host + ":" + file_srv_port

	path := "/files/" + strconv.Itoa(id)
	if isFolder {
		path = "/folders/" + strconv.Itoa(id)
	}

	exp := strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10)
	data := path + "|" + exp
	url_hash := crypto.HashSeed(data)
	query := "?sid=" + url_hash + "&exp=" + exp

	var err error
	if isFolder {
		updates := &db.UserFolder{TrustedSource: []byte(url_hash)}
		_, err = f.Repos.Folders.Update(r.Context(), &db.UserFolder{ID: id, OwnerID: reqCont.Requester.ID}, updates, reqCont.Requester)
	} else {
		updates := &db.UserFile{TrustedSource: []byte(url_hash)}
		_, err = f.Repos.Files.Update(r.Context(), &db.UserFile{ID: id, OwnerID: reqCont.Requester.ID}, updates, reqCont.Requester)
	}

	if err != nil {
		middleware.ErrorResponse(w, "Error updating shared source", http.StatusInternalServerError)
		return ""
	}

	w.Header().Set("HX-Trigger", utils.AlertEvent("File Shared", "Access link generated successfully", constants.Success))
	return baseUrl + path + query
}

func (f *FileActionHandler) Share(w http.ResponseWriter, r *http.Request) {
	id, ok := f.getID(w, r)
	if !ok {
		return
	}

	reqCont, ok := middleware.GetContext(w, r)
	if !ok {
		return
	}

	isFolder := r.URL.Query().Get("is_folder") == "true"

	link := f.performShare(w, r, id, isFolder, reqCont)

	if link != "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(link))
	}
}
