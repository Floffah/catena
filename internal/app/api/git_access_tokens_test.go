package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/gitauth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/zeebo/assert"
)

func TestListGitAccessTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 22, 16, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13d01")
	tokenID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13d02")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)

	t.Run("anonymous request returns unauthorized", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/git-access-tokens", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusUnauthorized)
	})

	t.Run("authenticated request lists active tokens", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expiresAt := createdAt.Add(24 * time.Hour)
		token := testGitAccessToken(tokenID, user.ID, "Local laptop", "ctn_pat_12345678", []string{gitauth.ScopeRepoRead}, createdAt, updatedAt)
		token.ExpiresAt = pgtype.Timestamptz{Time: expiresAt, Valid: true}
		mock.ExpectQuery("select (.+) from git_access_tokens").
			WithArgs(user.ID).
			WillReturnRows(gitAccessTokenRows().AddRow(
				token.ID,
				token.UserID,
				token.Name,
				token.TokenHash,
				token.TokenPrefix,
				token.Scopes,
				token.LastUsedAt,
				token.ExpiresAt,
				token.RevokedAt,
				token.CreatedAt,
				token.UpdatedAt,
			))

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/git-access-tokens", nil)
		request.Header.Set("Authorization", "Bearer "+testBearerToken)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body []GitAccessToken
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, len(body) == 1)
		assert.That(t, body[0].Id == tokenID)
		assert.That(t, body[0].Name == token.Name)
		assert.That(t, body[0].TokenPrefix == token.TokenPrefix)
		assert.That(t, len(body[0].Scopes) == 1)
		assert.That(t, body[0].Scopes[0] == gitauth.ScopeRepoRead)
		assert.That(t, body[0].ExpiresAt != nil && body[0].ExpiresAt.Equal(expiresAt))
	})
}

func TestCreateGitAccessToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 22, 17, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13d03")
	tokenID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13d04")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)
	token := testGitAccessToken(tokenID, user.ID, "Local laptop", "ctn_pat_12345678", []string{gitauth.ScopeRepoRead, gitauth.ScopeRepoWrite}, createdAt, updatedAt)

	t.Run("anonymous request returns unauthorized", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:        failDB{t: t},
			Auth:      testAuthProvider{user: user},
			GitTokens: testGitTokenIssuer{},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/v1/git-access-tokens", bytes.NewBufferString(`{"name":"Local laptop"}`))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusUnauthorized)
	})

	t.Run("blank name returns bad request", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:        failDB{t: t},
			Auth:      testAuthProvider{user: user},
			GitTokens: testGitTokenIssuer{},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/v1/git-access-tokens", bytes.NewBufferString(`{"name":"   "}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusBadRequest)
	})

	t.Run("invalid scope returns bad request", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:        failDB{t: t},
			Auth:      testAuthProvider{user: user},
			GitTokens: testGitTokenIssuer{},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/v1/git-access-tokens", bytes.NewBufferString(`{"name":"Local laptop","scopes":["repo:admin"]}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusBadRequest)

		var body BadRequestJSONResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Error == "token scopes are invalid")
	})

	t.Run("past expiry returns bad request", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:        failDB{t: t},
			Auth:      testAuthProvider{user: user},
			GitTokens: testGitTokenIssuer{},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/v1/git-access-tokens", bytes.NewBufferString(`{"name":"Local laptop","expiresAt":"2000-01-01T00:00:00Z"}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusBadRequest)
	})

	t.Run("valid request creates token", func(t *testing.T) {
		issuer := testGitTokenIssuer{
			rawToken: "ctn_pat_rawsecret",
			token:    token,
		}
		router := NewRouter(ServerDeps{
			DB:        failDB{t: t},
			Auth:      testAuthProvider{user: user},
			GitTokens: issuer,
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/v1/git-access-tokens", bytes.NewBufferString(`{"name":" Local laptop ","scopes":["repo:read","repo:write"]}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusCreated)

		var body CreateGitAccessTokenResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Token == issuer.rawToken)
		assert.That(t, body.AccessToken.Id == tokenID)
		assert.That(t, body.AccessToken.Name == token.Name)
		assert.That(t, body.AccessToken.TokenPrefix == token.TokenPrefix)
		assert.That(t, len(body.AccessToken.Scopes) == 2)
	})
}

func TestRevokeGitAccessToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 22, 18, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13d05")
	tokenID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13d06")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)

	t.Run("anonymous request returns unauthorized", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/v1/git-access-tokens/"+tokenID.String(), nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusUnauthorized)
	})

	t.Run("authenticated request revokes token", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		mock.ExpectExec("update git_access_tokens").
			WithArgs(UUIDToPgtype(tokenID), user.ID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/v1/git-access-tokens/"+tokenID.String(), nil)
		request.Header.Set("Authorization", "Bearer "+testBearerToken)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusNoContent)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.That(t, response.Body.Len() == 0)
	})
}

type testGitTokenIssuer struct {
	rawToken string
	token    db.GitAccessToken
	err      error
}

func (i testGitTokenIssuer) CreateToken(_ context.Context, _ db.User, _ string, _ []string, _ *time.Time) (string, db.GitAccessToken, error) {
	return i.rawToken, i.token, i.err
}

func testGitAccessToken(id uuid.UUID, userID pgtype.UUID, name string, prefix string, scopes []string, createdAt time.Time, updatedAt time.Time) db.GitAccessToken {
	return db.GitAccessToken{
		ID:          UUIDToPgtype(id),
		UserID:      userID,
		Name:        name,
		TokenHash:   []byte("hashed-token"),
		TokenPrefix: prefix,
		Scopes:      scopes,
		CreatedAt:   pgtype.Timestamptz{Time: createdAt, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: updatedAt, Valid: true},
	}
}

func gitAccessTokenRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"user_id",
		"name",
		"token_hash",
		"token_prefix",
		"scopes",
		"last_used_at",
		"expires_at",
		"revoked_at",
		"created_at",
		"updated_at",
	})
}
