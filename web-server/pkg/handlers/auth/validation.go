package auth

import "homeserver/pkg/globals"

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

func LoginValidation(form globals.LoginForm) (bool, string, string) {
	if form.Identificator == "" {
		return false, IdentificatorField, fieldset(IdentificatorField, IdentificatorError, form.Identificator, "text", IdentificatorPlaceholder, true)
	}
	if form.Password == "" {
		return false, PasswordField, fieldset(PasswordField, PasswordError, form.Password, "password", PasswordPlaceholder, true)
	}
	return true, "", ""
}

func RegisterValidation(form globals.RegisterForm) (bool, string, string) {
	if form.Username == "" {
		return false, UsernameField, fieldset(UsernameField, UsernameError, form.Username, "text", UsernamePlaceholder, true)
	}
	if form.Email == "" {
		return false, EmailField, fieldset(EmailField, EmailError, form.Email, "email", EmailPlaceholder, true)
	}
	if form.Password == "" {
		return false, PasswordField, fieldset(PasswordField, PasswordError, form.Password, "password", PasswordPlaceholder, true)
	}
	if form.Confirm == "" {
		return false, ConfirmPasswordField, fieldset(ConfirmPasswordField, ConfirmPasswordError, form.Confirm, "password", ConfirmPasswordPlaceholder, true)
	}
	if form.Password != form.Confirm {
		return false, ConfirmPasswordField, fieldset(ConfirmPasswordField, PasswordMismatchError, form.Confirm, "password", ConfirmPasswordPlaceholder, true)
	}
	return true, "", ""
}
