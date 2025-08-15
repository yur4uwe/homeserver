package filerepo

import (
	"fmt"
	fschema "homeserver/db/repositories/file_repo/schema"
	"os"
)

var (
	Stat = os.Stat
)

func PreventDuplicateFileName(file *fschema.UserFile) (*fschema.UserFile, error) {
	new_file := *file
	counter := 1
	for {
		path := new_file.AbsolutePath()

		if _, err := Stat(path); err != nil && os.IsNotExist(err) {
			return &new_file, nil
		}
		fmt.Printf("File %s already exists, trying to prevent duplicate...\n", new_file.AbsolutePath())

		new_file.Filename = fmt.Sprintf("%s_%d", new_file.Filename, counter)

		counter++
	}
}
