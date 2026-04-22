package auth

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db"
	"homeserver/internal/templates"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/utils"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

func LoginHandler(repos *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in LoginHandler\n")
			ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Login endpoint accessed")

		form := utils.LoginFormFromRequest(r)

		valid, target, fieldsetHTML := LoginValidation(form)
		if !valid {
			w.Header().Set("HX-Retarget", "#"+target+"-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")
			w.Write([]byte(fieldsetHTML))
			return
		}

		q_user := &db.User{
			Password: form.Password,
		}

		if strings.Contains(form.Identificator, "@") {
			q_user.Email = form.Identificator
		} else {
			q_user.Username = form.Identificator
		}

		user, err := repos.Users.Query(r.Context(), q_user)
		if err != nil {
			if strings.Contains(err.Error(), "no identifier") || strings.Contains(err.Error(), "no rows") {
				fmt.Fprintf(os.Stderr, "Login failed for user %s: %v\n", form.Identificator, err.Error())

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Header().Set("HX-Retarget", "#error-resp")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Unauthorized: Invalid credentials"))
				return
			}
			fmt.Fprintf(os.Stderr, "Failed to query user on login: %v\n", err.Error())
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if !user.EmailVerified {
			fmt.Fprintf(os.Stderr, "User %d has not verified their email\n", user.ID)
			w.Header().Set("HX-Retarget", "#error-resp")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden: Email not verified"))
			return
		}

		expires := time.Now().Add(24 * time.Hour)

		newToken, err := NewToken(user.Email, expires)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate new token: %v\n", err.Error())
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err = repos.Tokens.CreateToken(r.Context(), user.ID, newToken); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create token: %v\n", err.Error())
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "authTokenHS",
			Value:    newToken,
			Path:     "/",
			Expires:  expires,
			HttpOnly: true,
		})

		reqCont.RerenderBody = true
		reqCont.AuthToken = newToken

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("HX-Push-Url", "/")

		welcomeData := &templates.WelcomePageData{
			IsLoggedIn: true,
		}

		header, err := ParseTemplateToString(constants.Rootd+"static/templates/header.html", welcomeData)
		if err != nil {
			reqCont.LogHistory.ErrLog("Failed to parse header template: "+err.Error(), "LoginHandler:66")
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		reqCont.PageContent.Header = template.HTML(header)

		tmpl := Must(ParseFiles(constants.Rootd + "static/templates/welcome.html"))
		if err := tmpl.Execute(w, welcomeData); err != nil {
			reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "LoginHandler:66")
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
