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
	authUser, err := s.auth.GetAuthFromContext(ctx)
	if err != nil || authUser == nil {
		return CreateRepository401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	user, err := s.auth.GetUserFromAuth(ctx, authUser)
	if err != nil {
		return CreateRepository500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
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
		_ = tx.Rollback(ctx)

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
		_ = tx.Rollback(ctx)
		return CreateRepository500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to create git repository"},
		}, nil
	}

	repoResponse, err := RepositoryToAPI(repository, user.Name)
	if err != nil {
		_ = tx.Rollback(ctx)
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

	response := CreateRepository201JSONResponse(repoResponse)

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

func (s *Server) ListUserRepositoriesByName(ctx context.Context, request ListUserRepositoriesByNameRequestObject) (ListUserRepositoriesByNameResponseObject, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return ListUserRepositoriesByName404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{Error: "user not found"},
		}, nil
	}

	limit, ok := normalizedListLimit(request.Params.Limit)
	if !ok {
		return ListUserRepositoriesByName400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "limit must be between 1 and 50"},
		}, nil
	}

	sort := Updated
	if request.Params.Sort != nil {
		if !request.Params.Sort.Valid() {
			return ListUserRepositoriesByName400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: "repository sort is invalid"},
			}, nil
		}

		sort = *request.Params.Sort
	}

	visibility := db.RepositoryVisibilityPublic
	filterVisibility := false
	if request.Params.Visibility != nil {
		if !request.Params.Visibility.Valid() {
			return ListUserRepositoriesByName400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: "repository visibility is invalid"},
			}, nil
		}

		visibility = db.RepositoryVisibility(*request.Params.Visibility)
		filterVisibility = true
	}

	owner, err := s.repository.GetUserByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ListUserRepositoriesByName404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "user not found"},
			}, nil
		}

		return ListUserRepositoriesByName500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}

	includePrivate := false
	authUser, err := s.auth.GetAuthFromContext(ctx)
	if err == nil && authUser != nil {
		user, err := s.auth.GetUserFromAuth(ctx, authUser)
		if err == nil {
			includePrivate = user.ID.Valid && user.ID == owner.ID
		}
	}

	var repositories []db.Repository
	switch sort {
	case Featured:
		repositories, err = s.repository.ListRepositoriesByOwnerFeatured(ctx, db.ListRepositoriesByOwnerFeaturedParams{
			OwnerID:          owner.ID,
			FilterVisibility: filterVisibility,
			IncludePrivate:   includePrivate,
			Visibility:       visibility,
			ResultLimit:      limit,
		})
	case Updated:
		repositories, err = s.repository.ListRepositoriesByOwnerUpdated(ctx, db.ListRepositoriesByOwnerUpdatedParams{
			OwnerID:          owner.ID,
			FilterVisibility: filterVisibility,
			IncludePrivate:   includePrivate,
			Visibility:       visibility,
			ResultLimit:      limit,
		})
	default:
		return ListUserRepositoriesByName400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "repository sort is invalid"},
		}, nil
	}
	if err != nil {
		return ListUserRepositoriesByName500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to list repositories"},
		}, nil
	}

	responseRepositories := make([]Repository, 0, len(repositories))
	for _, repository := range repositories {
		responseRepository, err := RepositoryToAPI(repository, owner.Name)
		if err != nil {
			return ListUserRepositoriesByName500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode repository"},
			}, nil
		}

		responseRepositories = append(responseRepositories, responseRepository)
	}

	return ListUserRepositoriesByName200JSONResponse{Repositories: responseRepositories}, nil
}

func (s *Server) UpdateRepository(ctx context.Context, request UpdateRepositoryRequestObject) (UpdateRepositoryResponseObject, error) {
	authUser, err := s.auth.GetAuthFromContext(ctx)
	if err != nil || authUser == nil {
		return UpdateRepository401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	user, err := s.auth.GetUserFromAuth(ctx, authUser)
	if err != nil {
		return UpdateRepository500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}

	repository, accessErr := s.getAccessibleRepositoryForUser(ctx, request.Owner, request.Repository, &user)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return UpdateRepository401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return UpdateRepository404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return UpdateRepository500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	if repository.OwnerID != user.ID {
		return UpdateRepository403JSONResponse{
			ForbiddenJSONResponse: ForbiddenJSONResponse{Error: "forbidden"},
		}, nil
	}

	if request.Body == nil {
		return UpdateRepository400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "request body is required"},
		}, nil
	}

	description := repository.Description
	if request.Body.Description != nil {
		trimmedDescription := strings.TrimSpace(*request.Body.Description)
		if trimmedDescription == "" {
			description = nil
		} else {
			description = &trimmedDescription
		}
	}

	visibility := repository.Visibility
	if request.Body.Visibility != nil {
		if !request.Body.Visibility.Valid() {
			return UpdateRepository400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: "repository visibility is invalid"},
			}, nil
		}

		visibility = db.RepositoryVisibility(*request.Body.Visibility)
	}

	defaultBranch := repository.DefaultBranch
	if request.Body.DefaultBranch != nil {
		defaultBranch = strings.TrimSpace(*request.Body.DefaultBranch)
	}
	if defaultBranch == "" {
		return UpdateRepository400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "default branch is required"},
		}, nil
	}
	if request.Body.DefaultBranch != nil && defaultBranch != repository.DefaultBranch {
		branchExists, err := s.git.BranchExists(ctx, repository, defaultBranch)
		if err != nil {
			if errors.Is(err, gitstore.ErrInvalidRef) {
				return UpdateRepository400JSONResponse{
					BadRequestJSONResponse: BadRequestJSONResponse{Error: "default branch is not a valid ref"},
				}, nil
			}

			return UpdateRepository500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to validate default branch"},
			}, nil
		}
		if !branchExists {
			return UpdateRepository400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: "default branch does not exist"},
			}, nil
		}
	}

	updatedRepository, err := s.repository.UpdateRepository(ctx, db.UpdateRepositoryParams{
		ID:            repository.ID,
		Name:          repository.Name,
		Description:   description,
		Visibility:    visibility,
		DefaultBranch: defaultBranch,
	})
	if err != nil {
		return UpdateRepository500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to update repository"},
		}, nil
	}

	response, err := RepositoryToAPI(updatedRepository, request.Owner)
	if err != nil {
		return UpdateRepository500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode repository"},
		}, nil
	}

	return UpdateRepository200JSONResponse(response), nil
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

func (s *Server) GetRepositoryFile(ctx context.Context, request GetRepositoryFileRequestObject) (GetRepositoryFileResponseObject, error) {
	repository, accessErr := s.getAccessibleRepository(ctx, request.Owner, request.Repository)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return GetRepositoryFile401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return GetRepositoryFile404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return GetRepositoryFile500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	ref := ""
	if request.Params.Ref != nil {
		ref = strings.TrimSpace(*request.Params.Ref)
	}

	filePath := strings.TrimSpace(request.Params.Path)

	file, err := s.git.GetFile(ctx, repository, ref, filePath)
	if err != nil {
		switch {
		case errors.Is(err, gitstore.ErrInvalidPath), errors.Is(err, gitstore.ErrInvalidRef):
			return GetRepositoryFile400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: err.Error()},
			}, nil
		case errors.Is(err, gitstore.ErrFileNotFound):
			return GetRepositoryFile404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "file not found"},
			}, nil
		case errors.Is(err, gitstore.ErrFileTooLarge):
			return GetRepositoryFile400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: "file is too large"},
			}, nil
		default:
			return GetRepositoryFile500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load file"},
			}, nil
		}
	}

	return GetRepositoryFile200JSONResponse{
		CommitOid: file.CommitOID,
		Content:   file.Content,
		Encoding:  "utf-8",
		Name:      file.Name,
		Oid:       file.OID,
		Path:      file.Path,
		Ref:       file.Ref,
		Size:      file.Size,
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

	recursive := request.Params.Recursive != nil && *request.Params.Recursive

	tree, err := s.git.GetTree(ctx, repository, ref, directory, recursive)
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
		case errors.Is(err, gitstore.ErrTreeTooLarge):
			return GetRepositoryTree413JSONResponse{
				PayloadTooLargeJSONResponse: PayloadTooLargeJSONResponse{Error: err.Error()},
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
	return s.getAccessibleRepositoryForUser(ctx, ownerName, repositoryName, nil)
}

func (s *Server) getAccessibleRepositoryForUser(ctx context.Context, ownerName string, repositoryName string, user *db.User) (db.Repository, *repositoryAccessError) {
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

	if user == nil {
		authUser, err := s.auth.GetAuthFromContext(ctx)
		if err != nil || authUser == nil {
			return db.Repository{}, &repositoryAccessError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			}
		}

		dbUser, err := s.auth.GetUserFromAuth(ctx, authUser)
		if err != nil {
			return db.Repository{}, &repositoryAccessError{
				Status:  http.StatusNotFound,
				Message: "repository not found",
			}
		}
		user = &dbUser
	}

	if user.ID != repository.OwnerID {
		return db.Repository{}, &repositoryAccessError{
			Status:  http.StatusNotFound,
			Message: "repository not found",
		}
	}

	return repository, nil
}
