package auth

const (
	IdentificatorField   = "identificator"
	PasswordField        = "password"
	UsernameField        = "username"
	EmailField           = "email"
	ConfirmPasswordField = "confirm-password"
)

const (
	IdentificatorPlaceholder   = "Username or Email"
	PasswordPlaceholder        = "Password"
	UsernamePlaceholder        = "Username"
	EmailPlaceholder           = "Your Email"
	ConfirmPasswordPlaceholder = "Confirm Password"
	LoginFieldDesignation      = "login"
	RegisterFieldDesignation   = "register"
	IdentificatorError         = "Identificator is required"
	PasswordError              = "Password is required"
	UsernameError              = "Username is required"
	EmailError                 = "Email is required"
	ConfirmPasswordError       = "Confirm Password is required"
	PasswordMismatchError      = "Passwords do not match"
)

const (
	Login    = "login"
	Register = "register"
)
