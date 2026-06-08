package gitstore

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeebo/assert"
)

func TestNormalizeGitDirectory(t *testing.T) {
	tests := []struct {
		name      string
		directory string
		want      string
		wantErr   bool
	}{
		{
			name:      "empty path is root",
			directory: "",
			want:      "",
		},
		{
			name:      "slash path is root",
			directory: "/",
			want:      "",
		},
		{
			name:      "leading slash is repository relative",
			directory: "/docs",
			want:      "docs",
		},
		{
			name:      "nested path is preserved",
			directory: "docs/guides",
			want:      "docs/guides",
		},
		{
			name:      "dot segments are ignored",
			directory: "./docs/./guides",
			want:      "docs/guides",
		},
		{
			name:      "parent segment is rejected",
			directory: "../docs",
			wantErr:   true,
		},
		{
			name:      "nested parent segment is rejected",
			directory: "docs/../guides",
			wantErr:   true,
		},
		{
			name:      "nul byte is rejected",
			directory: "docs\x00guides",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeGitDirectory(tt.directory)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.That(t, got == tt.want)
		})
	}
}

func TestSortTreeEntries(t *testing.T) {
	size := int64(10)
	entries := []TreeEntry{
		{Name: "zeta.go", Path: "zeta.go", Type: "blob", Size: &size},
		{Name: "alpha.go", Path: "alpha.go", Type: "blob", Size: &size},
		{Name: "vendor", Path: "vendor", Type: "commit"},
		{Name: "Docs", Path: "Docs", Type: "tree"},
		{Name: "app", Path: "app", Type: "tree"},
		{Name: "README.md", Path: "README.md", Type: "blob", Size: &size},
	}

	sortTreeEntries(entries)

	got := make([]string, 0, len(entries))
	for _, entry := range entries {
		got = append(got, entry.Name)
	}

	want := []string{"app", "Docs", "vendor", "alpha.go", "README.md", "zeta.go"}
	assert.That(t, len(got) == len(want))

	for i := range want {
		assert.That(t, got[i] == want[i])
	}
}

func TestStoreCreateRepo(t *testing.T) {
	gitBin := requireGit(t)
	store, repository := newTestStore(t, gitBin)

	err := store.CreateRepo(repository)

	assert.Nil(t, err)
	assert.That(t, isDir(t, store.GetRepoPath(repository)))
	assert.That(t, fileContains(t, filepath.Join(store.GetRepoPath(repository), "HEAD"), "refs/heads/main"))
}

func TestStoreReadRepositoryContent(t *testing.T) {
	gitBin := requireGit(t)
	ctx := context.Background()
	store, repository := newPopulatedTestStore(t, gitBin)

	readme, err := store.GetReadme(ctx, repository, "", "")
	assert.Nil(t, err)
	assert.That(t, readme.Ref == "main")
	assert.That(t, readme.Name == "README.md")
	assert.That(t, readme.Path == "README.md")
	assert.That(t, readme.Content == "# Catena\n")

	nestedReadme, err := store.GetReadme(ctx, repository, "feature/deep-path", "docs")
	assert.Nil(t, err)
	assert.That(t, nestedReadme.Ref == "feature/deep-path")
	assert.That(t, nestedReadme.Name == "README")
	assert.That(t, nestedReadme.Path == "docs/README")
	assert.That(t, nestedReadme.Content == "Docs readme\n")

	file, err := store.GetFile(ctx, repository, "", "/src/main.go")
	assert.Nil(t, err)
	assert.That(t, file.Ref == "main")
	assert.That(t, file.Name == "main.go")
	assert.That(t, file.Path == "src/main.go")
	assert.That(t, file.Content == "package main\n")

	tree, err := store.GetTree(ctx, repository, "", "/", false)
	assert.Nil(t, err)
	assert.That(t, tree.Ref == "main")
	assert.That(t, tree.Path == "")
	assert.That(t, len(tree.Entries) == 3)
	assert.That(t, tree.Entries[0].Name == "docs")
	assert.That(t, tree.Entries[0].Type == "tree")
	assert.That(t, tree.Entries[1].Name == "src")
	assert.That(t, tree.Entries[1].Type == "tree")
	assert.That(t, tree.Entries[2].Name == "README.md")
	assert.That(t, tree.Entries[2].Type == "blob")

	recursiveTree, err := store.GetTree(ctx, repository, "", "/", true)
	assert.Nil(t, err)
	assert.That(t, recursiveTree.Ref == "main")
	assert.That(t, len(recursiveTree.Entries) == 6)
	assertTreeContains(t, recursiveTree, "docs", "tree")
	assertTreeContains(t, recursiveTree, "docs/README", "blob")
	assertTreeContains(t, recursiveTree, "docs/guide.md", "blob")
	assertTreeContains(t, recursiveTree, "src", "tree")
	assertTreeContains(t, recursiveTree, "src/main.go", "blob")
	assertTreeContains(t, recursiveTree, "README.md", "blob")
}

func TestStoreRecursiveTreeLimits(t *testing.T) {
	gitBin := requireGit(t)
	ctx := context.Background()
	store, repository := newPopulatedTestStore(t, gitBin)

	store.treeLimits.MaxEntries = 6
	tree, err := store.GetTree(ctx, repository, "", "/", true)
	assert.Nil(t, err)
	assert.That(t, len(tree.Entries) == store.treeLimits.MaxEntries)

	store.treeLimits.MaxEntries = 5
	_, err = store.GetTree(ctx, repository, "", "/", true)
	assert.That(t, errors.Is(err, ErrTreeTooLarge))

	store.treeLimits.MaxEntries = defaultTreeMaxEntries
	store.treeLimits.MaxBytes = 1
	_, err = store.GetTree(ctx, repository, "", "/", true)
	assert.That(t, errors.Is(err, ErrTreeTooLarge))

	store.treeLimits.MaxBytes = defaultTreeMaxBytes
	store.treeLimits.Timeout = time.Nanosecond
	_, err = store.GetTree(ctx, repository, "", "/", true)
	assert.That(t, errors.Is(err, ErrTreeTooLarge))
}

func TestStoreResolveGitPath(t *testing.T) {
	gitBin := requireGit(t)
	ctx := context.Background()
	store, repository := newPopulatedTestStore(t, gitBin)

	resolved, err := store.ResolveGitPath(ctx, repository, "main")
	assert.Nil(t, err)
	assert.That(t, resolved.Ref == "main")
	assert.That(t, resolved.Path == "")
	assert.That(t, resolved.PathType == "root")

	resolved, err = store.ResolveGitPath(ctx, repository, "feature/deep-path/docs/guide.md")
	assert.Nil(t, err)
	assert.That(t, resolved.Ref == "feature/deep-path")
	assert.That(t, resolved.Path == "docs/guide.md")
	assert.That(t, resolved.PathType == "blob")

	resolved, err = store.ResolveGitPath(ctx, repository, "feature/deep-path/docs")
	assert.Nil(t, err)
	assert.That(t, resolved.Ref == "feature/deep-path")
	assert.That(t, resolved.Path == "docs")
	assert.That(t, resolved.PathType == "tree")

	_, err = store.ResolveGitPath(ctx, repository, "missing/path")
	assert.That(t, errors.Is(err, ErrRefNotFound))
}

func TestStoreRefsAndBranchExists(t *testing.T) {
	gitBin := requireGit(t)
	ctx := context.Background()
	store, repository := newPopulatedTestStore(t, gitBin)

	refs, err := store.ListBranchRefs(ctx, repository)
	assert.Nil(t, err)
	assert.That(t, len(refs) == 2)
	assert.That(t, refs[0].Name == "main")
	assert.That(t, refs[0].IsDefault)
	assert.That(t, refs[1].Name == "feature/deep-path")
	assert.That(t, !refs[1].IsDefault)

	exists, err := store.BranchExists(ctx, repository, "feature/deep-path")
	assert.Nil(t, err)
	assert.That(t, exists)

	exists, err = store.BranchExists(ctx, repository, "v1.0.0")
	assert.Nil(t, err)
	assert.That(t, !exists)

	_, err = store.BranchExists(ctx, repository, "-bad-ref")
	assert.That(t, errors.Is(err, ErrInvalidRef))
}

func TestStoreLatestCommit(t *testing.T) {
	gitBin := requireGit(t)
	ctx := context.Background()
	store, repository := newPopulatedTestStore(t, gitBin)

	commit, err := store.GetLatestCommit(ctx, repository, "", "src/main.go")
	assert.Nil(t, err)
	assert.That(t, commit.Ref == "main")
	assert.That(t, commit.MessageHeadline == "Add source")
	assert.That(t, commit.AuthorName == "Catena Tests")
	assert.That(t, commit.AuthorEmail == "tests@catena.local")

	_, err = store.GetLatestCommit(ctx, repository, "", "missing.txt")
	assert.That(t, errors.Is(err, ErrCommitNotFound))
}

func TestStoreErrors(t *testing.T) {
	gitBin := requireGit(t)
	ctx := context.Background()
	store, repository := newPopulatedTestStore(t, gitBin)

	_, err := store.GetFile(ctx, repository, "", "/")
	assert.That(t, errors.Is(err, ErrInvalidPath))

	_, err = store.GetFile(ctx, repository, "-bad-ref", "README.md")
	assert.That(t, errors.Is(err, ErrInvalidRef))

	_, err = store.GetTree(ctx, repository, "", "../escape", false)
	assert.That(t, errors.Is(err, ErrInvalidPath))

	_, err = store.GetReadme(ctx, repository, "", "src")
	assert.That(t, errors.Is(err, ErrReadmeNotFound))
}

func assertTreeContains(t *testing.T, tree Tree, path string, entryType string) {
	t.Helper()

	for _, entry := range tree.Entries {
		if entry.Path == path {
			assert.That(t, entry.Type == entryType)
			return
		}
	}

	t.Fatalf("tree does not contain %q", path)
}

func requireGit(t *testing.T) string {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git binary is required for git-backed tests")
	}

	return gitBin
}

func newTestStore(t *testing.T, gitBin string) (Store, db.Repository) {
	t.Helper()

	id := uuid.New()
	return NewStore(t.TempDir(), gitBin), db.Repository{
		ID:            pgtype.UUID{Bytes: id, Valid: true},
		Name:          "catena",
		DefaultBranch: "main",
	}
}

func newPopulatedTestStore(t *testing.T, gitBin string) (Store, db.Repository) {
	t.Helper()

	store, repository := newTestStore(t, gitBin)
	assert.Nil(t, store.CreateRepo(repository))

	worktree := t.TempDir()
	runGit(t, worktree, "init", "--initial-branch=main")
	runGit(t, worktree, "config", "user.name", "Catena Tests")
	runGit(t, worktree, "config", "user.email", "tests@catena.local")
	runGit(t, worktree, "remote", "add", "origin", store.GetRepoPath(repository))

	writeTestFile(t, worktree, "README.md", "# Catena\n")
	writeTestFile(t, worktree, "docs/README", "Docs readme\n")
	writeTestFile(t, worktree, "docs/guide.md", "Guide\n")
	runGit(t, worktree, "add", ".")
	runGit(t, worktree, "commit", "-m", "Initial content")

	writeTestFile(t, worktree, "src/main.go", "package main\n")
	runGit(t, worktree, "add", ".")
	runGit(t, worktree, "commit", "-m", "Add source")
	runGit(t, worktree, "tag", "v1.0.0")
	runGit(t, worktree, "push", "origin", "main", "v1.0.0")

	runGit(t, worktree, "checkout", "-b", "feature/deep-path")
	writeTestFile(t, worktree, "docs/guide.md", "Feature guide\n")
	runGit(t, worktree, "add", ".")
	runGit(t, worktree, "commit", "-m", "Update guide")
	runGit(t, worktree, "push", "origin", "feature/deep-path")

	return store, repository
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Catena Tests",
		"GIT_AUTHOR_EMAIL=tests@catena.local",
		"GIT_COMMITTER_NAME=Catena Tests",
		"GIT_COMMITTER_EMAIL=tests@catena.local",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(output))
	}
}

func writeTestFile(t *testing.T, root string, name string, content string) {
	t.Helper()

	path := filepath.Join(root, filepath.FromSlash(name))
	assert.Nil(t, os.MkdirAll(filepath.Dir(path), 0750))
	assert.Nil(t, os.WriteFile(path, []byte(content), 0600))
}

func isDir(t *testing.T, path string) bool {
	t.Helper()

	info, err := os.Stat(path)
	assert.Nil(t, err)
	return info.IsDir()
}

func fileContains(t *testing.T, path string, needle string) bool {
	t.Helper()

	content, err := os.ReadFile(path)
	assert.Nil(t, err)
	return strings.Contains(string(content), needle)
}
