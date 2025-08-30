package auth

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db/repositories"
	"homeserver/internal/templates"
	"homeserver/pkg/crypto"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/parser"
	"homeserver/pkg/utils/errresp"
	"homeserver/pkg/utils/formatters"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

var (
	ErrorResponse = errresp.ErrorResponse
	NewAlert      = parser.NewAlert
	LoadIcon      = parser.LoadIcon
	VerFailed     = SendVerificationFailurePage
	Must          = template.Must
	ParseFiles    = template.ParseFiles
	ValidateToken = crypto.ValidateToken
)

func SendVerificationFailurePage(w http.ResponseWriter, verification *templates.VerificationPageData, errCode int) {
	tmpl := Must(ParseFiles(constants.Rootd + "static/templates/verification.html"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(errCode)
	verification.Icon = template.HTML(LoadIcon("x"))
	verification.Success = false
	tmpl.Execute(w, verification)
}

func SendVerificationEmail(form formatters.RegisterForm, verification_token string) {
	expires := time.Now().Add(24 * time.Hour)

	m := gomail.NewMessage()
	m.SetHeader("From", "tsg9515@gmail.com")
	m.SetHeader("To", form.Email)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/html",
		fmt.Sprintf(`
			<p>Thank you for registering, %s!</p>
			<p>Please verify your email address by clicking the link below:</p>
			<p><a href="http://localhost:8080/verification?vt=%s&exp=%d">Verify Email</a></p>
			<p>If you did not register, please ignore this email.</p>
			<p>Best regards,<br>Your HomeServer Team</p>
		`, form.Username, verification_token, expires.Unix()),
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
		}

		expirationStr := r.URL.Query().Get("exp")
		expiration, err := strconv.ParseInt(expirationStr, 10, 64)
		if err != nil {
			verification.Error = "Tampered verification link"
			VerFailed(w, verification, http.StatusBadRequest)
			return
		}

		if time.Now().Unix() > expiration {
			verification.Error = "Verification link expired"
			verification.Causes = append(verification.Causes, template.HTML("The link has expired. Please request a new verification email."))
			VerFailed(w, verification, http.StatusBadRequest)
			return
		}

		token := r.URL.Query().Get("vt")
		if token == "" {
			verification.Error = "Missing verification token"
			verification.Causes = append(verification.Causes, template.HTML("You shouldn't be here"))
			VerFailed(w, verification, http.StatusBadRequest)
			return
		}

		claims, valid := ValidateToken(token)
		if !valid {
			verification.Error = "Failed to verify token"
			verification.Causes = append(verification.Causes, template.HTML("Invalid or expired token"))
			VerFailed(w, verification, http.StatusBadRequest)
			return
		}

		reqCont.LogHistory.Log("[Handler] Verification token received: '" + token + "'")

		user, err := repos.TokenRepo.QueryToken(r.Context(), token, false)
		if err != nil {
			verification.Error = "Failed to check verification token"
			verification.Causes = append(verification.Causes, template.HTML("Invalid or expired token"))
			verification.Causes = append(verification.Causes, template.HTML("You refreshed the page with already verified token"))
			VerFailed(w, verification, http.StatusInternalServerError)
			return
		}

		if claims.Seed != user.Email {
			verification.Error = "Failed to verify email"
			verification.Causes = append(verification.Causes, template.HTML("The email associated with this token is invalid."))
			VerFailed(w, verification, http.StatusBadRequest)
		}

		if err := repos.TokenRepo.VerifyToken(r.Context(), token); err != nil {
			verification.Error = "Failed to verify email"
			verification.Causes = append(verification.Causes, template.HTML("The email associated with this token is invalid."))
			VerFailed(w, verification, http.StatusBadRequest)
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
