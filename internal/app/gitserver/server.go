package gitserver

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/gitauth"
	"github.com/floffah/catena/internal/pkg/gitstore"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	repository db.Queries
	git        gitstore.Store
	verifier   gitauth.GitCredentialVerifier
}

func NewHandler(conn db.DBTX, gitService gitstore.Store, verifier gitauth.GitCredentialVerifier) Handler {
	return Handler{
		repository: *db.New(conn),
		git:        gitService,
		verifier:   verifier,
	}
}

func (h Handler) Handle(c *gin.Context) {
	request, ok := parseRequestPath(c.Request.URL.Path)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	operation, ok := getOperation(request.GitPath, c.Request.URL.Query().Get("service"))
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	repository, err := h.repository.GetRepositoryByOwnerAndName(c.Request.Context(), db.GetRepositoryByOwnerAndNameParams{
		OwnerName:      request.Owner,
		RepositoryName: request.Repository,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.Status(http.StatusNotFound)
			return
		}

		c.String(http.StatusInternalServerError, "failed to load repository")
		return
	}

	remoteUser, ok := h.authorize(c, repository, operation)
	if !ok {
		return
	}

	err = h.serveGitHTTP(c, repository, request.GitPath, remoteUser)
	if err != nil {
		c.Error(err)
		if !c.Writer.Written() {
			c.String(http.StatusInternalServerError, "failed to serve git repository")
		}
	}
}

func (h Handler) authorize(c *gin.Context, repository db.Repository, operation gitOperation) (string, bool) {
	if operation == gitOperationUploadPack && repository.Visibility == db.RepositoryVisibilityPublic {
		return "", true
	}

	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Header("WWW-Authenticate", `Basic realm="Catena Git"`)
		c.Status(http.StatusUnauthorized)
		return "", false
	}

	principal, err := h.verifier.VerifyGitCredential(c.Request.Context(), username, password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, gitauth.ErrInvalidCredential) {
			status = http.StatusUnauthorized
			c.Header("WWW-Authenticate", `Basic realm="Catena Git"`)
		}

		c.Status(status)
		return "", false
	}

	if principal.User.ID != repository.OwnerID {
		c.Status(http.StatusForbidden)
		return "", false
	}

	requiredScope := gitauth.ScopeRepoRead
	if operation == gitOperationReceivePack {
		requiredScope = gitauth.ScopeRepoWrite
	}
	if !principal.HasScope(requiredScope) {
		c.Status(http.StatusForbidden)
		return "", false
	}

	return principal.User.Name, true
}

func (h Handler) serveGitHTTP(c *gin.Context, repository db.Repository, gitPath string, remoteUser string) error {
	repoPath := h.git.GetRepoPath(repository)
	projectRoot := filepath.Dir(repoPath)
	pathInfo := "/" + filepath.Base(repoPath) + gitPath

	cmd := exec.CommandContext(c.Request.Context(), h.git.GitBinaryPath(), "http-backend")
	cmd.Env = append(os.Environ(), buildCGIEnv(c, projectRoot, pathInfo, remoteUser)...)
	cmd.Stdin = c.Request.Body

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("create git-http-backend stdout pipe: %w", err)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start git-http-backend: %w", err)
	}

	reader := bufio.NewReader(stdout)
	statusCode, err := writeCGIHeaders(c, reader)
	if err != nil {
		cmd.Wait()
		return err
	}

	c.Status(statusCode)
	if _, err := io.Copy(c.Writer, reader); err != nil {
		cmd.Wait()
		return fmt.Errorf("stream git-http-backend response: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		stderrText := strings.TrimSpace(stderr.String())
		if stderrText != "" {
			return fmt.Errorf("git-http-backend failed: %w: %s", err, stderrText)
		}

		return fmt.Errorf("git-http-backend failed: %w", err)
	}

	return nil
}

type requestPath struct {
	Owner      string
	Repository string
	GitPath    string
}

type gitOperation string

const (
	gitOperationUploadPack  gitOperation = "git-upload-pack"
	gitOperationReceivePack gitOperation = "git-receive-pack"
)

func parseRequestPath(rawPath string) (requestPath, bool) {
	parts := strings.Split(strings.Trim(rawPath, "/"), "/")
	if len(parts) < 3 {
		return requestPath{}, false
	}

	repository := strings.TrimSuffix(parts[1], ".git")
	gitPath := "/" + strings.Join(parts[2:], "/")
	if gitPath != "/info/refs" && gitPath != "/git-upload-pack" && gitPath != "/git-receive-pack" {
		return requestPath{}, false
	}

	return requestPath{
		Owner:      parts[0],
		Repository: repository,
		GitPath:    gitPath,
	}, true
}

func getOperation(gitPath string, service string) (gitOperation, bool) {
	if gitPath == "/git-upload-pack" {
		return gitOperationUploadPack, true
	}
	if gitPath == "/git-receive-pack" {
		return gitOperationReceivePack, true
	}
	if gitPath == "/info/refs" {
		switch service {
		case "git-upload-pack":
			return gitOperationUploadPack, true
		case "git-receive-pack":
			return gitOperationReceivePack, true
		}
	}

	return "", false
}

func buildCGIEnv(c *gin.Context, projectRoot string, pathInfo string, remoteUser string) []string {
	contentLength := c.Request.Header.Get("Content-Length")
	if contentLength == "" && c.Request.ContentLength >= 0 {
		contentLength = strconv.FormatInt(c.Request.ContentLength, 10)
	}

	env := []string{
		"GIT_HTTP_EXPORT_ALL=1",
		"GIT_PROJECT_ROOT=" + projectRoot,
		"PATH_INFO=" + pathInfo,
		"REQUEST_METHOD=" + c.Request.Method,
		"QUERY_STRING=" + c.Request.URL.RawQuery,
		"CONTENT_TYPE=" + c.Request.Header.Get("Content-Type"),
		"CONTENT_LENGTH=" + contentLength,
		"REMOTE_ADDR=" + remoteAddr(c.Request),
		"SERVER_PROTOCOL=" + c.Request.Proto,
		"HTTP_HOST=" + c.Request.Host,
	}

	if gitProtocol := c.Request.Header.Get("Git-Protocol"); gitProtocol != "" {
		env = append(env, "HTTP_GIT_PROTOCOL="+gitProtocol, "GIT_PROTOCOL="+gitProtocol)
	}

	if userAgent := c.Request.Header.Get("User-Agent"); userAgent != "" {
		env = append(env, "HTTP_USER_AGENT="+userAgent)
	}
	if remoteUser != "" {
		env = append(env, "REMOTE_USER="+remoteUser)
	}

	return env
}

func writeCGIHeaders(c *gin.Context, reader *bufio.Reader) (int, error) {
	statusCode := http.StatusOK

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return 0, fmt.Errorf("read git-http-backend headers: %w", err)
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			return statusCode, nil
		}

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return 0, fmt.Errorf("invalid git-http-backend header: %q", line)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if strings.EqualFold(key, "Status") {
			code, err := parseStatusCode(value)
			if err != nil {
				return 0, err
			}

			statusCode = code
			continue
		}

		c.Header(key, value)
	}
}

func parseStatusCode(status string) (int, error) {
	fields := strings.Fields(status)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty git-http-backend status")
	}

	code, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid git-http-backend status %q: %w", status, err)
	}

	return code, nil
}

func remoteAddr(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
