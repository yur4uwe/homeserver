package scripts

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	ReadDir   = os.ReadDir
	ReadFile  = os.ReadFile
	WriteFile = os.WriteFile
)

const (
	IconDir = "G:\\projects\\Myprojects\\homeserver\\static\\assets\\icons"
)

func ChangeClassNames() {
	files, err := ReadDir(IconDir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".svg" {
			continue
		}

		svg_file, err := ReadFile(IconDir + "\\" + file.Name())
		if err != nil {
			fmt.Println("Error reading file:", file.Name(), err)
			continue
		}

		svg := string(svg_file)

		for i := 0; i < len(svg)-3; i++ {
			if svg[i:i+3] != "st0" {
				continue
			}

			svg = svg[:i] + "st" + file.Name()[:len(file.Name())-4] + svg[i+3:]
		}

		if err := WriteFile(IconDir+"\\"+file.Name(), []byte(svg), 0644); err != nil {
			fmt.Println("Error writing file:", file.Name(), err)
			continue
		}
	}
}
