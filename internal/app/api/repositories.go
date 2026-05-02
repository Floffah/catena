package api

import (
	"context"
	"errors"
	"strings"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/jackc/pgx/v5"
)

func (s *Server) CreateRepository(ctx context.Context, request CreateRepositoryRequestObject) (CreateRepositoryResponseObject, error) {
	authUser, user, err := s.auth.EnsureUserInContext(ctx)
	if err != nil || authUser == nil {
		return CreateRepository401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	if request.Body == nil {
		return CreateRepository400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "request body is required"},
		}, nil
	}

	name := strings.TrimSpace(request.Body.Name)
	if name == "" {
		return CreateRepository400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "repository name is required"},
		}, nil
	}

	var description *string
	if request.Body.Description != nil {
		description = request.Body.Description
	}

	visibility := db.RepositoryVisibilityPrivate
	if request.Body.Visibility != nil {
		if !request.Body.Visibility.Valid() {
			return CreateRepository400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: "repository visibility is invalid"},
			}, nil
		}
		visibility = db.RepositoryVisibility(*request.Body.Visibility)
	}

	defaultBranch := "main"
	if request.Body.DefaultBranch != nil {
		defaultBranch = strings.TrimSpace(*request.Body.DefaultBranch)
	}
	if defaultBranch == "" {
		return CreateRepository400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "default branch is required"},
		}, nil
	}

	tx, err := s.dbConn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return CreateRepository500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to start transaction"},
		}, nil
	}
	queriesTx := s.repository.WithTx(tx)

	repository, err := queriesTx.CreateRepository(ctx, db.CreateRepositoryParams{
		OwnerID:       user.ID,
		Name:          name,
		Description:   description,
		Visibility:    visibility,
		DefaultBranch: defaultBranch,
	})
	if err != nil {
		tx.Rollback(ctx)

		if db.IsUniqueViolation(err) {
			return CreateRepository409JSONResponse{
				ConflictJSONResponse: ConflictJSONResponse{Error: "repository already exists"},
			}, nil
		}

		return CreateRepository500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to create repository"},
		}, nil
	}

	err = s.git.CreateRepo(repository)
	if err != nil {
		tx.Rollback(ctx)
		return CreateRepository500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to create git repository"},
		}, nil
	}

	repoResponse, err := RepositoryToAPI(repository)
	if err != nil {
		tx.Rollback(ctx)
		return CreateRepository500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode repository"},
		}, nil
	}

	err = tx.Commit(ctx)
	if err != nil {
		return CreateRepository500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to commit transaction"},
		}, nil
	}

	response := CreateRepository201JSONResponse{
		CreatedAt:     repoResponse.CreatedAt,
		DefaultBranch: repoResponse.DefaultBranch,
		Description:   repoResponse.Description,
		Id:            repoResponse.Id,
		Name:          repoResponse.Name,
		OwnerId:       repoResponse.OwnerId,
		OwnerName:     user.Name,
		UpdatedAt:     repoResponse.UpdatedAt,
		Visibility:    repoResponse.Visibility,
	}

	return response, nil
}

func (s *Server) GetRepositoryByOwnerAndName(ctx context.Context, request GetRepositoryByOwnerAndNameRequestObject) (GetRepositoryByOwnerAndNameResponseObject, error) {
	repository, err := s.repository.GetRepositoryByOwnerAndName(ctx, db.GetRepositoryByOwnerAndNameParams{
		OwnerName:      request.Owner,
		RepositoryName: request.Repository,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetRepositoryByOwnerAndName404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "repository not found"},
			}, nil
		}

		return GetRepositoryByOwnerAndName500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load repository"},
		}, nil
	}

	if repository.Visibility == db.RepositoryVisibilityPrivate {
		authUser, err := s.auth.GetUserFromContext(ctx)
		if err != nil {
			return GetRepositoryByOwnerAndName401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
			}, nil
		}
		if authUser == nil {
			return GetRepositoryByOwnerAndName401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
			}, nil
		}

		user, err := s.repository.GetUserByClerkUserID(ctx, authUser.ID)
		if err != nil || user.ID != repository.OwnerID {
			return GetRepositoryByOwnerAndName404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "repository not found"},
			}, nil
		}
	}

	response, err := RepositoryToAPI(repository)
	if err != nil {
		return GetRepositoryByOwnerAndName500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode repository"},
		}, nil
	}

	return GetRepositoryByOwnerAndName200JSONResponse(response), nil
}
