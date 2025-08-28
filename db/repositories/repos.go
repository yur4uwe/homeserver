package repositories

import (
	filerepo "homeserver/db/repositories/file_repo"
	folderrepo "homeserver/db/repositories/folder_repo"
	tablerepo "homeserver/db/repositories/table_repo"
	tokenrepo "homeserver/db/repositories/token_repo"
	userrepo "homeserver/db/repositories/user_repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	UserRepo   userrepo.InUserRepo
	FileRepo   filerepo.InFileRepo
	FolderRepo folderrepo.InFolderRepo
	TokenRepo  tokenrepo.InTokenRepo
	TableRepo  tablerepo.InTableRepo
}

func NewRepositories(
	urepo userrepo.InUserRepo, frepo filerepo.InFileRepo,
	forepo folderrepo.InFolderRepo, tokrepo tokenrepo.InTokenRepo,
	tblrepo tablerepo.InTableRepo,
) *Repositories {
	return &Repositories{
		UserRepo:   urepo,
		FileRepo:   frepo,
		FolderRepo: forepo,
		TokenRepo:  tokrepo,
		TableRepo:  tblrepo,
	}
}

func InitRepositories(db *pgxpool.Pool) (userrepo.InUserRepo, filerepo.InFileRepo,
	folderrepo.InFolderRepo, tokenrepo.InTokenRepo,
	tablerepo.InTableRepo) {
	return userrepo.NewUserRepo(db),
		filerepo.NewFileRepo(db),
		folderrepo.NewFolderRepo(db),
		tokenrepo.NewTokenRepo(db),
		tablerepo.NewTableRepo(db)
}
