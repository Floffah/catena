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

func TestListRepositoryIssues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 22, 13, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13c01")
	repositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13c02")
	issueID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13c03")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)
	repository := testRepository(repositoryID, userID, "catena", nil, db.RepositoryVisibilityPublic, "main", createdAt, updatedAt)

	mock, err := pgxmock.NewPool()
	assert.Nil(t, err)
	defer mock.Close()

	expectRepositoryByOwnerAndName(mock, user.Name, repository.Name, repository)
	mock.ExpectQuery("select (.+) from repository_items").
		WithArgs(repository.ID).
		WillReturnRows(issueRows().AddRow(
			UUIDToPgtype(issueID),
			repository.ID,
			int64(1),
			db.RepositoryItemKindIssue,
			"First issue",
			ptr("Issue body"),
			user.ID,
			pgtype.Timestamptz{Time: createdAt, Valid: true},
			pgtype.Timestamptz{Time: updatedAt, Valid: true},
			pgtype.Timestamptz{Time: updatedAt, Valid: true},
			db.IssueStatusOpen,
		))

	router := NewRouter(ServerDeps{
		DB:   mock,
		Auth: testAuthProvider{user: user},
	})
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/repositories/floffah/catena/issues", nil)

	router.ServeHTTP(response, request)

	assert.That(t, response.Code == http.StatusOK)
	assert.Nil(t, mock.ExpectationsWereMet())

	var body ListIssuesResponse
	assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
	assert.That(t, len(body.Issues) == 1)
	assert.That(t, body.Issues[0].Id == issueID)
	assert.That(t, body.Issues[0].Reference == "I-1")
	assert.That(t, body.Issues[0].Status == IssueStatusOpen)
	assert.That(t, body.Issues[0].Title == "First issue")
}

func TestCreateRepositoryIssue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 22, 14, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13c04")
	repositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13c05")
	issueID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13c06")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)
	repository := testRepository(repositoryID, userID, "catena", nil, db.RepositoryVisibilityPublic, "main", createdAt, updatedAt)

	t.Run("anonymous request returns unauthorized", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectRepositoryByOwnerAndName(mock, user.Name, repository.Name, repository)

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/v1/repositories/floffah/catena/issues", bytes.NewBufferString(`{"title":"Issue"}`))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusUnauthorized)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("blank title returns bad request", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectRepositoryByOwnerAndName(mock, user.Name, repository.Name, repository)

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/v1/repositories/floffah/catena/issues", bytes.NewBufferString(`{"title":"   "}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusBadRequest)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body BadRequestJSONResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Error == "issue title is required")
	})

	t.Run("authenticated request creates issue", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		body := "Issue body"
		expectRepositoryByOwnerAndName(mock, user.Name, repository.Name, repository)
		mock.ExpectBegin()
		mock.ExpectQuery("update repositories").
			WithArgs(repository.ID).
			WillReturnRows(pgxmock.NewRows([]string{"number"}).AddRow(int64(1)))
		mock.ExpectQuery("insert into repository_items").
			WithArgs(repository.ID, int64(1), db.RepositoryItemKindIssue, "First issue", pgxmock.AnyArg(), user.ID).
			WillReturnRows(repositoryItemRows().AddRow(
				UUIDToPgtype(issueID),
				repository.ID,
				int64(1),
				db.RepositoryItemKindIssue,
				"First issue",
				&body,
				user.ID,
				pgtype.Timestamptz{Time: createdAt, Valid: true},
				pgtype.Timestamptz{Time: updatedAt, Valid: true},
				pgtype.Timestamptz{Time: updatedAt, Valid: true},
			))
		mock.ExpectQuery("insert into issues").
			WithArgs(UUIDToPgtype(issueID)).
			WillReturnRows(pgxmock.NewRows([]string{"repository_item_id", "kind", "status"}).
				AddRow(UUIDToPgtype(issueID), db.RepositoryItemKindIssue, db.IssueStatusOpen))
		mock.ExpectCommit()

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/v1/repositories/floffah/catena/issues", bytes.NewBufferString(`{"title":" First issue ","body":"Issue body"}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusCreated)
		assert.Nil(t, mock.ExpectationsWereMet())

		var responseBody Issue
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &responseBody))
		assert.That(t, responseBody.Id == issueID)
		assert.That(t, responseBody.Reference == "I-1")
		assert.That(t, responseBody.Title == "First issue")
		assert.That(t, responseBody.Body != nil && *responseBody.Body == body)
		assert.That(t, responseBody.Status == IssueStatusOpen)
	})
}

func TestGetRepositoryIssue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 22, 15, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13c07")
	repositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13c08")
	issueID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13c09")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)
	repository := testRepository(repositoryID, userID, "catena", nil, db.RepositoryVisibilityPublic, "main", createdAt, updatedAt)

	t.Run("invalid issue number returns bad request", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/repositories/floffah/catena/issues/0", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusBadRequest)
	})

	t.Run("existing issue returns issue", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectRepositoryByOwnerAndName(mock, user.Name, repository.Name, repository)
		mock.ExpectQuery("select (.+) from repository_items").
			WithArgs(repository.ID, int64(1)).
			WillReturnRows(issueRows().AddRow(
				UUIDToPgtype(issueID),
				repository.ID,
				int64(1),
				db.RepositoryItemKindIssue,
				"First issue",
				ptr("Issue body"),
				user.ID,
				pgtype.Timestamptz{Time: createdAt, Valid: true},
				pgtype.Timestamptz{Time: updatedAt, Valid: true},
				pgtype.Timestamptz{Time: updatedAt, Valid: true},
				db.IssueStatusInProgress,
			))

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/repositories/floffah/catena/issues/1", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body Issue
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Id == issueID)
		assert.That(t, body.Reference == "I-1")
		assert.That(t, body.Status == IssueStatusInProgress)
	})

	t.Run("missing issue returns not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectRepositoryByOwnerAndName(mock, user.Name, repository.Name, repository)
		mock.ExpectQuery("select (.+) from repository_items").
			WithArgs(repository.ID, int64(404)).
			WillReturnError(pgx.ErrNoRows)

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/repositories/floffah/catena/issues/404", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusNotFound)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func issueRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"repository_id",
		"number",
		"kind",
		"title",
		"body",
		"author_id",
		"created_at",
		"updated_at",
		"last_activity_at",
		"status",
	})
}

func repositoryItemRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"repository_id",
		"number",
		"kind",
		"title",
		"body",
		"author_id",
		"created_at",
		"updated_at",
		"last_activity_at",
	})
}
