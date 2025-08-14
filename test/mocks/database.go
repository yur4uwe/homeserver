package mocks

import (
	"context"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"

	"mime/multipart"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

// Basic Query Methods
func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgx.Rows), argsMock.Error(1)
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgx.Row)
}

func (m *MockDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgconn.CommandTag), argsMock.Error(1)
}

// Table Management
func (m *MockDB) CreateTable(ctx context.Context, table *globals.DBTable) error {
	args := m.Called(ctx, table)
	return args.Error(0)
}

func (m *MockDB) TableExists(ctx context.Context, tableName string) (bool, error) {
	args := m.Called(ctx, tableName)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) ClearEntries(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

// File Queries
func (m *MockDB) ExtractUserFiles(ctx context.Context, userID int) ([]*file_schema.UserFile, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*file_schema.UserFile), args.Error(1)
}

func (m *MockDB) SearchFiles(ctx context.Context, query string, requester *user_schema.User, directory []string) ([]*file_schema.UserFile, error) {
	args := m.Called(ctx, query, requester, directory)
	return args.Get(0).([]*file_schema.UserFile), args.Error(1)
}

func (m *MockDB) CreateFileEntry(ctx context.Context, file *file_schema.UserFile, requester *user_schema.User) (int, error) {
	args := m.Called(ctx, file, requester)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) ExtractFilesFromDirectory(ctx context.Context, directory string, requester *user_schema.User) ([]*file_schema.UserFile, error) {
	args := m.Called(ctx, directory, requester)
	return args.Get(0).([]*file_schema.UserFile), args.Error(1)
}

func (m *MockDB) ExtractAllFiles(ctx context.Context) ([]*file_schema.UserFile, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*file_schema.UserFile), args.Error(1)
}

func (m *MockDB) GetFile(ctx context.Context, file_data *file_schema.UserFile, requester *user_schema.User) (*file_schema.UserFile, error) {
	args := m.Called(ctx, file_data, requester)
	return args.Get(0).(*file_schema.UserFile), args.Error(1)
}

func (m *MockDB) DeleteFile(ctx context.Context, file *file_schema.UserFile, userName string) (int, error) {
	args := m.Called(ctx, file, userName)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) RenameFile(ctx context.Context, oldFile *file_schema.UserFile, newPath string, requester *user_schema.User) (string, error) {
	args := m.Called(ctx, oldFile, newPath, requester)
	return args.String(0), args.Error(1)
}

func (m *MockDB) UpdateFileEntry(ctx context.Context, file *file_schema.UserFile, ownerName string) error {
	args := m.Called(ctx, file, ownerName)
	return args.Error(0)
}

func (m *MockDB) SaveFile(ctx context.Context, header *file_schema.UserFile, fileData *multipart.File, requester *user_schema.User) (int, string, error) {
	args := m.Called(ctx, header, fileData, requester)
	return args.Int(0), args.String(1), args.Error(2)
}

// Folder Queries
func (m *MockDB) CreateFolderEntry(ctx context.Context, folder *folder_schema.UserFolder, ownerName string) (*folder_schema.UserFolder, error) {
	args := m.Called(ctx, folder, ownerName)
	return args.Get(0).(*folder_schema.UserFolder), args.Error(1)
}

func (m *MockDB) CreateRootFolderEntry(ctx context.Context, folder *folder_schema.UserFolder, ownerName string) (*folder_schema.UserFolder, error) {
	args := m.Called(ctx, folder, ownerName)
	return args.Get(0).(*folder_schema.UserFolder), args.Error(1)
}

func (m *MockDB) ExtractUserFolders(ctx context.Context, userID int) ([]*folder_schema.UserFolder, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*folder_schema.UserFolder), args.Error(1)
}

func (m *MockDB) QueryFolderByName(ctx context.Context, folderName string, ownerID int, pathToFolder string) (*folder_schema.UserFolder, error) {
	args := m.Called(ctx, folderName, ownerID, pathToFolder)
	return args.Get(0).(*folder_schema.UserFolder), args.Error(1)
}

func (m *MockDB) FolderExists(ctx context.Context, folderName string, ownerID int, pathToFolder string) (bool, error) {
	args := m.Called(ctx, folderName, ownerID, pathToFolder)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) ExtractFoldersFromDirectory(ctx context.Context, directory []string, requester *user_schema.User) ([]*folder_schema.UserFolder, error) {
	args := m.Called(ctx, directory, requester)
	return args.Get(0).([]*folder_schema.UserFolder), args.Error(1)
}

func (m *MockDB) DeleteFolder(ctx context.Context, folder *folder_schema.UserFolder, ownerName string) (int, error) {
	args := m.Called(ctx, folder, ownerName)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) EnsureFolderExists(ctx context.Context, relativePath []string, owner *user_schema.User) error {
	args := m.Called(ctx, relativePath, owner)
	return args.Error(0)
}

func (m *MockDB) UpdateFolder(ctx context.Context, folder_data *folder_schema.UserFolder, ownerName string) (*folder_schema.UserFolder, error) {
	args := m.Called(ctx, folder_data, ownerName)
	return args.Get(0).(*folder_schema.UserFolder), args.Error(1)
}

func (m *MockDB) RenameFolder(ctx context.Context, folder_data *folder_schema.UserFolder, newName string, ownerName string) error {
	args := m.Called(ctx, folder_data, newName, ownerName)
	return args.Error(0)
}

// Token Queries
func (m *MockDB) CheckToken(ctx context.Context, token_str string) (int, error) {
	args := m.Called(ctx, token_str)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) CreateToken(ctx context.Context, userID int, token_str string) error {
	args := m.Called(ctx, userID, token_str)
	return args.Error(0)
}

func (m *MockDB) QueryToken(ctx context.Context, token_str string) (*user_schema.User, error) {
	args := m.Called(ctx, token_str)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_schema.User), args.Error(1)
}

func (m *MockDB) VerifyToken(ctx context.Context, token_str string) error {
	args := m.Called(ctx, token_str)
	return args.Error(0)
}

func (m *MockDB) DeleteToken(ctx context.Context, token_str string) error {
	args := m.Called(ctx, token_str)
	return args.Error(0)
}

// User Queries
func (m *MockDB) CreateUserEntry(ctx context.Context, user *user_schema.User) (*user_schema.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*user_schema.User), args.Error(1)
}

func (m *MockDB) ExtractUserEntries(ctx context.Context) ([]*user_schema.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*user_schema.User), args.Error(1)
}

func (m *MockDB) CheckUserFolder(ctx context.Context, user *user_schema.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockDB) CreateUser(ctx context.Context, username, email, password string) (*user_schema.User, error) {
	args := m.Called(ctx, username, email, password)
	return args.Get(0).(*user_schema.User), args.Error(1)
}

func (m *MockDB) QueryUserOnLogin(ctx context.Context, form globals.LoginForm) (*user_schema.User, bool, error) {
	args := m.Called(ctx, form)
	return args.Get(0).(*user_schema.User), args.Bool(1), args.Error(2)
}

func (m *MockDB) QueryUsername(ctx context.Context, username string) (*user_schema.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*user_schema.User), args.Error(1)
}

func (m *MockDB) QueryEmail(ctx context.Context, email string) (*user_schema.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*user_schema.User), args.Error(1)
}

func (m *MockDB) QueryPassword(ctx context.Context, password string) (*user_schema.User, error) {
	args := m.Called(ctx, password)
	return args.Get(0).(*user_schema.User), args.Error(1)
}

func (m *MockDB) DropTable(ctx context.Context, tableName string) error {
	args := m.Called(ctx, tableName)
	return args.Error(0)
}

func (m *MockDB) ModifyUserField(ctx context.Context, user *user_schema.User, field string, value interface{}) error {
	args := m.Called(ctx, user, field, value)
	return args.Error(0)
}
