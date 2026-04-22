package constants

import (
	"encoding/json"
	"os"
)

// Environment Variable Keys
const (
	WebSrvPort    = "MAIN_PORT"
	WebSrvHost    = "MAIN_HOST"
	FileSrvPort   = "SECONDARY_PORT"
	FileSrvHost   = "SECONDARY_HOST"
	StaticPort    = "STATIC_PORT"
	StaticHost    = "STATIC_HOST"
	DbURL         = "DB_URL"
	StoragePath   = "STORAGE_PATH"
	EmailPassword = "EMAIL_PASSWORD"
	EmailHost     = "EMAIL_HOST"
	EmailPort     = "EMAIL_PORT"
	EmailUser     = "EMAIL_USER"
	SourceSecret  = "SOURCE_SECRET"
	TokenSecret   = "TOKEN_SECRET"
)

// Alert Types
const (
	Success = "success"
	Warning = "warning"
	Error   = "error"
	Info    = "info"
)

// Context & Cookie Keys
type CtxKey string

const (
	RequesterKey CtxKey = "requester"
	EntryKey     CtxKey = "entry"
	StatusKey    CtxKey = "status"
)

const (
	AuthCookieName = "authTokenHS"
	SIDQueryKey    = "sid"
	ExpQueryKey    = "exp"
)

// Global Variables
var Rootd = ""
var AvailableFileTypes = []string{}

// Pages
var (
	Pages = []string{
		"/login",
		"/register",
		"/verification",
		"/settings",
		"/upload",
		"/search/",
	}

	ProtectedPages = []string{
		"/settings",
		"/api/v1/settings",
		"/upload",
		"/search/",
	}
)

func ReadAvailableFileTypes() {
	data, err := os.ReadFile(Rootd + "resources/catalog.json")
	if err != nil {
		panic("unable to read catalog.json: " + err.Error())
	}

	err = json.Unmarshal(data, &AvailableFileTypes)
	if err != nil {
		panic("unable to parse catalog.json: " + err.Error())
	}
}
