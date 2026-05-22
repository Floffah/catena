package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/zeebo/assert"
)

func TestGetAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 21, 14, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13afa")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)

	router := NewRouter(ServerDeps{
		DB:   failDB{t: t},
		Auth: testAuthProvider{user: user},
	})

	t.Run("anonymous request returns unauthorized", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/user", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusUnauthorized)

		var body UnauthorizedJSONResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Error == "unauthorized")
	})

	t.Run("authenticated request returns current user", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
		request.Header.Set("Authorization", "Bearer "+testBearerToken)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)

		var body User
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Id == userID)
		assert.That(t, body.Name == "floffah")
		assert.That(t, body.Email != nil && string(*body.Email) == "floffah@example.com")
		assert.That(t, body.DisplayName != nil && *body.DisplayName == "Floffah")
		assert.That(t, body.AvatarUrl != nil && *body.AvatarUrl == "https://example.com/avatar.png")
		assert.That(t, body.CreatedAt.Equal(createdAt))
		assert.That(t, body.UpdatedAt.Equal(updatedAt))
	})
}

func TestUpdateAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 21, 15, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13afb")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)

	t.Run("anonymous request returns unauthorized", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPatch, "/v1/user", bytes.NewBufferString(`{"displayName":"New Name"}`))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusUnauthorized)

		var body UnauthorizedJSONResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Error == "unauthorized")
	})

	t.Run("missing displayName returns current user without database update", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPatch, "/v1/user", bytes.NewBufferString(`{}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)

		var body User
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Id == userID)
		assert.That(t, body.DisplayName != nil && *body.DisplayName == "Floffah")
	})

	t.Run("blank displayName returns bad request", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPatch, "/v1/user", bytes.NewBufferString(`{"displayName":"   "}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusBadRequest)

		var body BadRequestJSONResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Error == "displayName must not be empty")
	})

	t.Run("displayName is trimmed and updated", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		newDisplayName := "New Floffah"
		updatedUser := testUser(userID, "floffah", newDisplayName, "https://example.com/avatar.png", createdAt, updatedAt)
		mock.ExpectQuery("update users").
			WithArgs(user.ID, user.Name, pgxmock.AnyArg(), user.AvatarUrl).
			WillReturnRows(userRows().AddRow(
				updatedUser.ID,
				updatedUser.ClerkUserID,
				updatedUser.Name,
				updatedUser.DisplayName,
				updatedUser.AvatarUrl,
				updatedUser.Email,
				updatedUser.CreatedAt,
				updatedUser.UpdatedAt,
			))

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPatch, "/v1/user", bytes.NewBufferString(`{"displayName":"  New Floffah  "}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body User
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Id == userID)
		assert.That(t, body.DisplayName != nil && *body.DisplayName == newDisplayName)
	})
}

func TestGetUserByClerkUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 21, 16, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13afc")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)

	router := NewRouter(ServerDeps{
		DB:   failDB{t: t},
		Auth: testAuthProvider{user: user},
	})

	t.Run("anonymous request returns unauthorized", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/clerk/user_123", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusUnauthorized)

		var body UnauthorizedJSONResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Error == "unauthorized")
	})

	t.Run("different clerk user id returns forbidden", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/clerk/user_other", nil)
		request.Header.Set("Authorization", "Bearer "+testBearerToken)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusForbidden)

		var body ForbiddenJSONResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Error == "forbidden")
	})

	t.Run("matching clerk user id returns current user", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/clerk/user_123", nil)
		request.Header.Set("Authorization", "Bearer "+testBearerToken)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)

		var body User
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Id == userID)
		assert.That(t, body.Name == "floffah")
		assert.That(t, body.DisplayName != nil && *body.DisplayName == "Floffah")
	})
}

func TestGetUserByName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 21, 17, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13afd")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)

	t.Run("existing user returns profile", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		mock.ExpectQuery("select (.+) from users").
			WithArgs(user.Name).
			WillReturnRows(userRows().AddRow(
				user.ID,
				user.ClerkUserID,
				user.Name,
				user.DisplayName,
				user.AvatarUrl,
				user.Email,
				user.CreatedAt,
				user.UpdatedAt,
			))

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/name/floffah", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body User
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Id == userID)
		assert.That(t, body.Name == "floffah")
		assert.That(t, body.DisplayName != nil && *body.DisplayName == "Floffah")
	})

	t.Run("missing user returns not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		mock.ExpectQuery("select (.+) from users").
			WithArgs("missing").
			WillReturnError(pgx.ErrNoRows)

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/name/missing", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusNotFound)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body NotFoundJSONResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Error == "user not found")
	})
}

func testUser(id uuid.UUID, name string, displayName string, avatarURL string, createdAt time.Time, updatedAt time.Time) db.User {
	email := name + "@example.com"

	return db.User{
		ID:          UUIDToPgtype(id),
		ClerkUserID: "user_123",
		Name:        name,
		DisplayName: &displayName,
		AvatarUrl:   &avatarURL,
		Email:       email,
		CreatedAt:   pgtype.Timestamptz{Time: createdAt, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: updatedAt, Valid: true},
	}
}

func userRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"clerk_user_id",
		"name",
		"display_name",
		"avatar_url",
		"email",
		"created_at",
		"updated_at",
	})
}
