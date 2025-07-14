package handlers

import (
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/setup"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/net/html"
	"gopkg.in/gomail.v2"
)

func sendVerificationFaliurePage(w http.ResponseWriter, reqCont *middleware.AdditionalRequestContext, message string, errCode int) {
	verification_failure_page, err := parser.ParseFile("../static/fragments/verification_failed.html")
	if err != nil {
		reqCont.LogHistory.ErrLog("Failed to parse verification failure HTML file: "+err.Error(), "VerificationHandler:58")
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	header_node := parser.GetNodeById(verification_failure_page, "ver-header")
	if header_node == nil {
		reqCont.LogHistory.ErrLog("Failed to find header node in verification failure HTML", "VerificationHandler:62")
		http.Error(w, "Internal Server Error: Failed to find header node", http.StatusInternalServerError)
		return
	}

	check_icon, err := parser.ParseFile("../static/assets/check.svg")
	if err != nil {
		reqCont.LogHistory.ErrLog("Failed to parse check icon: "+err.Error(), "VerificationHandler:68")
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	header_node.AppendChild(check_icon)
	header_node.FirstChild.Data = "Verification Failed"

	error_message := fmt.Sprintf("Verification failed: %v", message)

	error_node := parser.GetNodeById(verification_failure_page, "error-msg")
	if error_node == nil {
		reqCont.LogHistory.ErrLog("Failed to find error message node in verification failure HTML", "VerificationHandler:66")
		http.Error(w, "Internal Server Error: Failed to find error message node", http.StatusInternalServerError)
		return
	}

	error_node.Data = error_message

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(errCode)
	html.Render(w, verification_failure_page)
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

		reqCont.LogHistory.Log("[Handler] Verification endpoint accessed")

		token := r.URL.Query().Get("token")
		if token == "" {
			reqCont.LogHistory.ErrLog("Missing verification token in request\n", "VerificationHandler:58")
			sendVerificationFaliurePage(w, reqCont, "Missing verification token token", http.StatusBadRequest)
			return
		}

		err := connection.VerifyToken(r.Context(), token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check verification token: %v\n", err.Error())
			sendVerificationFaliurePage(w, reqCont, "Token not found", http.StatusInternalServerError)
			return
		}

		alert := parser.NewAlert("Your email has been verified successfully!", parser.AlertSuccess)

		index, err := parser.ParseFile("../static/index.html")
		if err != nil {
			reqCont.LogHistory.ErrLog("Failed to parse index file: "+err.Error(), "VerificationHandler:78")
			sendVerificationFaliurePage(w, reqCont, "Internal server error", http.StatusInternalServerError)
			return
		}

		body := parser.GetNodeByTagName(index, "body")
		if body == nil {
			fmt.Fprintf(os.Stderr, "Failed to find body tag in index HTML\n")
			sendVerificationFaliurePage(w, reqCont, "Internal server error", http.StatusInternalServerError)
			return
		}

		verification_page, err := parser.ParseFile("../static/fragments/verification_success.html")
		if err != nil {
			reqCont.LogHistory.ErrLog("Failed to parse verification HTML file: "+err.Error(), "VerificationHandler:90")
			sendVerificationFaliurePage(w, reqCont, "Internal server error", http.StatusInternalServerError)
			return
		}

		header_node := parser.GetNodeById(verification_page, "ver-header")
		if header_node == nil {
			reqCont.LogHistory.ErrLog("Failed to find header node in verification HTML", "VerificationHandler:95")
			sendVerificationFaliurePage(w, reqCont, "Internal server error", http.StatusInternalServerError)
			return
		}

		check_icon, err := parser.ParseFile("../static/assets/check.svg")
		if err != nil {
			reqCont.LogHistory.ErrLog("Failed to parse check icon: "+err.Error(), "VerificationHandler:101")
			sendVerificationFaliurePage(w, reqCont, "Internal server error", http.StatusInternalServerError)
			return
		}

		header_node.AppendChild(check_icon)
		header_node.Data = "Verification Successful"

		parser.InsertNodeIntoPage(body, verification_page, "main-content")
		parser.InsertNodeIntoPage(body, alert, "modal-content")

		reqCont.FullPageRender = index

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}
