package folconsts

const (
	// General
	NoFieldsSpecifiedError  = "no fields specified for validation"
	InvalidFolderFieldError = "field '%s' is not a valid folder field"

	// Folder-specific
	InvalidIDError           = "ID must be greater than 0"
	MissingNameError         = "name cannot be empty"
	MissingPathError         = "path cannot be empty"
	InvalidPathEndError      = "path cannot end with a slash"
	InvalidPathStartError    = "path cannot start with a slash"
	MissingOwnerFolderError  = "path must start with the owner's name"
	InvalidOwnerIDError      = "owner ID must be greater than 0"
	InvalidFolderParentError = "parent ID must be greater than 0"
)
