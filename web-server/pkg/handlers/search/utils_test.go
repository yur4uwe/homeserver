package search

import (
	fschema "homeserver/db/repositories/file_repo/schema"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/test/helpers"
	"homeserver/web-server/pkg/templates"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIsFileTypeAvailable struct {
	mock.Mock
}

func (m *MockIsFileTypeAvailable) IsFileTypeAvailable(extension string) bool {
	args := m.Called(extension)
	return args.Bool(0)
}

func TestParseQuerySlice(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"empty", []string{}, ""},
		{"single", []string{"foo"}, "foo"},
		{"multiple", []string{"foo", "bar"}, "foo/bar"},
		{"needs escaping", []string{"foo bar", "baz/qux"}, "foo+bar/baz%2Fqux"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseQuerySlice(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want string
	}{
		{"bytes", 512, "512.00 B"},
		{"kilobytes", 2048, "2.00 KB"},
		{"megabytes", 5 * 1024 * 1024, "5.00 MB"},
		{"gigabytes", 3 * 1024 * 1024 * 1024, "3.00 GB"},
		{"terabytes", 2 * 1024 * 1024 * 1024 * 1024, "2.00 TB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSize(tt.size)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFileToSearchEntry(t *testing.T) {
	mockIsFileTypeAvailable := new(MockIsFileTypeAvailable)
	cleanup := helpers.MockFunction(&IsFileTypeAvailable, mockIsFileTypeAvailable.IsFileTypeAvailable, func() {})
	defer cleanup()

	mockIsFileTypeAvailable.On("IsFileTypeAvailable", "pdf").Return(true)

	file := &fschema.UserFile{
		ID:        1,
		Filename:  "report",
		Extension: "pdf",
		Path:      "user/docs",
		Size:      2048,
	}

	got := FileToSearchEntry(file)
	assert.Equal(t, 1, got.ID)
	assert.Equal(t, "report", got.Name)
	assert.Equal(t, "docs/", got.Path)
	assert.Equal(t, "2.00 KB", got.Size)
	assert.Equal(t, ".pdf", got.Extension)
	assert.Equal(t, "docs/", got.URIPath)
	assert.Equal(t, "report", got.URIName)
	assert.False(t, got.IsFolder)
	assert.Equal(t, "icons/pdf", got.Type)

	file.Extension = "unknown"

	mockIsFileTypeAvailable.On("IsFileTypeAvailable", "unknown").Return(false)

	got = FileToSearchEntry(file)
	assert.Equal(t, 1, got.ID)
	assert.Equal(t, "report", got.Name)
	assert.Equal(t, "docs/", got.Path)
	assert.Equal(t, "2.00 KB", got.Size)
	assert.Equal(t, ".unknown", got.Extension)
	assert.Equal(t, "docs/", got.URIPath)
	assert.Equal(t, "report", got.URIName)
	assert.False(t, got.IsFolder)
	assert.Equal(t, "file", got.Type)

	mockIsFileTypeAvailable.AssertExpectations(t)
}

func TestFolderToSearchEntry(t *testing.T) {
	folder := &folschema.UserFolder{
		ID:   2,
		Name: "reports",
		Path: "user/docs/reports",
	}
	got := FolderToSearchEntry(folder)
	assert.Equal(t, 2, got.ID)
	assert.Equal(t, "reports", got.Name)
	assert.Equal(t, "docs/reports/", got.Path)
	assert.Equal(t, "docs/reports/", got.URIPath)
	assert.Equal(t, "reports", got.URIName)
	assert.True(t, got.IsFolder)
	assert.Equal(t, "folder", got.Type)
	assert.Equal(t, "", got.Size)
}

func TestTransformEntries(t *testing.T) {
	mockIsFileTypeAvailable := new(MockIsFileTypeAvailable)
	cleanup := helpers.MockFunction(&IsFileTypeAvailable, mockIsFileTypeAvailable.IsFileTypeAvailable, func() {})
	defer cleanup()

	mockIsFileTypeAvailable.On("IsFileTypeAvailable", "txt").Return(true)
	mockIsFileTypeAvailable.On("IsFileTypeAvailable", mock.AnythingOfType("string")).Return(false)

	file := fschema.UserFile{ID: 1, Filename: "file", Extension: "txt", Path: "user/folder/file.txt", Size: 100}
	folder := folschema.UserFolder{ID: 2, Name: "folder", Path: "user/folder"}
	files := []fschema.UserFile{file}
	folders := []folschema.UserFolder{folder}
	entries := TransformEntries(files, folders)
	assert.Len(t, entries, 2)
	assert.Equal(t, "file", entries[0].Name)
	assert.Equal(t, "folder", entries[1].Name)
}

func TestSortEntries(t *testing.T) {
	entries := []templates.SearchEntry{
		{Name: "file2", IsFolder: false},
		{Name: "folder1", IsFolder: true},
		{Name: "file1", IsFolder: false},
		{Name: "folder2", IsFolder: true},
	}

	SortEntries(entries)

	assert.Equal(t, "folder1", entries[0].Name)
	assert.Equal(t, "folder2", entries[1].Name)
	assert.Equal(t, "file1", entries[2].Name)
	assert.Equal(t, "file2", entries[3].Name)
}

func TestExcludeFilesFromParents(t *testing.T) {
	files := []fschema.UserFile{
		{
			Filename: "file1.txt",
			Path:     "root/parent",
		},
		{
			Filename: "file2.txt",
			Path:     "root/parent/folder",
		},
		{
			Filename: "file3.txt",
			Path:     "root/parent/folder/subfolder",
		},
		{
			Filename: "file4.txt",
			Path:     "root",
		},
		{
			Filename: "file5.txt",
			Path:     "root/parent2",
		},
	}

	directory := []string{"root", "parent", "folder"}

	filtered := ExcludeFilesFromParents(files, directory)

	// Only files whose path starts with "root/parent/folder" should remain
	expected := []fschema.UserFile{
		{
			Filename: "file2.txt",
			Path:     "root/parent/folder",
		},
		{
			Filename: "file3.txt",
			Path:     "root/parent/folder/subfolder",
		},
	}

	assert.Equal(t, len(expected), len(filtered))
	for i, f := range filtered {
		assert.Equal(t, expected[i].Filename, f.Filename)
		assert.Equal(t, expected[i].Path, f.Path)
	}
}
