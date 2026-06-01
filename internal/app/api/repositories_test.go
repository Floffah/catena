package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/gitstore"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/zeebo/assert"
)

func TestCreateRepository(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gitBin := requireGit(t)
	createdAt := time.Date(2026, 5, 22, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13b01")
	repositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13b02")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)
	store := gitstore.NewStore(t.TempDir(), gitBin)
	repository := testRepository(repositoryID, userID, "catena", nil, db.RepositoryVisibilityPublic, "main", createdAt, updatedAt)

	t.Run("anonymous request returns unauthorized", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
			Git:  store,
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/v1/repositories", bytes.NewBufferString(`{"name":"catena"}`))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusUnauthorized)
	})

	t.Run("authenticated request creates repository and bare git repo", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		mock.ExpectBegin()
		mock.ExpectQuery("insert into repositories").
			WithArgs(user.ID, repository.Name, pgxmock.AnyArg(), db.RepositoryVisibilityPublic, repository.DefaultBranch).
			WillReturnRows(repositoryRows().AddRow(
				repository.ID,
				repository.OwnerID,
				repository.Name,
				repository.Description,
				repository.Visibility,
				repository.DefaultBranch,
				repository.CreatedAt,
				repository.UpdatedAt,
				repository.ItemPrefix,
				repository.NextItemNumber,
			))
		mock.ExpectCommit()

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
			Git:  store,
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/v1/repositories", bytes.NewBufferString(`{"name":" catena ","visibility":"public"}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusCreated)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.That(t, isDir(t, store.GetRepoPath(repository)))

		var body Repository
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Id == repositoryID)
		assert.That(t, body.OwnerId == userID)
		assert.That(t, body.OwnerName == user.Name)
		assert.That(t, body.Name == repository.Name)
		assert.That(t, body.DefaultBranch == "main")
		assert.That(t, body.Visibility == Public)
	})
}

func TestGetRepositoryByOwnerAndName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 22, 11, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13b03")
	repositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13b04")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)
	publicRepository := testRepository(repositoryID, userID, "catena", nil, db.RepositoryVisibilityPublic, "main", createdAt, updatedAt)

	t.Run("public repository can be read anonymously", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectRepositoryByOwnerAndName(mock, user.Name, publicRepository.Name, publicRepository)

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/repositories/floffah/catena", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body Repository
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Id == repositoryID)
		assert.That(t, body.OwnerName == user.Name)
		assert.That(t, body.Visibility == Public)
	})

	t.Run("private repository requires authentication", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		privateRepository := publicRepository
		privateRepository.Visibility = db.RepositoryVisibilityPrivate
		expectRepositoryByOwnerAndName(mock, user.Name, privateRepository.Name, privateRepository)

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/repositories/floffah/catena", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusUnauthorized)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("private repository owner can read", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		privateRepository := publicRepository
		privateRepository.Visibility = db.RepositoryVisibilityPrivate
		expectRepositoryByOwnerAndName(mock, user.Name, privateRepository.Name, privateRepository)

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/repositories/floffah/catena", nil)
		request.Header.Set("Authorization", "Bearer "+testBearerToken)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body Repository
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Visibility == Private)
	})
}

func TestListUserRepositoriesByName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 22, 11, 30, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13ba1")
	otherUserID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13ba4")
	publicRepositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13ba2")
	privateRepositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13ba3")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)
	otherUser := testUser(otherUserID, "octo", "Octo", "https://example.com/octo.png", createdAt, updatedAt)
	publicRepository := testRepository(publicRepositoryID, userID, "catena", nil, db.RepositoryVisibilityPublic, "main", createdAt, updatedAt)
	privateRepository := testRepository(privateRepositoryID, userID, "secret", nil, db.RepositoryVisibilityPrivate, "main", createdAt, updatedAt.Add(-time.Hour))

	t.Run("anonymous request defaults to updated sort and max limit", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectUserByName(mock, user)
		mock.ExpectQuery("select (.+) from repositories").
			WithArgs(user.ID, false, false, db.RepositoryVisibilityPublic, int32(maxListLimit)).
			WillReturnRows(repositoryRows().AddRow(
				publicRepository.ID,
				publicRepository.OwnerID,
				publicRepository.Name,
				publicRepository.Description,
				publicRepository.Visibility,
				publicRepository.DefaultBranch,
				publicRepository.CreatedAt,
				publicRepository.UpdatedAt,
				publicRepository.ItemPrefix,
				publicRepository.NextItemNumber,
			))

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/name/floffah/repositories", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body ListRepositoriesResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, len(body.Repositories) == 1)
		assert.That(t, body.Repositories[0].Id == publicRepositoryID)
		assert.That(t, body.Repositories[0].OwnerName == user.Name)
		assert.That(t, body.Repositories[0].Visibility == Public)
	})

	t.Run("owner request can include private repositories with featured sort", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectUserByName(mock, user)
		mock.ExpectQuery("select (.+) from repositories").
			WithArgs(user.ID, false, true, db.RepositoryVisibilityPublic, int32(6)).
			WillReturnRows(repositoryRows().
				AddRow(
					publicRepository.ID,
					publicRepository.OwnerID,
					publicRepository.Name,
					publicRepository.Description,
					publicRepository.Visibility,
					publicRepository.DefaultBranch,
					publicRepository.CreatedAt,
					publicRepository.UpdatedAt,
					publicRepository.ItemPrefix,
					publicRepository.NextItemNumber,
				).
				AddRow(
					privateRepository.ID,
					privateRepository.OwnerID,
					privateRepository.Name,
					privateRepository.Description,
					privateRepository.Visibility,
					privateRepository.DefaultBranch,
					privateRepository.CreatedAt,
					privateRepository.UpdatedAt,
					privateRepository.ItemPrefix,
					privateRepository.NextItemNumber,
				))

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/name/floffah/repositories?sort=featured&limit=6", nil)
		request.Header.Set("Authorization", "Bearer "+testBearerToken)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body ListRepositoriesResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, len(body.Repositories) == 2)
		assert.That(t, body.Repositories[0].Id == publicRepositoryID)
		assert.That(t, body.Repositories[1].Id == privateRepositoryID)
		assert.That(t, body.Repositories[1].Visibility == Private)
	})

	t.Run("authenticated non-owner request cannot include private repositories", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectUserByName(mock, user)
		mock.ExpectQuery("select (.+) from repositories").
			WithArgs(user.ID, false, false, db.RepositoryVisibilityPublic, int32(maxListLimit)).
			WillReturnRows(repositoryRows().AddRow(
				publicRepository.ID,
				publicRepository.OwnerID,
				publicRepository.Name,
				publicRepository.Description,
				publicRepository.Visibility,
				publicRepository.DefaultBranch,
				publicRepository.CreatedAt,
				publicRepository.UpdatedAt,
				publicRepository.ItemPrefix,
				publicRepository.NextItemNumber,
			))

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: otherUser},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/name/floffah/repositories", nil)
		request.Header.Set("Authorization", "Bearer "+testBearerToken)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body ListRepositoriesResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, len(body.Repositories) == 1)
		assert.That(t, body.Repositories[0].Id == publicRepositoryID)
		assert.That(t, body.Repositories[0].Visibility == Public)
	})

	t.Run("owner request can explicitly filter public repositories", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectUserByName(mock, user)
		mock.ExpectQuery("select (.+) from repositories").
			WithArgs(user.ID, true, true, db.RepositoryVisibilityPublic, int32(6)).
			WillReturnRows(repositoryRows().AddRow(
				publicRepository.ID,
				publicRepository.OwnerID,
				publicRepository.Name,
				publicRepository.Description,
				publicRepository.Visibility,
				publicRepository.DefaultBranch,
				publicRepository.CreatedAt,
				publicRepository.UpdatedAt,
				publicRepository.ItemPrefix,
				publicRepository.NextItemNumber,
			))

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/name/floffah/repositories?sort=featured&limit=6&visibility=public", nil)
		request.Header.Set("Authorization", "Bearer "+testBearerToken)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body ListRepositoriesResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, len(body.Repositories) == 1)
		assert.That(t, body.Repositories[0].Id == publicRepositoryID)
		assert.That(t, body.Repositories[0].Visibility == Public)
	})

	t.Run("anonymous request cannot see explicit private repositories", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectUserByName(mock, user)
		mock.ExpectQuery("select (.+) from repositories").
			WithArgs(user.ID, true, false, db.RepositoryVisibilityPrivate, int32(maxListLimit)).
			WillReturnRows(repositoryRows())

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/name/floffah/repositories?visibility=private", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body ListRepositoriesResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, len(body.Repositories) == 0)
	})

	t.Run("invalid limit returns bad request", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/name/floffah/repositories?limit=0", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusBadRequest)
	})

	t.Run("invalid sort returns bad request", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/name/floffah/repositories?sort=random", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusBadRequest)
	})

	t.Run("invalid visibility returns bad request", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/name/floffah/repositories?visibility=secret", nil)

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusBadRequest)
	})
}

func TestUpdateRepository(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Date(2026, 5, 22, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13b05")
	otherUserID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13b06")
	repositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13b07")
	user := testUser(userID, "floffah", "Floffah", "https://example.com/avatar.png", createdAt, updatedAt)
	repository := testRepository(repositoryID, userID, "catena", ptr("A Git platform"), db.RepositoryVisibilityPublic, "main", createdAt, updatedAt)

	t.Run("anonymous request returns unauthorized", func(t *testing.T) {
		router := NewRouter(ServerDeps{
			DB:   failDB{t: t},
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPatch, "/v1/repositories/floffah/catena", bytes.NewBufferString(`{"description":"New"}`))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusUnauthorized)
	})

	t.Run("non-owner cannot update public repository", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		otherRepository := repository
		otherRepository.OwnerID = UUIDToPgtype(otherUserID)
		expectRepositoryByOwnerAndName(mock, user.Name, repository.Name, otherRepository)

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPatch, "/v1/repositories/floffah/catena", bytes.NewBufferString(`{"description":"New"}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusForbidden)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("owner can update metadata", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		updatedRepository := repository
		updatedRepository.Description = nil
		updatedRepository.Visibility = db.RepositoryVisibilityPrivate
		updatedRepository.UpdatedAt = pgtype.Timestamptz{Time: updatedAt.Add(time.Hour), Valid: true}

		expectRepositoryByOwnerAndName(mock, user.Name, repository.Name, repository)
		mock.ExpectQuery("update repositories").
			WithArgs(repository.ID, repository.Name, pgxmock.AnyArg(), db.RepositoryVisibilityPrivate, repository.DefaultBranch).
			WillReturnRows(repositoryRows().AddRow(
				updatedRepository.ID,
				updatedRepository.OwnerID,
				updatedRepository.Name,
				updatedRepository.Description,
				updatedRepository.Visibility,
				updatedRepository.DefaultBranch,
				updatedRepository.CreatedAt,
				updatedRepository.UpdatedAt,
				updatedRepository.ItemPrefix,
				updatedRepository.NextItemNumber,
			))

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPatch, "/v1/repositories/floffah/catena", bytes.NewBufferString(`{"description":"","visibility":"private"}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusOK)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body Repository
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Description == nil)
		assert.That(t, body.Visibility == Private)
	})

	t.Run("default branch must exist when changed", func(t *testing.T) {
		gitBin := requireGit(t)
		store := gitstore.NewStore(t.TempDir(), gitBin)
		assert.Nil(t, store.CreateRepo(repository))

		mock, err := pgxmock.NewPool()
		assert.Nil(t, err)
		defer mock.Close()

		expectRepositoryByOwnerAndName(mock, user.Name, repository.Name, repository)

		router := NewRouter(ServerDeps{
			DB:   mock,
			Auth: testAuthProvider{user: user},
			Git:  store,
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPatch, "/v1/repositories/floffah/catena", bytes.NewBufferString(`{"defaultBranch":"develop"}`))
		request.Header.Set("Authorization", "Bearer "+testBearerToken)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(response, request)

		assert.That(t, response.Code == http.StatusBadRequest)
		assert.Nil(t, mock.ExpectationsWereMet())

		var body BadRequestJSONResponse
		assert.Nil(t, json.Unmarshal(response.Body.Bytes(), &body))
		assert.That(t, body.Error == "default branch does not exist")
	})
}

func testRepository(id uuid.UUID, ownerID uuid.UUID, name string, description *string, visibility db.RepositoryVisibility, defaultBranch string, createdAt time.Time, updatedAt time.Time) db.Repository {
	return db.Repository{
		ID:             UUIDToPgtype(id),
		OwnerID:        UUIDToPgtype(ownerID),
		Name:           name,
		Description:    description,
		Visibility:     visibility,
		DefaultBranch:  defaultBranch,
		CreatedAt:      pgtype.Timestamptz{Time: createdAt, Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: updatedAt, Valid: true},
		ItemPrefix:     "I",
		NextItemNumber: 1,
	}
}

func expectRepositoryByOwnerAndName(mock pgxmock.PgxPoolIface, ownerName string, repositoryName string, repository db.Repository) {
	mock.ExpectQuery("select repositories").
		WithArgs(ownerName, repositoryName).
		WillReturnRows(repositoryRows().AddRow(
			repository.ID,
			repository.OwnerID,
			repository.Name,
			repository.Description,
			repository.Visibility,
			repository.DefaultBranch,
			repository.CreatedAt,
			repository.UpdatedAt,
			repository.ItemPrefix,
			repository.NextItemNumber,
		))
}

func expectUserByName(mock pgxmock.PgxPoolIface, user db.User) {
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
			user.Description,
		))
}

func repositoryRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"owner_id",
		"name",
		"description",
		"visibility",
		"default_branch",
		"created_at",
		"updated_at",
		"item_prefix",
		"next_item_number",
	})
}

func requireGit(t *testing.T) string {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git binary is required for repository endpoint tests")
	}

	return gitBin
}

func ptr(value string) *string {
	return &value
}
