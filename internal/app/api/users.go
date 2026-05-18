package api

import (
	"context"
	"errors"
	"strings"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/jackc/pgx/v5"
)

func (s *Server) GetAuthenticatedUser(ctx context.Context, request GetAuthenticatedUserRequestObject) (GetAuthenticatedUserResponseObject, error) {
	authUser, err := s.auth.GetAuthFromContext(ctx)
	if err != nil {
		return GetAuthenticatedUser500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load auth user"},
		}, nil
	}

	if authUser == nil {
		return GetAuthenticatedUser401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	user, err := s.auth.GetUserFromAuth(ctx, authUser)
	if err != nil {
		return GetAuthenticatedUser500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}

	if !user.ID.Valid {
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

func (s *Server) UpdateAuthenticatedUser(ctx context.Context, request UpdateAuthenticatedUserRequestObject) (UpdateAuthenticatedUserResponseObject, error) {
	authUser, err := s.auth.GetAuthFromContext(ctx)
	if err != nil {
		return UpdateAuthenticatedUser500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load auth user"},
		}, nil
	}

	if authUser == nil {
		return UpdateAuthenticatedUser401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	user, err := s.auth.GetUserFromAuth(ctx, authUser)
	if err != nil {
		return UpdateAuthenticatedUser500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}

	if !user.ID.Valid {
		return UpdateAuthenticatedUser401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	if request.Body == nil {
		return UpdateAuthenticatedUser400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "request body is required"},
		}, nil
	}

	if request.Body.DisplayName == nil {
		response, err := UserToAPI(user)
		if err != nil {
			return UpdateAuthenticatedUser500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode user"},
			}, nil
		}

		return UpdateAuthenticatedUser200JSONResponse(response), nil
	}

	trimmedDisplayName := strings.TrimSpace(*request.Body.DisplayName)
	if trimmedDisplayName == "" {
		return UpdateAuthenticatedUser400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "displayName must not be empty"},
		}, nil
	}

	user, err = s.repository.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		ID:          user.ID,
		Name:        user.Name,
		DisplayName: &trimmedDisplayName,
		AvatarUrl:   user.AvatarUrl,
	})
	if err != nil {
		return UpdateAuthenticatedUser500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to update user"},
		}, nil
	}

	response, err := UserToAPI(user)
	if err != nil {
		return UpdateAuthenticatedUser500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode user"},
		}, nil
	}

	return UpdateAuthenticatedUser200JSONResponse(response), nil
}

func (s *Server) GetUserByClerkUserId(ctx context.Context, request GetUserByClerkUserIdRequestObject) (GetUserByClerkUserIdResponseObject, error) {
	authUser, err := s.auth.GetAuthFromContext(ctx)
	if err != nil {
		return GetUserByClerkUserId500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load auth user"},
		}, nil
	}

	if authUser == nil {
		return GetUserByClerkUserId401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	if authUser.ClerkUserID != request.ClerkUserId {
		return GetUserByClerkUserId403JSONResponse{
			ForbiddenJSONResponse: ForbiddenJSONResponse{Error: "forbidden"},
		}, nil
	}

	user, err := s.auth.GetUserFromAuth(ctx, authUser)
	if err != nil {
		return GetUserByClerkUserId500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}

	if !user.ID.Valid {
		return GetUserByClerkUserId401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
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
