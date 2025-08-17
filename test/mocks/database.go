package mocks

import (
	repomocks "homeserver/test/mocks/repos"
)

func InitRepoMocks() (
	*repomocks.MockFileRepo, *repomocks.MockUserRepo,
	*repomocks.MockFolderRepo, *repomocks.MockTokenRepo,
	*repomocks.MockTableRepo,
) {
	return repomocks.NewMockFileRepo(),
		repomocks.NewMockUserRepo(),
		repomocks.NewMockFolderRepo(),
		repomocks.NewMockTokenRepo(),
		repomocks.NewMockTableRepo()
}
