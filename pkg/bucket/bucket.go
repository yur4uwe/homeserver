package bucket

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Bucket defines the interface for file storage operations.
type Bucket interface {
	Save(path string, data io.Reader) error
	Delete(path string) error
	Move(oldPath, newPath string) error
	CreateDir(path string) error
	Open(path string) (io.ReadCloser, error)
	Exists(path string) (bool, error)
}

// FSBucket is a local filesystem implementation of the Bucket interface.
type FSBucket struct {
	basePath string
}

func NewFSBucket(basePath string) (*FSBucket, error) {
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("invalid base path: %w", err)
	}
	return &FSBucket{basePath: absBase}, nil
}

func (b *FSBucket) fullPath(relPath string) string {
	return filepath.Join(b.basePath, relPath)
}

func (b *FSBucket) Save(path string, data io.Reader) error {
	fullPath := b.fullPath(path)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, data)
	return err
}

func (b *FSBucket) Delete(path string) error {
	return os.RemoveAll(b.fullPath(path))
}

func (b *FSBucket) Move(oldPath, newPath string) error {
	return os.Rename(b.fullPath(oldPath), b.fullPath(newPath))
}

func (b *FSBucket) CreateDir(path string) error {
	return os.MkdirAll(b.fullPath(path), 0755)
}

func (b *FSBucket) Open(path string) (io.ReadCloser, error) {
	return os.Open(b.fullPath(path))
}

func (b *FSBucket) Exists(path string) (bool, error) {
	_, err := os.Stat(b.fullPath(path))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
