package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"homeserver/pkg/auth"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"

	"golang.org/x/net/html"
)

func LoginHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := middleware.AddLog(r, "[Handler] Login endpoint accessed")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		form, err := globals.LoginFormFromRequest(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse login form: %v\n", err.Error())
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		if form.Identificator == "" {
			w.Header().Set("HX-Retarget", "#identificator-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")

			new_fieldset := fieldset("identificator", "Identificator is required", "field-error", "", "text", "Username or Email", "login")

			w.Write([]byte(new_fieldset))
			return
		}

		if form.Password == "" {
			w.Header().Set("HX-Retarget", "#password-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")

			new_fieldset := fieldset("password", "Password is required", "field-error", "", "password", "Password", "login")

			w.Write([]byte(new_fieldset))
			return
		}

		userId, err := connection.QueryUserOnLogin(r.Context(), form)
		if err != nil {
			if err.Error() == "no rows in result set" {
				fmt.Fprintf(os.Stderr, "Login failed for user %s: %v\n", form.Identificator, err.Error())

				w.Header().Set("HX-Retarget", "#error-resp")
				w.Header().Set("HX-Reswap", "innerHTML")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized: Invalid credentials"))
				return
			}
			fmt.Fprintf(os.Stderr, "Failed to query user on login: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		newToken := auth.NewToken(userId)

		err = connection.CreateToken(r.Context(), userId, newToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create token: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "authTokenHS",
			Value:    newToken,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		w.Header().Set("HX-Push-Url", "/")
		w.WriteHeader(http.StatusSeeOther)
		w.Write([]byte(""))
	}
}

func RegisterHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in RegisterHandler\n")
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory = append(reqCont.LogHistory, "[Handler] Register endpoint accessed")

		form, err := globals.RegisterFormFromRequest(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse registration form: %v\n", err.Error())
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		if form.Username == "" {
			w.Header().Set("HX-Retarget", "#username-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")

			new_fieldset := fieldset("username", "Username is required", "field-error", form.Username, "text", "Username", "register")
			w.Write([]byte(new_fieldset))
			return
		}

		if form.Email == "" {
			w.Header().Set("HX-Retarget", "#email-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")
			new_fieldset := fieldset("email", "Email is required", "field-error", form.Email, "email", "Your Email", "register")
			w.Write([]byte(new_fieldset))
			return
		}

		if form.Password == "" {
			w.Header().Set("HX-Retarget", "#password-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")
			new_fieldset := fieldset("password", "Password is required", "field-error", form.Password, "password", "Password", "register")
			w.Write([]byte(new_fieldset))
			return
		}

		if form.Confirm == "" {
			w.Header().Set("HX-Retarget", "#confirm-password-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")
			new_fieldset := fieldset("confirm-password", "Confirm Password is required", "field-error", form.Confirm, "password", "Confirm Password", "register")
			w.Write([]byte(new_fieldset))
			return
		}

		if form.Password != form.Confirm {
			w.Header().Set("HX-Retarget", "#confirm-password-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")
			new_fieldset := fieldset("confirm-password", "Passwords do not match", "field-error", form.Confirm, "password", "Confirm Password", "register")
			w.Write([]byte(new_fieldset))
			return
		}

		// _, err = connection.CreateUser(r.Context(), form.Username, form.Email, form.Password)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to create user: %v\n", err.Error()), "login"
		// 	http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		alert, err := parser.ParseString(`
			<div class="alert alert-success" hx-get="/nothing" hx-trigger="load delay:3s" hx-swap="delete">
				<strong>Success!</strong> You have registered successfully. Now log in
			</div>
		`)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse alert HTML: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		index, err := parser.ParseFile("../static/index.html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse index HTML file: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		welcomeNode, err := parser.ParseFile("../static/fragments/welcome.html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse welcome HTML file: %v\n", err.Error())
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		parser.InsertNodeIntoPage(index, welcomeNode, "main-content")
		parser.InsertNodeIntoPage(index, alert, "modal-content")

		reqCont.FullPageRender = true

		w.Header().Set("HX-Push-Url", "/")
		w.Header().Set("HX-Retarget", "body")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		html.Render(w, index)
	}
}
