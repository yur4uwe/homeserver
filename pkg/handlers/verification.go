package handlers

import (
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/setup"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func sendVerificationEmail(form globals.RegisterForm, verification_token string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "tsg9515@gmail.com")
	m.SetHeader("To", form.Email)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/html",
		fmt.Sprintf(`
			<p>Thank you for registering, %s!</p>
			<p>Please verify your email address by clicking the link below:</p>
			<p><a href="http://localhost:8080/api/v1/verification?token=%s">Verify Email</a></p>
			<p>If you did not register, please ignore this email.</p>
			<p>Best regards,<br>Your HomeServer Team</p>
		`, form.Username, verification_token),
	)

	port, err := strconv.Atoi(os.Getenv(setup.EnvEmailPort))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid email port: %v\n", err.Error())
		return
	}

	d := gomail.NewDialer(os.Getenv(setup.EnvEmailHost), port, os.Getenv(setup.EnvEmailUser), os.Getenv(setup.EnvEmailPassword))
	if err := d.DialAndSend(m); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send verification email: %v\n", err.Error())
		return
	}
}

func VerificationHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in VerificationHandler\n")
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.AddLog("[Handler] Verification endpoint accessed")

		token := r.URL.Query().Get("token")
		if token == "" {
			fmt.Fprintf(os.Stderr, "Missing verification token in request\n")
			http.Error(w, "Bad Request: Missing verification token", http.StatusBadRequest)
			return
		}

		_, err := connection.CheckToken(r.Context(), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check verification token: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Handle verification logic here
		// For example, you might check a verification token or send a verification email
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Verification endpoint reached"))
	}
}
