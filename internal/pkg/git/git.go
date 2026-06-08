package git

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type TreeEntry struct {
	Mode string
	Type string
	OID  string
	Size *int64
	Path string
}

var ErrLsTreeLimitExceeded = errors.New("ls-tree limit exceeded")

type CommitSummary struct {
	OID             string
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

type Ref struct {
	FullName string
	Name     string
	Type     string
	OID      string
}

// Git client, currently just a backend for gitstore and git binary, but eventually will incorporate go-git
type Git struct {
	BinaryPath string
}

func NewGit(binaryPath string) Git {
	return Git{
		BinaryPath: binaryPath,
	}
}

func (g Git) Init(repoPath string) error {
	err := exec.Command(g.BinaryPath, "init", repoPath).Run()
	if err != nil {
		return err
	}

	return nil
}

func (g Git) InitBare(repoPath string, defaultBranch string) error {
	err := exec.Command(g.BinaryPath, "init", "--bare", "--initial-branch="+defaultBranch, repoPath).Run()
	if err != nil {
		return err
	}

	return nil
}

func (g Git) Path() string {
	return g.BinaryPath
}

func (g Git) ResolveCommit(ctx context.Context, repoPath string, ref string) (string, error) {
	cmd := exec.CommandContext(ctx, g.BinaryPath, "-C", repoPath, "rev-parse", "--verify", "--end-of-options", ref+"^{commit}")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func (g Git) ListRefs(ctx context.Context, repoPath string) ([]Ref, error) {
	cmd := exec.CommandContext(ctx, g.BinaryPath, "-C", repoPath, "for-each-ref", "--format=%(refname)%09%(refname:short)%09%(objecttype)%09%(objectname)", "refs/heads", "refs/tags")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	records := strings.Split(strings.TrimSpace(string(output)), "\n")
	refs := make([]Ref, 0, len(records))
	for _, record := range records {
		if record == "" {
			continue
		}

		fullName, rest, ok := strings.Cut(record, "\t")
		if !ok {
			return nil, fmt.Errorf("failed to parse for-each-ref output")
		}
		name, rest, ok := strings.Cut(rest, "\t")
		if !ok {
			return nil, fmt.Errorf("failed to parse for-each-ref output")
		}
		refType, oid, ok := strings.Cut(rest, "\t")
		if !ok {
			return nil, fmt.Errorf("failed to parse for-each-ref output")
		}

		refs = append(refs, Ref{
			FullName: fullName,
			Name:     name,
			Type:     refType,
			OID:      oid,
		})
	}

	return refs, nil
}

func (g Git) LsTreePath(ctx context.Context, repoPath string, commitOID string, filePath string) (*TreeEntry, error) {
	cmd := exec.CommandContext(ctx, g.BinaryPath, "-C", repoPath, "ls-tree", "-l", commitOID, "--", filePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(output)) == "" {
		return nil, nil
	}

	metadata, _, found := strings.Cut(string(output), "\t")
	if !found {
		return nil, fmt.Errorf("failed to parse ls-tree output")
	}

	fields := strings.Fields(metadata)
	if len(fields) != 4 {
		return nil, fmt.Errorf("failed to parse ls-tree metadata")
	}

	size, err := parseLsTreeSize(fields[3])
	if err != nil {
		return nil, err
	}

	return &TreeEntry{
		Mode: fields[0],
		Type: fields[1],
		OID:  fields[2],
		Size: size,
		Path: strings.TrimRight(strings.TrimPrefix(string(output), metadata+"\t"), "\n"),
	}, nil
}

func (g Git) LsTree(ctx context.Context, repoPath string, treeish string) ([]TreeEntry, error) {
	cmd := exec.CommandContext(ctx, g.BinaryPath, "-C", repoPath, "ls-tree", "-l", "-z", treeish)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	records := strings.Split(strings.TrimRight(string(output), "\x00"), "\x00")
	entries := make([]TreeEntry, 0, len(records))
	for _, record := range records {
		if record == "" {
			continue
		}

		entry, err := parseLsTreeRecord(record)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (g Git) LsTreeRecursive(ctx context.Context, repoPath string, treeish string, maxEntries int, maxOutputBytes int64) ([]TreeEntry, error) {
	commandCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := exec.CommandContext(commandCtx, g.BinaryPath, "-C", repoPath, "ls-tree", "-r", "-t", "-l", "-z", treeish)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	reader := bufio.NewReader(io.LimitReader(stdout, maxOutputBytes+1))
	entries := make([]TreeEntry, 0)
	var outputBytes int64

	for {
		record, readErr := reader.ReadString('\x00')
		outputBytes += int64(len(record))
		if outputBytes > maxOutputBytes {
			cancel()
			_ = cmd.Wait()
			return nil, ErrLsTreeLimitExceeded
		}

		record = strings.TrimSuffix(record, "\x00")
		if record != "" {
			entry, parseErr := parseLsTreeRecord(record)
			if parseErr != nil {
				cancel()
				_ = cmd.Wait()
				return nil, parseErr
			}
			entries = append(entries, entry)
			if len(entries) > maxEntries {
				cancel()
				_ = cmd.Wait()
				return nil, ErrLsTreeLimitExceeded
			}
		}

		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			cancel()
			_ = cmd.Wait()
			return nil, readErr
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (g Git) CatFileBlob(ctx context.Context, repoPath string, oid string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, g.BinaryPath, "-C", repoPath, "cat-file", "blob", oid)
	return cmd.Output()
}

func (g Git) LogLatest(ctx context.Context, repoPath string, ref string, filePath string) (*CommitSummary, error) {
	format := "%H%x00%h%x00%s%x00%B%x00%an%x00%ae%x00%aI%x00%cn%x00%ce%x00%cI"
	args := []string{"-C", repoPath, "log", "-1", "--format=" + format, ref, "--"}
	if filePath != "" {
		args = append(args, filePath)
	}

	cmd := exec.CommandContext(ctx, g.BinaryPath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(output)) == "" {
		return nil, nil
	}

	return parseCommitSummary(string(output))
}

func parseLsTreeRecord(record string) (TreeEntry, error) {
	metadata, entryPath, found := strings.Cut(record, "\t")
	if !found {
		return TreeEntry{}, fmt.Errorf("failed to parse ls-tree output")
	}

	fields := strings.Fields(metadata)
	if len(fields) != 4 {
		return TreeEntry{}, fmt.Errorf("failed to parse ls-tree metadata")
	}

	size, err := parseLsTreeSize(fields[3])
	if err != nil {
		return TreeEntry{}, err
	}

	return TreeEntry{
		Mode: fields[0],
		Type: fields[1],
		OID:  fields[2],
		Size: size,
		Path: path.Clean(entryPath),
	}, nil
}

func parseLsTreeSize(raw string) (*int64, error) {
	if raw == "-" {
		return nil, nil
	}

	size, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, err
	}

	return &size, nil
}

func parseCommitSummary(output string) (*CommitSummary, error) {
	fields := strings.Split(strings.TrimRight(output, "\n"), "\x00")
	if len(fields) != 10 {
		return nil, fmt.Errorf("failed to parse git log output")
	}

	authoredAt, err := time.Parse(time.RFC3339, fields[6])
	if err != nil {
		return nil, err
	}

	committedAt, err := time.Parse(time.RFC3339, fields[9])
	if err != nil {
		return nil, err
	}

	return &CommitSummary{
		OID:             fields[0],
		ShortOID:        fields[1],
		MessageHeadline: fields[2],
		Message:         strings.TrimRight(fields[3], "\n"),
		AuthorName:      fields[4],
		AuthorEmail:     fields[5],
		AuthoredAt:      authoredAt,
		CommitterName:   fields[7],
		CommitterEmail:  fields[8],
		CommittedAt:     committedAt,
	}, nil
}
