package handlers

import (
	"fmt"
	"homeserver/pkg/db"
	"net/http"
	"strings"
)

func fieldset(fieldId, message, class, value, inputType, placeholder, field_designation string) string {
	return `
	<div id="` + fieldId + `-wrapper">
		<fieldset class="` + class + `">
			<legend>` + placeholder + `</legend>
			<input 
				type="` + inputType + `" 
				hx-post="/api/v1/` + field_designation + `/check?field=` + fieldId + `" 
				hx-trigger="keyup changed delay:500ms" 
				name="` + fieldId + `" 
				placeholder="Your ` + placeholder + `" 
				value="` + value + `"
			>
		</fieldset>
		<div class="field-error-message">` + message + `</div>
	</div>`
}

func LoginValidationHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fieldId := r.URL.Query().Get("field")
		fmt.Println("LoginValidationHandler called for field:", fieldId)
		value := r.FormValue(fieldId)
		var fieldsetHTML string

		switch fieldId {
		case "identificator":
			fieldsetHTML = fieldset("identificator", "", "", value, "text", "Username or Email", "login")
			if value == "" {
				fieldsetHTML = fieldset("identificator", "Identificator is required", "field-error", value, "text", "Username or Email", "login")
			}
		case "password":
			fieldsetHTML = fieldset("password", "", "", value, "password", "Password", "login")
			if value == "" {
				fieldsetHTML = fieldset("password", "Password is required", "field-error", value, "password", "Password", "login")
			}
		default:
			http.Error(w, "Bad Request: Invalid field", http.StatusBadRequest)
			return
		}

		w.Header().Set("HX-Retarget", "#"+fieldId+"-wrapper")
		w.Header().Set("HX-Reswap", "outerHTML")
		w.Write([]byte(fieldsetHTML))
	}
}

func RegisterValidationHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fieldId := strings.TrimSpace(r.URL.Query().Get("field"))
		fmt.Println("RegisterValidationHandler called for field:", fieldId)
		value := r.FormValue(fieldId)
		var fieldsetHTML string

		switch fieldId {
		case "username":
			fieldsetHTML = fieldset("username", "", "", value, "text", "Username", "register")
			if value == "" {
				fieldsetHTML = fieldset("username", "Username is required", "field-error", value, "text", "Username", "register")
			}
		case "email":
			fieldsetHTML = fieldset("email", "", "", value, "email", "Your Email", "register")
			if value == "" {
				fieldsetHTML = fieldset("email", "Email is required", "field-error", value, "email", "Your Email", "register")
			}
		case "password":
			fieldsetHTML = fieldset("password", "", "", value, "password", "Password", "register")
			if value == "" {
				fieldsetHTML = fieldset("password", "Password is required", "field-error", value, "password", "Password", "register")
			}
		case "confirm-password":
			fieldsetHTML = fieldset("confirm-password", "", "", value, "password", "Confirm Password", "register")
			if value == "" {
				fieldsetHTML = fieldset("confirm-password", "Confirm Password is required", "field-error", value, "password", "Confirm Password", "register")
			} else if r.FormValue("password") != value {
				fieldsetHTML = fieldset("confirm-password", "Passwords do not match", "field-error", value, "password", "Confirm Password", "register")
			}
		default:
			http.Error(w, "Bad Request: Invalid field", http.StatusBadRequest)
			return
		}

		w.Header().Set("HX-Retarget", "#"+fieldId+"-wrapper")
		w.Header().Set("HX-Reswap", "outerHTML")
		w.Write([]byte(fieldsetHTML))
	}
}
