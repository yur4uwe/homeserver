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
			w.Header().Set("HX-Retarget", "#id-fieldset-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")

			new_fieldset := `
			<div id="id-fieldset-wrapper">
			<fieldset id="identificator-fieldset" class="field-error">	
				<legend>Username or Email</legend>
				<input type="text" name="identificator" placeholder="Your Username or Email">
			</fieldset>
			<div class="field-error-message">Identificator is required</div>
			</div>`

			w.Write([]byte(new_fieldset))
			return
		}

		if form.Password == "" {
			w.Header().Set("HX-Retarget", "#pwd-fieldset-wrapper")
			w.Header().Set("HX-Reswap", "outerHTML")

			new_fieldset := `
			<div id="pwd-fieldset-wrapper">
			<fieldset id="password-fieldset" class="field-error">
				<legend>Password</legend>
				<input type="password" name="password" placeholder="Your Password">
			</fieldset>
			<div class="field-error-message">Password is required</div>
			</div>`

			w.Write([]byte(new_fieldset))
			return
		}

		userId, err := connection.QueryUserOnLogin(r.Context(), form)
		if err != nil {
			if err.Error() == "no rows in result set" {
				fmt.Fprintf(os.Stderr, "Login failed for user %s: %v\n", form.Identificator, err.Error())
				http.Error(w, "Unauthorized: Invalid credentials", http.StatusUnauthorized)
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

		http.Redirect(w, r, "/", http.StatusSeeOther)
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

		if form.Password != form.Confirm {
			fmt.Fprintf(os.Stderr, "Passwords do not match: %s != %s\n", form.Password, form.Confirm)
			http.Error(w, "Bad Request: Passwords do not match", http.StatusBadRequest)
			return
		}

		// _, err = connection.CreateUser(r.Context(), form.Username, form.Email, form.Password)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to create user: %v\n", err.Error())
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
