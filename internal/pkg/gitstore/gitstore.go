package gitstore

import (
	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/environment"
	"github.com/floffah/catena/internal/pkg/git"
)

// Store is not a git backend, but the git orchestrator. Business logic for git operations, using the git package as the backend
type Store struct {
	git git.Git

	root string
}

func NewStore(root string, gitBin string) Store {
	return Store{
		git:  git.NewGit(gitBin),
		root: root,
	}
}

func NewStoreFromEnv(env environment.Environment) Store {
	return NewStore(env.Config.CatenaGitRoot, env.GitBin)
}

func (s Store) CreateRepo(dbRepo db.Repository) error {
	repoPath := s.GetRepoPath(dbRepo)

	err := s.git.Init(repoPath)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) GetRepoPath(dbRepo db.Repository) string {
	byte1 := dbRepo.ID.String()[:2]
	byte2 := dbRepo.ID.String()[2:4]

	return s.root + "/" + string(byte1) + "/" + string(byte2) + "/" + dbRepo.ID.String() + ".git"
}
