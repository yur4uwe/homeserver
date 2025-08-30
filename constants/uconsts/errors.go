package uconsts

const (
	// General
	NoFieldsSpecifiedError = "no fields specified for validation"
	InvalidUserFieldError  = "field '%s' is not a valid user field"

	// User-specific
	InvalidIDError            = "ID must be greater than 0"
	MissingUsernameError      = "username cannot be empty"
	InvalidPasswordError      = "password cannot be empty"
	InvalidEmailError         = "email cannot be empty"
	InvalidEmailFormatError   = "email must contain '@'"
	InvalidEmailVerifiedError = "email is not verified"
	InvalidAuthTokenError     = "auth token cannot be empty"
	InvalidCreatedAtError     = "created at cannot be empty"
)
