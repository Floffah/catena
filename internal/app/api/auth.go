package api

import "context"

func (s *Server) CreateClerkSignInToken(ctx context.Context, request CreateClerkSignInTokenRequestObject) (CreateClerkSignInTokenResponseObject, error) {
	_, user, err := s.auth.EnsureUserInContext(ctx)
	if err != nil {
		return CreateClerkSignInToken500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load user"},
		}, nil
	}
	if !user.ID.Valid {
		return CreateClerkSignInToken401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
		}, nil
	}

	token, err := s.auth.CreateClerkSignInToken(user)
	if err != nil {
		return CreateClerkSignInToken500JSONResponse{
			InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to create Clerk sign-in token"},
		}, nil
	}

	return CreateClerkSignInToken200JSONResponse{
		Token: token,
	}, nil
}
