package api

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/gitauth"
)

func (s *Server) ListGitAccessTokens(ctx context.Context, request ListGitAccessTokensRequestObject) (ListGitAccessTokensResponseObject, error) {
	authUser, err := s.auth.GetAuthFromContext(ctx)
	if err != nil {
		return ListGitAccessTokens500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load auth user"},
		}, nil
	}
	if authUser == nil {
		return ListGitAccessTokens401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	user, err := s.auth.GetUserFromAuth(ctx, authUser)
	if err != nil {
		return ListGitAccessTokens500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}
	if !user.ID.Valid {
		return ListGitAccessTokens401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	tokens, err := s.repository.ListGitAccessTokensByUserID(ctx, user.ID)
	if err != nil {
		return ListGitAccessTokens500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to list git access tokens"},
		}, nil
	}

	response := make([]GitAccessToken, 0, len(tokens))
	for _, token := range tokens {
		apiToken, err := GitAccessTokenToAPI(token)
		if err != nil {
			return ListGitAccessTokens500JSONResponse{
				InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode git access token"},
			}, nil
		}
		response = append(response, apiToken)
	}

	return ListGitAccessTokens200JSONResponse(response), nil
}

func (s *Server) CreateGitAccessToken(ctx context.Context, request CreateGitAccessTokenRequestObject) (CreateGitAccessTokenResponseObject, error) {
	authUser, err := s.auth.GetAuthFromContext(ctx)
	if err != nil {
		return CreateGitAccessToken500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load auth user"},
		}, nil
	}
	if authUser == nil {
		return CreateGitAccessToken401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	user, err := s.auth.GetUserFromAuth(ctx, authUser)
	if err != nil {
		return CreateGitAccessToken500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}
	if !user.ID.Valid {
		return CreateGitAccessToken401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	if request.Body == nil {
		return CreateGitAccessToken400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "request body is required"},
		}, nil
	}

	name := strings.TrimSpace(request.Body.Name)
	if name == "" {
		return CreateGitAccessToken400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "token name is required"},
		}, nil
	}

	var scopes []string
	if request.Body.Scopes != nil {
		scopes = *request.Body.Scopes
	}
	scopes = gitauth.NormalizeScopes(scopes)
	if err := gitauth.ValidateScopes(scopes); err != nil {
		return CreateGitAccessToken400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "token scopes are invalid"},
		}, nil
	}

	if request.Body.ExpiresAt != nil && request.Body.ExpiresAt.Before(time.Now()) {
		return CreateGitAccessToken400JSONResponse{
			BadRequestJSONResponse: BadRequestJSONResponse{Error: "expiresAt must be in the future"},
		}, nil
	}

	rawToken, token, err := s.gitTokens.CreateToken(ctx, user, name, scopes, request.Body.ExpiresAt)
	if err != nil {
		if errors.Is(err, gitauth.ErrInvalidScope) {
			return CreateGitAccessToken400JSONResponse{
				BadRequestJSONResponse: BadRequestJSONResponse{Error: "token scopes are invalid"},
			}, nil
		}

		return CreateGitAccessToken500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to create git access token"},
		}, nil
	}

	apiToken, err := GitAccessTokenToAPI(token)
	if err != nil {
		return CreateGitAccessToken500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to encode git access token"},
		}, nil
	}

	return CreateGitAccessToken201JSONResponse{
		AccessToken: apiToken,
		Token:       rawToken,
	}, nil
}

func (s *Server) RevokeGitAccessToken(ctx context.Context, request RevokeGitAccessTokenRequestObject) (RevokeGitAccessTokenResponseObject, error) {
	authUser, err := s.auth.GetAuthFromContext(ctx)
	if err != nil {
		return RevokeGitAccessToken500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load auth user"},
		}, nil
	}
	if authUser == nil {
		return RevokeGitAccessToken401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	user, err := s.auth.GetUserFromAuth(ctx, authUser)
	if err != nil {
		return RevokeGitAccessToken500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}
	if !user.ID.Valid {
		return RevokeGitAccessToken401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	err = s.repository.RevokeGitAccessToken(ctx, db.RevokeGitAccessTokenParams{
		ID:     UUIDToPgtype(request.Id),
		UserID: user.ID,
	})
	if err != nil {
		return RevokeGitAccessToken500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to revoke git access token"},
		}, nil
	}

	return RevokeGitAccessToken204Response{}, nil
}
