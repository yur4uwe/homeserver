package zipper

import (
	"archive/zip"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func ZipFile(zw *zip.Writer, filePath, zipPath string) error {
	fileToZip, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	zipFile, err := zw.Create(zipPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFile, fileToZip)
	return err
}

func ZipFolder(w http.ResponseWriter, folder *folschema.UserFolder) error {
	zw := zip.NewWriter(w)
	defer zw.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		zipPath, _ := filepath.Rel(folder.AbsolutePath(), path)
		return ZipFile(zw, path, zipPath)
	}

	if err := filepath.Walk(folder.AbsolutePath(), walker); err != nil {
		return err
	}

	return nil
}
