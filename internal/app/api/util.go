package api

import (
	"time"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/repositoryitems"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const maxListLimit = 50

func normalizedListLimit(limit *int) (int32, bool) {
	if limit == nil {
		return maxListLimit, true
	}

	if *limit < 1 || *limit > maxListLimit {
		return 0, false
	}

	return int32(*limit), true
}

func RepositoryToAPI(repository db.Repository, ownerName string) (Repository, error) {
	id, err := uuid.FromBytes(repository.ID.Bytes[:])
	if err != nil {
		return Repository{}, err
	}

	ownerID, err := uuid.FromBytes(repository.OwnerID.Bytes[:])
	if err != nil {
		return Repository{}, err
	}

	return Repository{
		CreatedAt:     repository.CreatedAt.Time,
		DefaultBranch: repository.DefaultBranch,
		Description:   repository.Description,
		Id:            id,
		Name:          repository.Name,
		OwnerId:       ownerID,
		OwnerName:     ownerName,
		UpdatedAt:     repository.UpdatedAt.Time,
		Visibility:    RepositoryVisibility(repository.Visibility),
	}, nil
}

func UserToAPI(user db.User) (User, error) {
	id, err := uuid.FromBytes(user.ID.Bytes[:])
	if err != nil {
		return User{}, err
	}

	return User{
		AvatarUrl:   user.AvatarUrl,
		CreatedAt:   user.CreatedAt.Time,
		Description: user.Description,
		DisplayName: user.DisplayName,
		Email:       new(openapi_types.Email(user.Email)),
		Id:          id,
		Name:        user.Name,
		UpdatedAt:   user.UpdatedAt.Time,
	}, nil
}

func GitAccessTokenToAPI(token db.GitAccessToken) (GitAccessToken, error) {
	id, err := uuid.FromBytes(token.ID.Bytes[:])
	if err != nil {
		return GitAccessToken{}, err
	}

	return GitAccessToken{
		CreatedAt:   token.CreatedAt.Time,
		ExpiresAt:   timestampPtr(token.ExpiresAt),
		Id:          id,
		LastUsedAt:  timestampPtr(token.LastUsedAt),
		Name:        token.Name,
		RevokedAt:   timestampPtr(token.RevokedAt),
		Scopes:      token.Scopes,
		TokenPrefix: token.TokenPrefix,
		UpdatedAt:   token.UpdatedAt.Time,
	}, nil
}

func IssueToAPI(repository db.Repository, item db.RepositoryItem, status db.IssueStatus) (Issue, error) {
	id, err := uuid.FromBytes(item.ID.Bytes[:])
	if err != nil {
		return Issue{}, err
	}

	repositoryID, err := uuid.FromBytes(item.RepositoryID.Bytes[:])
	if err != nil {
		return Issue{}, err
	}

	authorID, err := nullableUUIDPtr(item.AuthorID)
	if err != nil {
		return Issue{}, err
	}

	return Issue{
		AuthorId:       authorID,
		Body:           item.Body,
		CreatedAt:      item.CreatedAt.Time,
		Id:             id,
		Kind:           IssueKindIssue,
		LastActivityAt: item.LastActivityAt.Time,
		Number:         item.Number,
		Reference:      repositoryitems.Reference(repository, item),
		RepositoryId:   repositoryID,
		Status:         IssueStatus(status),
		Title:          item.Title,
		UpdatedAt:      item.UpdatedAt.Time,
	}, nil
}

func GetIssueByRepositoryAndNumberRowToAPI(repository db.Repository, row db.GetIssueByRepositoryAndNumberRow) (Issue, error) {
	return IssueToAPI(repository, db.RepositoryItem{
		ID:             row.ID,
		RepositoryID:   row.RepositoryID,
		Number:         row.Number,
		Kind:           row.Kind,
		Title:          row.Title,
		Body:           row.Body,
		AuthorID:       row.AuthorID,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
		LastActivityAt: row.LastActivityAt,
	}, row.Status)
}

func ListIssuesByRepositoryRowToAPI(repository db.Repository, row db.ListIssuesByRepositoryRow) (Issue, error) {
	return IssueToAPI(repository, db.RepositoryItem{
		ID:             row.ID,
		RepositoryID:   row.RepositoryID,
		Number:         row.Number,
		Kind:           row.Kind,
		Title:          row.Title,
		Body:           row.Body,
		AuthorID:       row.AuthorID,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
		LastActivityAt: row.LastActivityAt,
	}, row.Status)
}

func UUIDToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func nullableUUIDPtr(id pgtype.UUID) (*openapi_types.UUID, error) {
	if !id.Valid {
		return nil, nil
	}

	parsed, err := uuid.FromBytes(id.Bytes[:])
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func timestampPtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}

	return &ts.Time
}
