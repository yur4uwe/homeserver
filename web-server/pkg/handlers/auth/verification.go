package auth

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db/repositories"
	"homeserver/web-server/pkg/middleware"
	"homeserver/web-server/pkg/parser"
	"homeserver/web-server/pkg/templates"
	"homeserver/web-server/pkg/utils/errresp"
	"homeserver/web-server/pkg/utils/formatters"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

var (
	ErrorResponse = errresp.ErrorResponse
	NewAlert      = parser.NewAlert
	LoadIcon      = parser.LoadIcon
	VerFailed     = SendVerificationFailurePage
	Must          = template.Must
	ParseFiles    = template.ParseFiles
)

func SendVerificationFailurePage(w http.ResponseWriter, verification *templates.VerificationPageData, errCode int) {
	tmpl := Must(ParseFiles(constants.Rootd + "static/templates/verification.html"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(errCode)
	verification.Success = false
	tmpl.Execute(w, verification)
}

func SendVerificationEmail(form formatters.RegisterForm, verification_token string) {
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

	port, _ := strconv.Atoi(os.Getenv(constants.EmailPort))

	d := gomail.NewDialer(os.Getenv(constants.EmailHost), port, os.Getenv(constants.EmailUser), os.Getenv(constants.EmailPassword))
	if err := d.DialAndSend(m); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send verification email: %v\n", err.Error())
		return
	}
}

func VerificationHandler(repos *repositories.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in VerificationHandler\n")
			ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Verification endpoint accessed")

		verification := &templates.VerificationPageData{
			Success: false,
			Error:   "",
			Causes:  []template.HTML{},
			Icon:    template.HTML(LoadIcon("x")),
		}

		token := r.URL.Query().Get("token")
		if token == "" {
			reqCont.LogHistory.ErrLog("Missing verification token in request\n", "VerificationHandler:58")
			verification.Error = "Missing verification token"
			verification.Causes = append(verification.Causes, template.HTML("You shouldn't be here"))
			VerFailed(w, verification, http.StatusBadRequest)
			return
		}

		reqCont.LogHistory.Log("[Handler] Verification token received: '" + token + "'")

		err := repos.TokenRepo.VerifyToken(r.Context(), token)
		if err != nil {
			reqCont.LogHistory.ErrLog("Failed to check verification token: "+err.Error(), "VerificationHandler:64")
			verification.Error = "Failed to check verification token"
			verification.Causes = append(verification.Causes, template.HTML("Invalid or expired token"))
			verification.Causes = append(verification.Causes, template.HTML("You refreshed the page with already verified token"))
			VerFailed(w, verification, http.StatusInternalServerError)
			return
		}

		verification.Success = true
		verification.Icon = template.HTML(LoadIcon("check"))

		alert := NewAlert("Verification success", "check", "Your email has been verified successfully!", constants.Success)

		reqCont.PageContent.ModalContent = template.HTML(alert)
		reqCont.RerenderBody = true

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		tmpl := Must(ParseFiles(constants.Rootd + "static/templates/verification.html"))
		tmpl.Execute(w, verification)
	}
}
