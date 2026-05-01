package auth

import (
	"context"
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
)

const UserContextKey = "user"

type AuthService struct {
	ClerkJwks *jwks.Client
	ClerkUser *user.Client
}

func NewAuthService(clerkSecretKey string) AuthService {
	clerkConf := &clerk.ClientConfig{}
	clerkConf.Key = &clerkSecretKey
	clerkJwks := jwks.NewClient(clerkConf)
	clerkUser := user.NewClient(clerkConf)

	return AuthService{
		ClerkJwks: clerkJwks,
		ClerkUser: clerkUser,
	}
}

func (s AuthService) GetUserFromContext(ctx context.Context) (*clerk.User, error) {
	ginCtx, okCtx := ctx.(*gin.Context)
	if okCtx {
		cachedUser, exists := ginCtx.Get(UserContextKey)
		if exists {
			if authuser, ok := cachedUser.(*clerk.User); ok {
				return authuser, nil
			}

			return nil, nil
		}

		auth := ginCtx.GetHeader("Authorization")

		var token string
		if len(auth) > 7 && auth[:7] == "Bearer " {
			token = auth[7:]
		}

		if token != "" {
			claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
				Token:      token,
				JWKSClient: s.ClerkJwks,
			})
			if err != nil {
				return nil, err
			}

			authuser, err := s.ClerkUser.Get(ctx, claims.Subject)
			if err != nil {
				return nil, err
			}

			return authuser, nil
		}
	}

	return nil, nil
}

func (s AuthService) EnsureUserInContext(ctx context.Context) (*clerk.User, error) {
	authuser, err := s.GetUserFromContext(ctx)
	ginCtx, okCtx := ctx.(*gin.Context)

	if okCtx {
		if err != nil {
			ginCtx.AbortWithStatus(401)
			return nil, err
		}

		if authuser == nil {
			ginCtx.AbortWithStatus(401)
			return nil, fmt.Errorf("unauthorized")
		}
	} else {
		return nil, err
	}

	return authuser, nil
}

func (s AuthService) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authuser, err := s.GetUserFromContext(c)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		c.Set(UserContextKey, authuser)
		c.Next()
	}
}
