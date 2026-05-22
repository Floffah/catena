package gitauth

import (
	"crypto/sha256"
	"errors"
	"strings"
	"testing"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/zeebo/assert"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken()

	assert.Nil(t, err)
	assert.That(t, strings.HasPrefix(token, TokenPrefix))
	assert.That(t, len(token) > len(TokenPrefix))
}

func TestHashToken(t *testing.T) {
	token := "ctn_pat_test-token"
	want := sha256.Sum256([]byte(token))

	got := HashToken(token)

	assert.That(t, string(got) == string(want[:]))
}

func TestDisplayPrefix(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{
			name:  "long token is shortened",
			token: TokenPrefix + "1234567890abcdef",
			want:  TokenPrefix + "12345678",
		},
		{
			name:  "short token is preserved",
			token: TokenPrefix + "123",
			want:  TokenPrefix + "123",
		},
		{
			name:  "boundary length token is preserved",
			token: TokenPrefix + "12345678",
			want:  TokenPrefix + "12345678",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.That(t, DisplayPrefix(tt.token) == tt.want)
		})
	}
}

func TestNormalizeScopes(t *testing.T) {
	tests := []struct {
		name   string
		scopes []string
		want   []string
	}{
		{
			name: "empty scopes default to read and write",
			want: []string{ScopeRepoRead, ScopeRepoWrite},
		},
		{
			name:   "scopes are trimmed and deduplicated",
			scopes: []string{" repo:read ", "", "repo:write", "repo:read"},
			want:   []string{ScopeRepoRead, ScopeRepoWrite},
		},
		{
			name:   "order is preserved",
			scopes: []string{"repo:write", "repo:read"},
			want:   []string{ScopeRepoWrite, ScopeRepoRead},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeScopes(tt.scopes)

			assert.That(t, len(got) == len(tt.want))
			for i := range tt.want {
				assert.That(t, got[i] == tt.want[i])
			}
		})
	}
}

func TestValidateScopes(t *testing.T) {
	tests := []struct {
		name    string
		scopes  []string
		wantErr bool
	}{
		{
			name:    "empty scopes are invalid",
			wantErr: true,
		},
		{
			name:   "known scopes are valid",
			scopes: []string{ScopeRepoRead, ScopeRepoWrite},
		},
		{
			name:    "unknown scope is invalid",
			scopes:  []string{ScopeRepoRead, "admin:everything"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateScopes(tt.scopes)
			if tt.wantErr {
				assert.That(t, errors.Is(err, ErrInvalidScope))
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestGitPrincipalHasScope(t *testing.T) {
	principal := GitPrincipal{
		Token: db.GitAccessToken{
			Scopes: []string{ScopeRepoRead},
		},
	}

	assert.That(t, principal.HasScope(ScopeRepoRead))
	assert.That(t, !principal.HasScope(ScopeRepoWrite))
}
