package mocks

import (
	repomocks "homeserver/test/mocks/repos"
)

func InitRepoMocks() (
	*repomocks.MockUserRepo, *repomocks.MockFileRepo,
	*repomocks.MockFolderRepo, *repomocks.MockTokenRepo,
	*repomocks.MockTableRepo,
) {
	return repomocks.NewMockUserRepo(),
		repomocks.NewMockFileRepo(),
		repomocks.NewMockFolderRepo(),
		repomocks.NewMockTokenRepo(),
		repomocks.NewMockTableRepo()
}
