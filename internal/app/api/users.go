package api

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (s *Server) GetAuthenticatedUser(ctx context.Context, request GetAuthenticatedUserRequestObject) (GetAuthenticatedUserResponseObject, error) {
	_, user, err := s.auth.EnsureUserInContext(ctx)
	if err != nil {
		return GetAuthenticatedUser500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}

	if user.ID.Valid == false {
		return GetAuthenticatedUser401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	response, err := UserToAPI(user)
	if err != nil {
		return GetAuthenticatedUser500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode user"},
		}, nil
	}

	return GetAuthenticatedUser200JSONResponse(response), nil
}

func (s *Server) GetUserByClerkUserId(ctx context.Context, request GetUserByClerkUserIdRequestObject) (GetUserByClerkUserIdResponseObject, error) {
	authUser, user, err := s.auth.EnsureUserInContext(ctx)
	if err != nil {
		return GetUserByClerkUserId500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}

	if authUser == nil || user.ID.Valid == false {
		return GetUserByClerkUserId401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	if authUser.ID != request.ClerkUserId {
		return GetUserByClerkUserId403JSONResponse{
			ForbiddenJSONResponse: ForbiddenJSONResponse{Error: "forbidden"},
		}, nil
	}

	response, err := UserToAPI(user)
	if err != nil {
		return GetUserByClerkUserId500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode user"},
		}, nil
	}

	return GetUserByClerkUserId200JSONResponse(response), nil
}

func (s *Server) GetUserByName(ctx context.Context, request GetUserByNameRequestObject) (GetUserByNameResponseObject, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return GetUserByName404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{Error: "user not found"},
		}, nil
	}

	user, err := s.repository.GetUserByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetUserByName404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{Error: "user not found"},
			}, nil
		}

		return GetUserByName500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}

	response, err := UserToAPI(user)
	if err != nil {
		return GetUserByName500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode user"},
		}, nil
	}

	return GetUserByName200JSONResponse(response), nil
}
