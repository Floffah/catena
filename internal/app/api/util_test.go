package api

import (
	"testing"
	"time"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/gitauth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/runtime/types"
	"github.com/zeebo/assert"
)

func TestRepositoryToAPI(t *testing.T) {
	repositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13af1")
	ownerID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13af2")
	description := "A tiny repository"
	createdAt := time.Date(2026, 5, 21, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)

	got, err := RepositoryToAPI(db.Repository{
		ID:            UUIDToPgtype(repositoryID),
		OwnerID:       UUIDToPgtype(ownerID),
		Name:          "catena",
		Description:   &description,
		Visibility:    db.RepositoryVisibilityPublic,
		DefaultBranch: "main",
		CreatedAt:     pgtype.Timestamptz{Time: createdAt, Valid: true},
		UpdatedAt:     pgtype.Timestamptz{Time: updatedAt, Valid: true},
	}, "floffah")

	assert.Nil(t, err)
	assert.That(t, got.Id == repositoryID)
	assert.That(t, got.OwnerId == ownerID)
	assert.That(t, got.OwnerName == "floffah")
	assert.That(t, got.Name == "catena")
	assert.That(t, got.Description != nil && *got.Description == description)
	assert.That(t, got.Visibility == Public)
	assert.That(t, got.DefaultBranch == "main")
	assert.That(t, got.CreatedAt.Equal(createdAt))
	assert.That(t, got.UpdatedAt.Equal(updatedAt))
}

func TestUserToAPI(t *testing.T) {
	userID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13af3")
	displayName := "Ramsay"
	description := "Building Catena"
	avatarURL := "https://example.com/avatar.png"
	createdAt := time.Date(2026, 5, 21, 11, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)

	got, err := UserToAPI(db.User{
		ID:          UUIDToPgtype(userID),
		Name:        "floffah",
		DisplayName: &displayName,
		Description: &description,
		AvatarUrl:   &avatarURL,
		Email:       "floffah@example.com",
		CreatedAt:   pgtype.Timestamptz{Time: createdAt, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: updatedAt, Valid: true},
	})

	assert.Nil(t, err)
	assert.That(t, got.Id == userID)
	assert.That(t, got.Name == "floffah")
	assert.That(t, got.DisplayName != nil && *got.DisplayName == displayName)
	assert.That(t, got.Description != nil && *got.Description == description)
	assert.That(t, got.AvatarUrl != nil && *got.AvatarUrl == avatarURL)
	assert.That(t, got.Email != nil && *got.Email == types.Email("floffah@example.com"))
	assert.That(t, got.CreatedAt.Equal(createdAt))
	assert.That(t, got.UpdatedAt.Equal(updatedAt))
}

func TestGitAccessTokenToAPI(t *testing.T) {
	tokenID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13af4")
	createdAt := time.Date(2026, 5, 21, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	lastUsedAt := createdAt.Add(10 * time.Minute)
	expiresAt := createdAt.Add(24 * time.Hour)

	got, err := GitAccessTokenToAPI(db.GitAccessToken{
		ID:          UUIDToPgtype(tokenID),
		Name:        "Local laptop",
		TokenPrefix: "ctn_pat_12345678",
		Scopes:      []string{gitauth.ScopeRepoRead, gitauth.ScopeRepoWrite},
		LastUsedAt:  pgtype.Timestamptz{Time: lastUsedAt, Valid: true},
		ExpiresAt:   pgtype.Timestamptz{Time: expiresAt, Valid: true},
		CreatedAt:   pgtype.Timestamptz{Time: createdAt, Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: updatedAt, Valid: true},
	})

	assert.Nil(t, err)
	assert.That(t, got.Id == tokenID)
	assert.That(t, got.Name == "Local laptop")
	assert.That(t, got.TokenPrefix == "ctn_pat_12345678")
	assert.That(t, got.LastUsedAt != nil && got.LastUsedAt.Equal(lastUsedAt))
	assert.That(t, got.ExpiresAt != nil && got.ExpiresAt.Equal(expiresAt))
	assert.That(t, got.RevokedAt == nil)
	assert.That(t, len(got.Scopes) == 2)
	assert.That(t, got.Scopes[0] == gitauth.ScopeRepoRead)
	assert.That(t, got.Scopes[1] == gitauth.ScopeRepoWrite)
	assert.That(t, got.CreatedAt.Equal(createdAt))
	assert.That(t, got.UpdatedAt.Equal(updatedAt))
}

func TestIssueToAPI(t *testing.T) {
	repositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13af5")
	itemID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13af6")
	authorID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13af7")
	body := "Issue body"
	createdAt := time.Date(2026, 5, 21, 13, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	lastActivityAt := createdAt.Add(2 * time.Hour)

	got, err := IssueToAPI(db.Repository{
		ID:         UUIDToPgtype(repositoryID),
		ItemPrefix: "CAT",
	}, db.RepositoryItem{
		ID:             UUIDToPgtype(itemID),
		RepositoryID:   UUIDToPgtype(repositoryID),
		Number:         42,
		Kind:           db.RepositoryItemKindIssue,
		Title:          "Test issue",
		Body:           &body,
		AuthorID:       UUIDToPgtype(authorID),
		CreatedAt:      pgtype.Timestamptz{Time: createdAt, Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: updatedAt, Valid: true},
		LastActivityAt: pgtype.Timestamptz{Time: lastActivityAt, Valid: true},
	}, db.IssueStatusInProgress)

	assert.Nil(t, err)
	assert.That(t, got.Id == itemID)
	assert.That(t, got.RepositoryId == repositoryID)
	assert.That(t, got.AuthorId != nil && *got.AuthorId == authorID)
	assert.That(t, got.Kind == IssueKindIssue)
	assert.That(t, got.Status == IssueStatusInProgress)
	assert.That(t, got.Reference == "CAT-42")
	assert.That(t, got.Title == "Test issue")
	assert.That(t, got.Body != nil && *got.Body == body)
	assert.That(t, got.Number == 42)
	assert.That(t, got.CreatedAt.Equal(createdAt))
	assert.That(t, got.UpdatedAt.Equal(updatedAt))
	assert.That(t, got.LastActivityAt.Equal(lastActivityAt))
}

func TestIssueToAPIWithDeletedAuthor(t *testing.T) {
	repositoryID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13af8")
	itemID := uuid.MustParse("019deb10-dafc-743f-8cfc-289a80c13af9")

	got, err := IssueToAPI(db.Repository{
		ID:         UUIDToPgtype(repositoryID),
		ItemPrefix: "I",
	}, db.RepositoryItem{
		ID:             UUIDToPgtype(itemID),
		RepositoryID:   UUIDToPgtype(repositoryID),
		Number:         1,
		Kind:           db.RepositoryItemKindIssue,
		Title:          "Deleted author issue",
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		LastActivityAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}, db.IssueStatusOpen)

	assert.Nil(t, err)
	assert.That(t, got.AuthorId == nil)
	assert.That(t, got.Reference == "I-1")
}
