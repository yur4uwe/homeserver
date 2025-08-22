package scripts

import (
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks/mocksos"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestChangeClassNames(t *testing.T) {
	mockOS := new(mocksos.MockOS)
	cleanupReadFile := helpers.MockFunction(&ReadFile, mockOS.ReadFile, func() {})
	cleanupReadDir := helpers.MockFunction(&ReadDir, mockOS.ReadDir, func() {})
	cleanupWriteFile := helpers.MockFunction(&WriteFile, mockOS.WriteFile, func() {})

	mockSvg := []byte(`<svg><style>.st0{fill:#000;}</style></svg>`)

	mockDirEntry := new(mocksos.MockDirEntry)

	resetter := helpers.NewMocksResetter(
		&mockOS.Mock,
		&mockDirEntry.Mock,
	)

	defer func() {
		cleanupReadFile()
		cleanupReadDir()
		cleanupWriteFile()
	}()

	t.Run("Changes name of the class", func(t *testing.T) {
		resetter()

		mockDirEntry.On("Name").Return("test.svg")
		mockDirEntry.On("IsDir").Return(false)
		mockDirEntry.On("Type").Return(os.ModePerm)
		mockDirEntry.On("Info").Return((os.FileInfo)(nil), nil)

		mockOS.On("ReadFile", IconDir+"\\test.svg").Return(mockSvg, nil)
		mockOS.On("ReadDir", IconDir).Return([]os.DirEntry{mockDirEntry}, nil)
		mockOS.On("WriteFile", IconDir+"\\test.svg", mock.Anything, os.FileMode(0644)).Return(nil)

		ChangeClassNames()

		mockOS.AssertExpectations(t)
	})
}
