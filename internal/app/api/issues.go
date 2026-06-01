package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/repositoryitems"
	"github.com/jackc/pgx/v5"
)

func (s *Server) ListRepositoryIssues(ctx context.Context, request ListRepositoryIssuesRequestObject) (ListRepositoryIssuesResponseObject, error) {
	limit, ok := normalizedListLimit(request.Params.Limit)
	if !ok {
		return ListRepositoryIssues400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "limit must be between 1 and 50"},
		}, nil
	}

	repository, accessErr := s.getAccessibleRepository(ctx, request.Owner, request.Repository)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return ListRepositoryIssues401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return ListRepositoryIssues404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return ListRepositoryIssues500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	rows, err := s.repository.ListIssuesByRepository(ctx, db.ListIssuesByRepositoryParams{
		RepositoryID: repository.ID,
		ResultLimit:  limit,
	})
	if err != nil {
		return ListRepositoryIssues500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to list issues"},
		}, nil
	}

	issues := make([]Issue, 0, len(rows))
	for _, row := range rows {
		issue, err := ListIssuesByRepositoryRowToAPI(repository, row)
		if err != nil {
			return ListRepositoryIssues500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode issue"},
			}, nil
		}

		issues = append(issues, issue)
	}

	return ListRepositoryIssues200JSONResponse{Issues: issues}, nil
}

func (s *Server) CreateRepositoryIssue(ctx context.Context, request CreateRepositoryIssueRequestObject) (CreateRepositoryIssueResponseObject, error) {
	repository, accessErr := s.getAccessibleRepository(ctx, request.Owner, request.Repository)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return CreateRepositoryIssue401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return CreateRepositoryIssue404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return CreateRepositoryIssue500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	authUser, err := s.auth.GetAuthFromContext(ctx)
	if err != nil || authUser == nil {
		return CreateRepositoryIssue401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	user, err := s.auth.GetUserFromAuth(ctx, authUser)
	if err != nil {
		return CreateRepositoryIssue500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}
	if !user.ID.Valid {
		return CreateRepositoryIssue401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	if request.Body == nil {
		return CreateRepositoryIssue400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "request body is required"},
		}, nil
	}

	title := strings.TrimSpace(request.Body.Title)
	if title == "" {
		return CreateRepositoryIssue400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "issue title is required"},
		}, nil
	}

	tx, err := s.dbConn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return CreateRepositoryIssue500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to start transaction"},
		}, nil
	}
	queriesTx := s.repository.WithTx(tx)
	itemService := repositoryitems.NewService(tx)

	item, err := itemService.Create(ctx, repositoryitems.CreateParams{
		Repository: repository,
		Kind:       db.RepositoryItemKindIssue,
		Title:      title,
		Body:       request.Body.Body,
		AuthorID:   user.ID,
	})
	if err != nil {
		_ = tx.Rollback(ctx)
		return CreateRepositoryIssue500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to create repository item"},
		}, nil
	}

	issue, err := queriesTx.CreateIssue(ctx, item.ID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return CreateRepositoryIssue500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: err.Error()},
		}, nil
	}

	response, err := IssueToAPI(repository, item, issue.Status)
	if err != nil {
		_ = tx.Rollback(ctx)
		return CreateRepositoryIssue500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode issue"},
		}, nil
	}

	err = tx.Commit(ctx)
	if err != nil {
		return CreateRepositoryIssue500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to commit transaction"},
		}, nil
	}

	return CreateRepositoryIssue201JSONResponse(response), nil
}

func (s *Server) GetRepositoryIssue(ctx context.Context, request GetRepositoryIssueRequestObject) (GetRepositoryIssueResponseObject, error) {
	if request.Number <= 0 {
		return GetRepositoryIssue400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "issue number is invalid"},
		}, nil
	}

	repository, accessErr := s.getAccessibleRepository(ctx, request.Owner, request.Repository)
	if accessErr != nil {
		switch accessErr.Status {
		case http.StatusUnauthorized:
			return GetRepositoryIssue401JSONResponse{
				UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: accessErr.Message},
			}, nil
		case http.StatusNotFound:
			return GetRepositoryIssue404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: accessErr.Message},
			}, nil
		default:
			return GetRepositoryIssue500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: accessErr.Message},
			}, nil
		}
	}

	row, err := s.repository.GetIssueByRepositoryAndNumber(ctx, db.GetIssueByRepositoryAndNumberParams{
		RepositoryID: repository.ID,
		Number:       request.Number,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetRepositoryIssue404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "issue not found"},
			}, nil
		}

		return GetRepositoryIssue500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load issue"},
		}, nil
	}

	issue, err := GetIssueByRepositoryAndNumberRowToAPI(repository, row)
	if err != nil {
		return GetRepositoryIssue500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode issue"},
		}, nil
	}

	return GetRepositoryIssue200JSONResponse(issue), nil
}
