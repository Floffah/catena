package api

import (
	"github.com/floffah/catena/internal/pkg/db"
	"github.com/google/uuid"
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
