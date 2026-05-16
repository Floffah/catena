package api

import "context"

func (s *Server) CreateClerkSignInToken(ctx context.Context, request CreateClerkSignInTokenRequestObject) (CreateClerkSignInTokenResponseObject, error) {
	return CreateClerkSignInToken403JSONResponse{
		ForbiddenJSONResponse{
			Error: "sign-in token creation is disabled",
		},
	}, nil
}
