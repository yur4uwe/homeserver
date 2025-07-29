package main

import (
	"context"
	"fmt"
	"homeserver/pkg/db"
	"strings"
)

func SeparateExtensionsFromNames(connection *db.Database) error {
	ctx := context.Background()

	files, err := connection.ExtractAllFiles(ctx)
	if err != nil {
		return fmt.Errorf("error extracting files: %v", err)
	}

	for _, file := range files {
		file_sig := strings.Split(file.Filename, ".")
		if len(file_sig) > 1 && file.Extension == "" {
			file.Filename = strings.Join(file_sig[:len(file_sig)-1], ".")
			file.Extension = file_sig[len(file_sig)-1]
		}

		if err := connection.UpdateFileEntry(ctx, file); err != nil {
			return fmt.Errorf("error updating file entry: %v", err)
		}
	}

	return nil
}
