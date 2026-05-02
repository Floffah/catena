package api

import (
	"time"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func RepositoryToAPI(repository db.Repository) (Repository, error) {
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
		DisplayName: user.DisplayName,
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

func UUIDToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func timestampPtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}

	return &ts.Time
}
