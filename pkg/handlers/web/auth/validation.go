package auth

import (
	"homeserver/pkg/utils"
)

func fieldset(fieldId, message, value, inputType, placeholder string, isError bool) string {
	class := ""
	if isError {
		class = "field-error"
	}

	return `
	<div id="` + fieldId + `-wrapper">
		<fieldset class="` + class + `">
			<legend>` + placeholder + `</legend>
			<input 
				type="` + inputType + `" 
				name="` + fieldId + `" 
				placeholder="Your ` + placeholder + `" 
				value="` + value + `"
			>
		</fieldset>
		<div class="field-error-message">` + message + `</div>
	</div>`
}

func LoginValidation(form utils.LoginForm) (bool, string, string) {
	if form.Identificator == "" {
		return false, "identificator", fieldset("identificator", "Identificator is required", form.Identificator, "text", "Identificator", true)
	}
	if form.Password == "" {
		return false, "password", fieldset("password", "Password is required", form.Password, "password", "Password", true)
	}
	return true, "", ""
}

func RegisterValidation(form utils.RegisterForm) (bool, string, string) {
	if form.Username == "" {
		return false, "username", fieldset("username", "Username is required", form.Username, "text", "Username", true)
	}
	if form.Email == "" {
		return false, "email", fieldset("email", "Email is required", form.Email, "email", "Email", true)
	}
	if form.Password == "" {
		return false, "password", fieldset("password", "Password is required", form.Password, "password", "Password", true)
	}
	if form.Confirm == "" {
		return false, "confirm-password", fieldset("confirm-password", "Please confirm your password", form.Confirm, "password", "Confirm Password", true)
	}
	if form.Password != form.Confirm {
		return false, "confirm-password", fieldset("confirm-password", "Passwords do not match", form.Confirm, "password", "Confirm Password", true)
	}
	return true, "", ""
}
