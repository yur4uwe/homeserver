package auth

import (
	"fmt"
	"homeserver/pkg/auth"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/templates"
	"homeserver/pkg/utils"
	"html/template"
	"net/http"
	"os"
)

func RegisterHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in RegisterHandler\n")
			utils.ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse registration form: %v\n", err.Error())
			utils.ErrorResponse(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
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
			utils.ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		verification_token, err := auth.NewToken()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate verification token: %v\n", err.Error())
			utils.ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = connection.CreateToken(r.Context(), inserted_user.ID, verification_token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create verification token: %v\n", err.Error())
			utils.ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		go sendVerificationEmail(form, verification_token)

		alert := parser.NewAlert("Register successful", "check", "You have registered successfully. Now verify your email address.", utils.AlertSuccess)

		reqCont.PageContent.ModalContent = template.HTML(alert)

		reqCont.RerenderBody = true

		w.Header().Set("HX-Push-Url", "/")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		tmpl := template.Must(template.ParseFiles("../static/templates/welcome.html"))
		if err := tmpl.Execute(w, &templates.WelcomePageData{IsLoggedIn: false}); err != nil {
			reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "AuthenticationHandler:196")
			utils.ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
