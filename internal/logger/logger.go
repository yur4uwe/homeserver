package logger

import (
	"fmt"
	"homeserver/constants"
	"net/http"
)

type LoggerInt interface {
	Error(msg string)
	Info(msg string)
	Warn(msg string)
	Request(r *http.Request)
}

type Logger struct{}

func (l *Logger) Error(msg string) {
	// Implement error logging
}

func (l *Logger) Info(msg string) {
	// Implement info logging
}

func (l *Logger) Warn(msg string) {
	// Implement warn logging
}

func (l *Logger) Request(r *http.Request) {
	status, ok := r.Context().Value(constants.StatusKey).(int)
	if !ok {
		status = http.StatusOK
	}
	fmt.Printf("%s %s %d\n", r.Method, r.URL.Path, status)
}
