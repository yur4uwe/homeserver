package auth

import (
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/templates"
	"homeserver/pkg/utils"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func sendVerificationFailurePage(w http.ResponseWriter, verification *templates.VerificationPageData, errCode int) {
	tmpl := template.Must(template.ParseFiles("../static/templates/verification.html"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(errCode)
	verification.Success = false
	tmpl.Execute(w, verification)
}

func sendVerificationEmail(form globals.RegisterForm, verification_token string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "tsg9515@gmail.com")
	m.SetHeader("To", form.Email)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/html",
		fmt.Sprintf(`
			<p>Thank you for registering, %s!</p>
			<p>Please verify your email address by clicking the link below:</p>
			<p><a href="http://localhost:8080/verification?token=%s">Verify Email</a></p>
			<p>If you did not register, please ignore this email.</p>
			<p>Best regards,<br>Your HomeServer Team</p>
		`, form.Username, verification_token),
	)

	port, err := strconv.Atoi(os.Getenv(globals.EnvEmailPort))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid email port: %v\n", err.Error())
		return
	}

	d := gomail.NewDialer(os.Getenv(globals.EnvEmailHost), port, os.Getenv(globals.EnvEmailUser), os.Getenv(globals.EnvEmailPassword))
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
			utils.ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Verification endpoint accessed")

		verification := &templates.VerificationPageData{
			Success: false,
			Error:   "",
			Causes:  []template.HTML{},
			Icon:    template.HTML(parser.LoadIcon("x")),
		}

		token := r.URL.Query().Get("token")
		if token == "" {
			reqCont.LogHistory.ErrLog("Missing verification token in request\n", "VerificationHandler:58")
			verification.Error = "Missing verification token"
			verification.Causes = append(verification.Causes, template.HTML("You shouldn't be here"))
			sendVerificationFailurePage(w, verification, http.StatusBadRequest)
			return
		}

		reqCont.LogHistory.Log("[Handler] Verification token received: '" + token + "'")

		err := connection.VerifyToken(r.Context(), token)
		if err != nil {
			reqCont.LogHistory.ErrLog("Failed to check verification token: "+err.Error(), "VerificationHandler:64")
			verification.Error = "Failed to check verification token"
			verification.Causes = append(verification.Causes, template.HTML("Invalid or expired token"))
			verification.Causes = append(verification.Causes, template.HTML("You refreshed the page with already verified token"))
			sendVerificationFailurePage(w, verification, http.StatusInternalServerError)
			return
		}

		verification.Success = true
		verification.Icon = template.HTML(parser.LoadIcon("check"))

		alert := parser.NewAlert("Verification success", "check", "Your email has been verified successfully!", utils.AlertSuccess)

		reqCont.PageContent.ModalContent = template.HTML(alert)
		reqCont.RerenderBody = true

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		tmpl := template.Must(template.ParseFiles("../static/templates/verification.html"))
		tmpl.Execute(w, verification)
	}
}
