package auth

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db"
	"homeserver/internal/templates"
	"homeserver/pkg/crypto"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/utils"
	"net/http"
	"os"
	"time"
)

var (
	VerSuccess = SendVerificationEmail
	NewToken   = crypto.NewToken
)

func RegisterHandler(repos *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in RegisterHandler\n")
			ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Register endpoint accessed")

		form := utils.RegisterFormFromRequest(r)

		valid, target, fieldsetHTML := RegisterValidation(form)
		if !valid {
			w.Header().Set("HX-Retarget", "#"+target+"-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")
			w.Write([]byte(fieldsetHTML))
			return
		}

		user := &db.User{
			Username: form.Username,
			Email:    form.Email,
			Password: form.Password,
		}

		inserted_user, err := repos.Users.Create(r.Context(), user)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create user: %v\n", err.Error())
			ErrorResponse(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		verification_hash, err := NewToken(inserted_user.Email, time.Now().Add(24*time.Hour))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create verification token: %v\n", err.Error())
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = repos.Tokens.CreateToken(r.Context(), inserted_user.ID, verification_hash)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create verification token: %v\n", err.Error())
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		go VerSuccess(form, verification_hash)

		w.Header().Set("HX-Trigger", utils.AlertEventBuilder("Success!", "Registration successful! Please check your email for verification.", constants.Success))
		w.Header().Set("HX-Push-Url", "/")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		tmpl := Must(ParseFiles(constants.Rootd + "static/templates/welcome.html"))
		if err := tmpl.Execute(w, &templates.WelcomePageData{IsLoggedIn: false}); err != nil {
			reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "AuthenticationHandler:196")
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
