package fconsts

const (
	// General
	NoFieldsSpecifiedError = "no fields specified for validation"
	InvalidFileFieldError  = "field '%s' is not a valid file field"

	// File-specific
	InvalidIDError          = "ID must be greater than 0"
	MissingNameError        = "name cannot be empty"
	MissingPathError        = "path cannot be empty"
	InvalidPathEndError     = "path cannot end with a slash"
	InvalidPathStartError   = "path cannot start with a slash"
	MissingOwnerFolderError = "path must start with the owner's name"
	InvalidOwnerIDError     = "owner ID must be greater than 0"
	InvalidExtensionError   = "extension cannot start with a dot"
	InvalidDirectoryError   = "directory ID must be greater than 0"
	InvalidSizeError        = "file size must be greater than 0"
	InvalidContentTypeError = "content type cannot be empty"
	MissingMetadataError    = "file metadata is missing"
)
