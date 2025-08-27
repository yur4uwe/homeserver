package auth

import (
	"homeserver/constants/authconsts"
	"homeserver/pkg/utils/formatters"
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

func LoginValidation(form formatters.LoginForm) (bool, string, string) {
	if form.Identificator == "" {
		return false, authconsts.IdentificatorField, fieldset(authconsts.IdentificatorField, authconsts.IdentificatorError, form.Identificator, "text", authconsts.IdentificatorPlaceholder, true)
	}
	if form.Password == "" {
		return false, authconsts.PasswordField, fieldset(authconsts.PasswordField, authconsts.PasswordError, form.Password, "password", authconsts.PasswordPlaceholder, true)
	}
	return true, "", ""
}

func RegisterValidation(form formatters.RegisterForm) (bool, string, string) {
	if form.Username == "" {
		return false, authconsts.UsernameField, fieldset(authconsts.UsernameField, authconsts.UsernameError, form.Username, "text", authconsts.UsernamePlaceholder, true)
	}
	if form.Email == "" {
		return false, authconsts.EmailField, fieldset(authconsts.EmailField, authconsts.EmailError, form.Email, "email", authconsts.EmailPlaceholder, true)
	}
	if form.Password == "" {
		return false, authconsts.PasswordField, fieldset(authconsts.PasswordField, authconsts.PasswordError, form.Password, "password", authconsts.PasswordPlaceholder, true)
	}
	if form.Confirm == "" {
		return false, authconsts.ConfirmPasswordField, fieldset(authconsts.ConfirmPasswordField, authconsts.ConfirmPasswordError, form.Confirm, "password", authconsts.ConfirmPasswordPlaceholder, true)
	}
	if form.Password != form.Confirm {
		return false, authconsts.ConfirmPasswordField, fieldset(authconsts.ConfirmPasswordField, authconsts.PasswordMismatchError, form.Confirm, "password", authconsts.ConfirmPasswordPlaceholder, true)
	}
	return true, "", ""
}
