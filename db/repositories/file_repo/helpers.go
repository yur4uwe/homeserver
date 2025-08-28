package filerepo

import (
	"homeserver/constants/fconsts"
	fschema "homeserver/db/repositories/file_repo/schema"
	"os"
)

var (
	Stat = os.Stat
)

// OrderFileUpdateFields returns the values for the given fields in the correct order for SQL updates.
func OrderFileUpdateFields(fields_to_update []string, file *fschema.UserFile) []any {
	var args = make([]any, 0, len(fields_to_update))
	for _, field := range fields_to_update {
		switch field {
		case fconsts.ID:
			args = append(args, file.ID)
		case fconsts.Filename:
			args = append(args, file.Filename)
		case fconsts.Extension:
			args = append(args, file.Extension)
		case fconsts.Path:
			args = append(args, file.Path)
		case fconsts.Directory:
			args = append(args, file.Directory)
		case fconsts.Size:
			args = append(args, file.Size)
		case fconsts.ContentType:
			args = append(args, file.ContentType)
		case fconsts.OwnerID:
			args = append(args, file.OwnerID)
		case fconsts.CreatedAt:
			args = append(args, file.CreatedAt)
		case fconsts.UpdatedAt:
			args = append(args, file.UpdatedAt)
		case fconsts.Metadata:
			args = append(args, file.Metadata)
			// Add more fields as needed
		}
	}
	return args
}
