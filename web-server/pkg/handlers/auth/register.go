package auth

import (
	"fmt"
	"homeserver/pkg/auth"
	"homeserver/pkg/globals"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/templates"
	"net/http"
	"os"
)

var (
	verSuccess = sendVerificationEmail
	NewToken   = auth.NewToken
)

func RegisterHandler(connection interfaces.DatabaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in RegisterHandler\n")
			ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Register endpoint accessed")

		form := globals.RegisterFormFromRequest(r)

		valid, target, fieldsetHTML := RegisterValidation(form)
		if !valid {
			w.Header().Set("HX-Retarget", "#"+target+"-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")
			w.Write([]byte(fieldsetHTML))
			return
		}

		user := user_schema.NewUserBuilder().
			WithUsername(form.Username).
			WithEmail(form.Email).
			WithPassword(form.Password)

		inserted_user, err := connection.CreateUserEntry(r.Context(), user)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create user: %v\n", err.Error())
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		verification_token, err := NewToken()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate verification token: %v\n", err.Error())
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = connection.CreateToken(r.Context(), inserted_user.ID, verification_token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create verification token: %v\n", err.Error())
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		go verSuccess(form, verification_token)

		w.Header().Set("HX-Trigger", `{"raise-alert": {"title": "Success!", "value": "Registration successful! Please check your email for verification.", "status": "success", "icon": "check"}}`)
		w.Header().Set("HX-Push-Url", "/")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		tmpl := Must(ParseFiles("../static/templates/welcome.html"))
		if err := tmpl.Execute(w, &templates.WelcomePageData{IsLoggedIn: false}); err != nil {
			reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "AuthenticationHandler:196")
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
