package gitauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	ScopeRepoRead  = "repo:read"
	ScopeRepoWrite = "repo:write"
	TokenPrefix    = "ctn_pat_"
)

var (
	ErrInvalidCredential = errors.New("invalid git credential")
	ErrInvalidScope      = errors.New("invalid git access token scope")
)

type GitCredentialVerifier interface {
	VerifyGitCredential(ctx context.Context, username string, secret string) (GitPrincipal, error)
}

type GitPrincipal struct {
	User  db.User
	Token db.GitAccessToken
}

func (p GitPrincipal) HasScope(scope string) bool {
	for _, tokenScope := range p.Token.Scopes {
		if tokenScope == scope {
			return true
		}
	}

	return false
}

type TokenIssuer interface {
	CreateToken(context.Context, db.User, string, []string, *time.Time) (string, db.GitAccessToken, error)
}

type Service struct {
	repository db.Queries
}

func NewService(conn db.DBTX) Service {
	return Service{
		repository: *db.New(conn),
	}
}

func (s Service) CreateToken(ctx context.Context, user db.User, name string, scopes []string, expiresAt *time.Time) (string, db.GitAccessToken, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", db.GitAccessToken{}, err
	}

	scopes = NormalizeScopes(scopes)
	if err := ValidateScopes(scopes); err != nil {
		return "", db.GitAccessToken{}, err
	}

	var expires pgtype.Timestamptz
	if expiresAt != nil {
		expires = pgtype.Timestamptz{Time: *expiresAt, Valid: true}
	}

	created, err := s.repository.CreateGitAccessToken(ctx, db.CreateGitAccessTokenParams{
		UserID:      user.ID,
		Name:        name,
		TokenHash:   HashToken(token),
		TokenPrefix: DisplayPrefix(token),
		Scopes:      scopes,
		ExpiresAt:   expires,
	})
	if err != nil {
		return "", db.GitAccessToken{}, err
	}

	return token, created, nil
}

func (s Service) VerifyGitCredential(ctx context.Context, username string, secret string) (GitPrincipal, error) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(secret) == "" {
		return GitPrincipal{}, ErrInvalidCredential
	}

	token, err := s.repository.GetGitAccessTokenByHash(ctx, HashToken(secret))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GitPrincipal{}, ErrInvalidCredential
		}

		return GitPrincipal{}, err
	}

	if token.RevokedAt.Valid {
		return GitPrincipal{}, ErrInvalidCredential
	}

	if token.ExpiresAt.Valid && time.Now().After(token.ExpiresAt.Time) {
		return GitPrincipal{}, ErrInvalidCredential
	}

	user, err := s.repository.GetUserByID(ctx, token.UserID)
	if err != nil {
		return GitPrincipal{}, err
	}

	if user.Name != username {
		return GitPrincipal{}, ErrInvalidCredential
	}

	if err := s.repository.TouchGitAccessTokenLastUsed(ctx, token.ID); err != nil {
		return GitPrincipal{}, err
	}

	return GitPrincipal{User: user, Token: token}, nil
}

func GenerateToken() (string, error) {
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		return "", fmt.Errorf("generate git access token: %w", err)
	}

	return TokenPrefix + base64.RawURLEncoding.EncodeToString(secret[:]), nil
}

func HashToken(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}

func DisplayPrefix(token string) string {
	if len(token) <= len(TokenPrefix)+8 {
		return token
	}

	return token[:len(TokenPrefix)+8]
}

func NormalizeScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return []string{ScopeRepoRead, ScopeRepoWrite}
	}

	seen := map[string]struct{}{}
	normalized := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		if _, ok := seen[scope]; ok {
			continue
		}

		seen[scope] = struct{}{}
		normalized = append(normalized, scope)
	}

	return normalized
}

func ValidateScopes(scopes []string) error {
	if len(scopes) == 0 {
		return ErrInvalidScope
	}

	for _, scope := range scopes {
		switch scope {
		case ScopeRepoRead, ScopeRepoWrite:
		default:
			return fmt.Errorf("%w: %s", ErrInvalidScope, scope)
		}
	}

	return nil
}
