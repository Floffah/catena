package auth

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(ctx context.Context, clerkJwks *jwks.Client, clerkUser *user.Client) (*clerk.User, error) {
	ginCtx, okCtx := ctx.(*gin.Context)
	if okCtx {
		auth := ginCtx.GetHeader("Authorization")

		var token string
		if len(auth) > 7 && auth[:7] == "Bearer " {
			token = auth[7:]
		}

		if token != "" {
			claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
				Token:      token,
				JWKSClient: clerkJwks,
			})
			if err != nil {
				return nil, err
			}

			authuser, err := clerkUser.Get(ctx, claims.Subject)
			if err != nil {
				return nil, err
			}

			return authuser, nil
		}
	}

	return nil, nil
}
