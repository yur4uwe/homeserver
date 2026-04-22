package middleware

import (
	"bytes"
	"homeserver/db"
	"homeserver/internal/templates"
	"homeserver/pkg/parser"
	"homeserver/pkg/utils"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

// Global variables for mocking
var (
	ParseFiles            = template.ParseFiles
	Must                  = template.Must
	IconsInjector         = utils.IconsInjector
	ParseTemplateToString = parser.ParseTemplateToString
	ErrorResponse         = utils.ErrorResponse
	Now                   = time.Now
	Since                 = time.Since
	IsProtectedPath       = utils.IsProtectedPath
)

// Context key for storing AdditionalRequestContext in request context
type contextKey string

const ReqContKey contextKey = "add_req_context"

// LogHistoryInterface defines logging methods
type LogHistoryInterface interface {
	Log(entry string)
	ErrLog(entry, at string)
	UnloadLogs(header string)
}

// LogHistory implements LogHistoryInterface
type LogHistory struct {
	Entries []string
}

func (lh *LogHistory) Log(entry string) {
	lh.Entries = append(lh.Entries, entry)
}

func (lh *LogHistory) ErrLog(entry, at string) {
	lh.Entries = append(lh.Entries, "[Err at "+at+"] "+entry)
}

func (lh *LogHistory) UnloadLogs(header string) {
	logs := strings.Join(lh.Entries, "\n") + "\n\n"
	log.Println("\n\n" + header + logs)
}

// AdditionalRequestContext holds request-specific data
type AdditionalRequestContext struct {
	LogHistory    LogHistoryInterface
	ToRestorePage bool
	RerenderBody  bool
	AuthToken     string
	Requester     *db.User
	PageContent   *templates.LayoutPageData
}

// Default context initialization
var DefaultContext = &AdditionalRequestContext{
	LogHistory:    &LogHistory{Entries: []string{}},
	ToRestorePage: false,
	RerenderBody:  false,
	AuthToken:     "",
	Requester:     nil,
	PageContent:   &templates.LayoutPageData{},
}

// Response interceptor for capturing response details
type responseInterceptor struct {
	statusCode int
	body       *bytes.Buffer
	http.ResponseWriter
}

func InterceptorConstructor(w http.ResponseWriter) *responseInterceptor {
	return &responseInterceptor{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           new(bytes.Buffer),
	}
}

func (ri *responseInterceptor) WriteHeader(code int) {
	ri.statusCode = code
}

func (ri *responseInterceptor) Write(b []byte) (int, error) {
	return ri.body.Write(b)
}
