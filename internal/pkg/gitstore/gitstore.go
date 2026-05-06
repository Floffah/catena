package gitstore

import (
	"context"
	"errors"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"
	"time"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/environment"
	"github.com/floffah/catena/internal/pkg/git"
)

const maxReadmeBytes = 1024 * 1024

var (
	ErrInvalidPath    = errors.New("invalid path")
	ErrInvalidRef     = errors.New("invalid ref")
	ErrReadmeNotFound = errors.New("readme not found")
	ErrReadmeTooLarge = errors.New("readme too large")
	ErrCommitNotFound = errors.New("commit not found")
	ErrTreeNotFound   = errors.New("tree not found")
)

type Readme struct {
	Ref       string
	CommitOID string
	Path      string
	Name      string
	OID       string
	Size      int64
	Content   string
}

type Tree struct {
	Ref       string
	CommitOID string
	Path      string
	Entries   []TreeEntry
}

type TreeEntry struct {
	Name string
	Path string
	Type string
	OID  string
	Size *int64
}

type LatestCommit struct {
	Ref             string
	CommitOID       string
	ShortOID        string
	MessageHeadline string
	Message         string
	AuthorName      string
	AuthorEmail     string
	AuthoredAt      time.Time
	CommitterName   string
	CommitterEmail  string
	CommittedAt     time.Time
}

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

	err := os.MkdirAll(filepath.Dir(repoPath), 0755)
	if err != nil {
		return err
	}

	err = s.git.InitBare(repoPath, dbRepo.DefaultBranch)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) GetRepoPath(dbRepo db.Repository) string {
	byte1 := dbRepo.ID.String()[:2]
	byte2 := dbRepo.ID.String()[2:4]

	return filepath.Join(s.root, string(byte1), string(byte2), dbRepo.ID.String()+".git")
}

func (s Store) GitBinaryPath() string {
	return s.git.Path()
}

func (s Store) GetReadme(ctx context.Context, dbRepo db.Repository, ref string, directory string) (Readme, error) {
	if ref == "" {
		ref = dbRepo.DefaultBranch
	}
	if !isSafeRef(ref) {
		return Readme{}, ErrInvalidRef
	}

	directory, err := normalizeGitDirectory(directory)
	if err != nil {
		return Readme{}, err
	}

	repoPath := s.GetRepoPath(dbRepo)
	commitOID, err := s.git.ResolveCommit(ctx, repoPath, ref)
	if err != nil {
		return Readme{}, ErrReadmeNotFound
	}

	for _, name := range []string{"README.md", "README"} {
		readmePath := name
		if directory != "" {
			readmePath = pathpkg.Join(directory, name)
		}

		entry, err := s.git.LsTreePath(ctx, repoPath, commitOID, readmePath)
		if err != nil {
			return Readme{}, err
		}
		if entry == nil || entry.Type != "blob" {
			continue
		}
		if entry.Size == nil {
			return Readme{}, ErrReadmeNotFound
		}
		if *entry.Size > maxReadmeBytes {
			return Readme{}, ErrReadmeTooLarge
		}

		content, err := s.git.CatFileBlob(ctx, repoPath, entry.OID)
		if err != nil {
			return Readme{}, err
		}

		return Readme{
			Ref:       ref,
			CommitOID: commitOID,
			Path:      readmePath,
			Name:      name,
			OID:       entry.OID,
			Size:      *entry.Size,
			Content:   string(content),
		}, nil
	}

	return Readme{}, ErrReadmeNotFound
}

func (s Store) GetTree(ctx context.Context, dbRepo db.Repository, ref string, directory string) (Tree, error) {
	if ref == "" {
		ref = dbRepo.DefaultBranch
	}
	if !isSafeRef(ref) {
		return Tree{}, ErrInvalidRef
	}

	directory, err := normalizeGitDirectory(directory)
	if err != nil {
		return Tree{}, err
	}

	repoPath := s.GetRepoPath(dbRepo)
	commitOID, err := s.git.ResolveCommit(ctx, repoPath, ref)
	if err != nil {
		return Tree{}, ErrTreeNotFound
	}

	treeish := commitOID
	if directory != "" {
		treeish = commitOID + ":" + directory
	}

	gitEntries, err := s.git.LsTree(ctx, repoPath, treeish)
	if err != nil {
		return Tree{}, ErrTreeNotFound
	}

	entries := make([]TreeEntry, 0, len(gitEntries))
	for _, entry := range gitEntries {
		entryPath := entry.Path
		if directory != "" {
			entryPath = pathpkg.Join(directory, entry.Path)
		}

		entries = append(entries, TreeEntry{
			Name: pathpkg.Base(entry.Path),
			Path: entryPath,
			Type: entry.Type,
			OID:  entry.OID,
			Size: entry.Size,
		})
	}

	return Tree{
		Ref:       ref,
		CommitOID: commitOID,
		Path:      directory,
		Entries:   entries,
	}, nil
}

func (s Store) GetLatestCommit(ctx context.Context, dbRepo db.Repository, ref string, path string) (LatestCommit, error) {
	if ref == "" {
		ref = dbRepo.DefaultBranch
	}
	if !isSafeRef(ref) {
		return LatestCommit{}, ErrInvalidRef
	}

	path, err := normalizeGitDirectory(path)
	if err != nil {
		return LatestCommit{}, err
	}

	repoPath := s.GetRepoPath(dbRepo)
	commit, err := s.git.LogLatest(ctx, repoPath, ref, path)
	if err != nil {
		return LatestCommit{}, ErrCommitNotFound
	}
	if commit == nil {
		return LatestCommit{}, ErrCommitNotFound
	}

	return LatestCommit{
		Ref:             ref,
		CommitOID:       commit.OID,
		ShortOID:        commit.ShortOID,
		MessageHeadline: commit.MessageHeadline,
		Message:         commit.Message,
		AuthorName:      commit.AuthorName,
		AuthorEmail:     commit.AuthorEmail,
		AuthoredAt:      commit.AuthoredAt,
		CommitterName:   commit.CommitterName,
		CommitterEmail:  commit.CommitterEmail,
		CommittedAt:     commit.CommittedAt,
	}, nil
}

func normalizeGitDirectory(directory string) (string, error) {
	directory = strings.TrimSpace(strings.ReplaceAll(directory, "\\", "/"))
	if directory == "" || directory == "." || directory == "/" {
		return "", nil
	}

	if strings.ContainsAny(directory, "\x00\r\n") {
		return "", ErrInvalidPath
	}

	segments := make([]string, 0)
	for _, segment := range strings.Split(directory, "/") {
		switch segment {
		case "", ".":
			continue
		case "..":
			return "", ErrInvalidPath
		default:
			segments = append(segments, segment)
		}
	}

	if len(segments) == 0 {
		return "", nil
	}

	return pathpkg.Join(segments...), nil
}

func isSafeRef(ref string) bool {
	ref = strings.TrimSpace(ref)
	if ref == "" || strings.HasPrefix(ref, "-") {
		return false
	}

	return !strings.ContainsAny(ref, "\x00\r\n")
}
