package folderrepo

import (
	"homeserver/constants/folconsts"
	folschema "homeserver/db/repositories/folder_repo/schema"
)

func OrderUpdateFields(fields_to_update []string, folder *folschema.UserFolder) []any {
	var args = make([]any, 0, len(fields_to_update))
	for _, field := range fields_to_update {
		switch field {
		case folconsts.FolderName:
			args = append(args, folder.Name)
		case folconsts.Path:
			args = append(args, folder.Path)
		case folconsts.Parent:
			args = append(args, folder.Parent)
		}
	}
	return args
}
