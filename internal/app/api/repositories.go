package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/gitstore"
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

	repoResponse, err := RepositoryToAPI(repository, user.Name)
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
		OwnerName:     repoResponse.OwnerName,
		UpdatedAt:     repoResponse.UpdatedAt,
		Visibility:    repoResponse.Visibility,
	}

	return response, nil
}

func (s *Server) GetRepositoryByOwnerAndName(ctx context.Context, request GetRepositoryByOwnerAndNameRequestObject) (GetRepositoryByOwnerAndNameResponseObject, error) {
	repository, accessErr := s.getAccessibleRepository(ctx, request.Owner, request.Repository)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return GetRepositoryByOwnerAndName401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return GetRepositoryByOwnerAndName404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return GetRepositoryByOwnerAndName500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	response, err := RepositoryToAPI(repository, request.Owner)
	if err != nil {
		return GetRepositoryByOwnerAndName500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode repository"},
		}, nil
	}

	return GetRepositoryByOwnerAndName200JSONResponse(response), nil
}

func (s *Server) GetRepositoryReadme(ctx context.Context, request GetRepositoryReadmeRequestObject) (GetRepositoryReadmeResponseObject, error) {
	repository, accessErr := s.getAccessibleRepository(ctx, request.Owner, request.Repository)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return GetRepositoryReadme401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return GetRepositoryReadme404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return GetRepositoryReadme500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	ref := ""
	if request.Params.Ref != nil {
		ref = strings.TrimSpace(*request.Params.Ref)
	}

	directory := ""
	if request.Params.Path != nil {
		directory = strings.TrimSpace(*request.Params.Path)
	}

	readme, err := s.git.GetReadme(ctx, repository, ref, directory)
	if err != nil {
		switch {
		case errors.Is(err, gitstore.ErrInvalidPath), errors.Is(err, gitstore.ErrInvalidRef):
			return GetRepositoryReadme400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: err.Error()},
			}, nil
		case errors.Is(err, gitstore.ErrReadmeNotFound):
			return GetRepositoryReadme404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "readme not found"},
			}, nil
		case errors.Is(err, gitstore.ErrReadmeTooLarge):
			return GetRepositoryReadme400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: "readme is too large"},
			}, nil
		default:
			return GetRepositoryReadme500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load readme"},
			}, nil
		}
	}

	return GetRepositoryReadme200JSONResponse{
		CommitOid: readme.CommitOID,
		Content:   readme.Content,
		Encoding:  "utf-8",
		Name:      readme.Name,
		Oid:       readme.OID,
		Path:      readme.Path,
		Ref:       readme.Ref,
		Size:      readme.Size,
	}, nil
}

func (s *Server) GetRepositoryLatestCommit(ctx context.Context, request GetRepositoryLatestCommitRequestObject) (GetRepositoryLatestCommitResponseObject, error) {
	repository, accessErr := s.getAccessibleRepository(ctx, request.Owner, request.Repository)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return GetRepositoryLatestCommit401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return GetRepositoryLatestCommit404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return GetRepositoryLatestCommit500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	ref := ""
	if request.Params.Ref != nil {
		ref = strings.TrimSpace(*request.Params.Ref)
	}

	path := ""
	if request.Params.Path != nil {
		path = strings.TrimSpace(*request.Params.Path)
	}

	commit, err := s.git.GetLatestCommit(ctx, repository, ref, path)
	if err != nil {
		switch {
		case errors.Is(err, gitstore.ErrInvalidPath), errors.Is(err, gitstore.ErrInvalidRef):
			return GetRepositoryLatestCommit400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: err.Error()},
			}, nil
		case errors.Is(err, gitstore.ErrCommitNotFound):
			return GetRepositoryLatestCommit404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "commit not found"},
			}, nil
		default:
			return GetRepositoryLatestCommit500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load latest commit"},
			}, nil
		}
	}

	return GetRepositoryLatestCommit200JSONResponse{
		AuthorEmail:     commit.AuthorEmail,
		AuthorName:      commit.AuthorName,
		AuthoredAt:      commit.AuthoredAt,
		CommitOid:       commit.CommitOID,
		CommittedAt:     commit.CommittedAt,
		CommitterEmail:  commit.CommitterEmail,
		CommitterName:   commit.CommitterName,
		Message:         commit.Message,
		MessageHeadline: commit.MessageHeadline,
		Ref:             commit.Ref,
		ShortOid:        commit.ShortOID,
	}, nil
}

func (s *Server) ResolveRepositoryGitPath(ctx context.Context, request ResolveRepositoryGitPathRequestObject) (ResolveRepositoryGitPathResponseObject, error) {
	repository, accessErr := s.getAccessibleRepository(ctx, request.Owner, request.Repository)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return ResolveRepositoryGitPath401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return ResolveRepositoryGitPath404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return ResolveRepositoryGitPath500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	resolved, err := s.git.ResolveGitPath(ctx, repository, request.Params.Path)
	if err != nil {
		switch {
		case errors.Is(err, gitstore.ErrInvalidPath):
			return ResolveRepositoryGitPath400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: err.Error()},
			}, nil
		case errors.Is(err, gitstore.ErrRefNotFound):
			return ResolveRepositoryGitPath404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "ref not found"},
			}, nil
		case errors.Is(err, gitstore.ErrPathNotFound):
			return ResolveRepositoryGitPath404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "path not found"},
			}, nil
		default:
			return ResolveRepositoryGitPath500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to resolve git path"},
			}, nil
		}
	}

	return ResolveRepositoryGitPath200JSONResponse{
		CommitOid: resolved.CommitOID,
		Path:      resolved.Path,
		PathType:  ResolvedRepositoryGitPathPathType(resolved.PathType),
		Ref:       resolved.Ref,
	}, nil
}

func (s *Server) ListRepositoryRefs(ctx context.Context, request ListRepositoryRefsRequestObject) (ListRepositoryRefsResponseObject, error) {
	repository, accessErr := s.getAccessibleRepository(ctx, request.Owner, request.Repository)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return ListRepositoryRefs401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return ListRepositoryRefs404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return ListRepositoryRefs500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	refType := Branch
	if request.Params.Type != nil {
		refType = *request.Params.Type
	}
	if refType != Branch {
		return ListRepositoryRefs400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "only branch refs are supported"},
		}, nil
	}

	refs, err := s.git.ListBranchRefs(ctx, repository)
	if err != nil {
		return ListRepositoryRefs500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load refs"},
		}, nil
	}

	responseRefs := make([]RepositoryRef, 0, len(refs))
	for _, ref := range refs {
		responseRefs = append(responseRefs, RepositoryRef{
			CommitOid: ref.CommitOID,
			IsDefault: ref.IsDefault,
			Name:      ref.Name,
			Type:      RepositoryRefType(ref.Type),
		})
	}

	return ListRepositoryRefs200JSONResponse{
		Refs: responseRefs,
	}, nil
}

func (s *Server) GetRepositoryTree(ctx context.Context, request GetRepositoryTreeRequestObject) (GetRepositoryTreeResponseObject, error) {
	repository, accessErr := s.getAccessibleRepository(ctx, request.Owner, request.Repository)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return GetRepositoryTree401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return GetRepositoryTree404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return GetRepositoryTree500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	ref := ""
	if request.Params.Ref != nil {
		ref = strings.TrimSpace(*request.Params.Ref)
	}

	directory := ""
	if request.Params.Path != nil {
		directory = strings.TrimSpace(*request.Params.Path)
	}

	tree, err := s.git.GetTree(ctx, repository, ref, directory)
	if err != nil {
		switch {
		case errors.Is(err, gitstore.ErrInvalidPath), errors.Is(err, gitstore.ErrInvalidRef):
			return GetRepositoryTree400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: err.Error()},
			}, nil
		case errors.Is(err, gitstore.ErrTreeNotFound):
			return GetRepositoryTree404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "tree not found"},
			}, nil
		default:
			return GetRepositoryTree500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load tree"},
			}, nil
		}
	}

	entries := make([]RepositoryTreeEntry, 0, len(tree.Entries))
	for _, entry := range tree.Entries {
		entries = append(entries, RepositoryTreeEntry{
			Name: entry.Name,
			Oid:  entry.OID,
			Path: entry.Path,
			Size: entry.Size,
			Type: RepositoryTreeEntryType(entry.Type),
		})
	}

	return GetRepositoryTree200JSONResponse{
		CommitOid: tree.CommitOID,
		Entries:   entries,
		Path:      tree.Path,
		Ref:       tree.Ref,
	}, nil
}

type repositoryAccessError struct {
	Status  int
	Message string
}

func (s *Server) getAccessibleRepository(ctx context.Context, ownerName string, repositoryName string) (db.Repository, *repositoryAccessError) {
	repository, err := s.repository.GetRepositoryByOwnerAndName(ctx, db.GetRepositoryByOwnerAndNameParams{
		OwnerName:      ownerName,
		RepositoryName: repositoryName,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Repository{}, &repositoryAccessError{
				Status:  http.StatusNotFound,
				Message: "repository not found",
			}
		}

		return db.Repository{}, &repositoryAccessError{
			Status:  http.StatusInternalServerError,
			Message: "failed to load repository",
		}
	}

	if repository.Visibility != db.RepositoryVisibilityPrivate {
		return repository, nil
	}

	authUser, err := s.auth.GetUserFromContext(ctx)
	if err != nil || authUser == nil {
		return db.Repository{}, &repositoryAccessError{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		}
	}

	user, err := s.repository.GetUserByClerkUserID(ctx, authUser.ID)
	if err != nil || user.ID != repository.OwnerID {
		return db.Repository{}, &repositoryAccessError{
			Status:  http.StatusNotFound,
			Message: "repository not found",
		}
	}

	return repository, nil
}
