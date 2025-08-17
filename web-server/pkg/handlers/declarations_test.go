package handlers

import (
	"context"
	"homeserver/db/repositories"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/web-server/pkg/middleware"
	"homeserver/web-server/pkg/templates"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type UploadMock struct {
	mock.Mock
}

func (m *UploadMock) DownloadFile(repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	args := m.Called(repos, reqCont)
	return args.Get(0).(http.HandlerFunc)
}

func (m *UploadMock) DeleteFile(repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	args := m.Called(repos, reqCont)
	return args.Get(0).(http.HandlerFunc)
}

func (m *UploadMock) PatchFile(repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	args := m.Called(repos, reqCont)
	return args.Get(0).(http.HandlerFunc)
}

type ActionsMock struct {
	mock.Mock
}

func (m *ActionsMock) DeleteFileHandler(w http.ResponseWriter, r *http.Request, repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName string) {
	m.Called(w, r, repos, reqCont, relativePath, fileName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted successfully"))
}

func (m *ActionsMock) RenameFileHandler(w http.ResponseWriter, r *http.Request, repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName, desiredName string, isFolder bool) {
	m.Called(w, r, repos, reqCont, relativePath, fileName, desiredName, isFolder)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File renamed successfully"))
}

func (m *ActionsMock) OpenFileHandler(w http.ResponseWriter, r *http.Request, repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName string) {
	m.Called(w, r, repos, reqCont, relativePath, fileName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File opened successfully"))
}

type SearchMock struct {
	mock.Mock
}

func (m *SearchMock) FileSearchHandler(ctx context.Context, repos *repositories.Repositories, query string, directory []string, requester *uschema.User) (*templates.SearchPageData, error) {
	args := m.Called(ctx, repos, query, directory, requester)
	return args.Get(0).(*templates.SearchPageData), args.Error(1)
}

func (m *SearchMock) BrowseHandler(ctx context.Context, repos *repositories.Repositories, directory []string, requester *uschema.User) (*templates.SearchPageData, error) {
	args := m.Called(ctx, repos, directory, requester)
	return args.Get(0).(*templates.SearchPageData), args.Error(1)
}

type IsPageMock struct {
	mock.Mock
}

func (m *IsPageMock) IsPage(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

type WelcomeMock struct {
	mock.Mock
}

func (m *WelcomeMock) WelcomePageHandler(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("WelcomePage"))
}
