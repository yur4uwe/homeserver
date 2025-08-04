package schemas

const (
	UserID        = "id"
	Username      = "username"
	Password      = "password"
	Email         = "email"
	AuthToken     = "auth_token"
	EmailVerified = "email_verified"
	CreatedAt     = "created_at"
	FileID        = "id"
	Filename      = "filename"
	Extension     = "extension"
	Path          = "path"
	Directory     = "directory_id"
	Size          = "size"
	ContentType   = "content_type"
	OwnerID       = "owner_id"
	UpdatedAt     = "updated_at"
	Metadata      = "metadata"
	FolderID      = "id"
	FolderName    = "name"
	FolderParent  = "parent"
)

const (
	NoFieldsSpecifiedError  = "no fields specified for validation"
	ValidationErrorPrefix   = "validation errors: "
	InvalidUserFieldError   = "field '%s' is not a valid user field"
	InvalidFileFieldError   = "field '%s' is not a valid file field"
	InvalidFolderFieldError = "field '%s' is not a valid folder field"
)

// UserFolder error messages
const (
	InvalidFolderIDError        = "folder ID must be greater than 0"
	InvalidFolderNameError      = "folder name cannot be empty"
	InvalidFolderPathError      = "path cannot be empty"
	InvalidFolderPathEndError   = "path cannot end with a slash"
	InvalidFolderPathStartError = "path cannot start with a slash"
	InvalidFolderPathOwnerError = "path must start with the owner's name '%s'"
	InvalidFolderParentError    = "parent ID must be greater than 0"
	InvalidFolderOwnerIDError   = "owner ID must be greater than 0"
)

// UserFile error messages
const (
	InvalidFileIDError        = "file ID must be greater than 0"
	InvalidFilenameError      = "filename cannot be empty"
	InvalidExtensionError     = "extension cannot start with a dot"
	InvalidFilePathError      = "path cannot start with a slash"
	InvalidFilePathEndError   = "path cannot end with a slash"
	InvalidFilePathOwnerError = "path must start with the owner's name '%s'"
	InvalidDirectoryError     = "directory ID must be greater than 0"
	InvalidFileSizeError      = "file size must be greater than 0"
	InvalidContentTypeError   = "content type cannot be empty"
	InvalidFileOwnerIDError   = "owner ID must be greater than 0"
	InvalidFileMetadataError  = "metadata cannot be empty"
)

// User error messages
const (
	InvalidUserIDError        = "user ID must be greater than 0"
	InvalidUsernameError      = "username cannot be empty"
	InvalidPasswordError      = "password cannot be empty"
	InvalidEmailError         = "email cannot be empty"
	InvalidEmailFormatError   = "email must contain '@'"
	InvalidEmailVerifiedError = "email is not verified"
	InvalidAuthTokenError     = "auth token cannot be empty"
	InvalidCreatedAtError     = "created at timestamp cannot be empty"
	InvalidUpdatedAtError     = "updated at timestamp cannot be empty"
)
