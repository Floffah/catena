package api

import "context"

func (s *Server) CreateClerkSignInToken(ctx context.Context, request CreateClerkSignInTokenRequestObject) (CreateClerkSignInTokenResponseObject, error) {
	return CreateClerkSignInToken403JSONResponse{
		ForbiddenJSONResponse{
			Error: "sign-in token creation is disabled",
		},
	}, nil

	//authUser, err := s.auth.GetAuthFromContext(ctx)
	//if err != nil {
	//	return CreateClerkSignInToken500JSONResponse{
	//		InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to load auth user"},
	//	}, nil
	//}
	//if authUser == nil {
	//	return CreateClerkSignInToken401JSONResponse{
	//		UnauthorizedJSONResponse: UnauthorizedJSONResponse{Error: "unauthorized"},
	//	}, nil
	//}
	//
	//token, err := s.auth.CreateClerkSignInToken(authUser)
	//if err != nil {
	//	return CreateClerkSignInToken500JSONResponse{
	//		InternalServerErrorJSONResponse: InternalServerErrorJSONResponse{Error: "failed to create Clerk sign-in token"},
	//	}, nil
	//}
	//
	//return CreateClerkSignInToken200JSONResponse{
	//	Token: token,
	//}, nil
}
