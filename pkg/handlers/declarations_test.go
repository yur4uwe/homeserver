package handlers

import (
	"homeserver/pkg/interfaces"
	"homeserver/pkg/middleware"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type UploadMock struct {
	mock.Mock
}

func (m *UploadMock) DownloadFile(connection interfaces.DatabaseInterface, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	args := m.Called(connection, reqCont)
	return args.Get(0).(http.HandlerFunc)
}

func (m *UploadMock) DeleteFile(connection interfaces.DatabaseInterface, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	args := m.Called(connection, reqCont)
	return args.Get(0).(http.HandlerFunc)
}

func (m *UploadMock) PatchFile(connection interfaces.DatabaseInterface, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	args := m.Called(connection, reqCont)
	return args.Get(0).(http.HandlerFunc)
}

type ActionsMock struct {
	mock.Mock
}

func (m *ActionsMock) DeleteFileHandler(w http.ResponseWriter, r *http.Request, connection interfaces.DatabaseInterface, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName string, isFolder bool) {
	m.Called(w, r, connection, reqCont, relativePath, fileName, isFolder)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted successfully"))
}

func (m *ActionsMock) RenameFileHandler(w http.ResponseWriter, r *http.Request, connection interfaces.DatabaseInterface, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName, desiredName string, isFolder bool) {
	m.Called(w, r, connection, reqCont, relativePath, fileName, desiredName, isFolder)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File renamed successfully"))
}

func (m *ActionsMock) OpenFileHandler(w http.ResponseWriter, r *http.Request, connection interfaces.DatabaseInterface, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName string) {
	m.Called(w, r, connection, reqCont, relativePath, fileName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File opened successfully"))
}
