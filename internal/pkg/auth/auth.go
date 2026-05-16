package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/signintoken"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/floffah/catena/internal/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

const AuthContextKey = "authUser"

var (
	ClerkPublishableKey string
	ClerkFrontendApiUrl string
)

type AuthService struct {
	repository       db.Queries
	ClerkJwks        *jwks.Client
	ClerkUser        *user.Client
	ClerkSignInToken *signintoken.Client
}

type Auth struct {
	ClerkUserID string
}

func NewAuthService(clerkSecretKey string, conn db.DBTX) *AuthService {
	clerkConf := &clerk.ClientConfig{}
	clerkConf.Key = &clerkSecretKey
	clerkJwks := jwks.NewClient(clerkConf)
	clerkUser := user.NewClient(clerkConf)
	clerkSIT := signintoken.NewClient(clerkConf)

	return &AuthService{
		repository:       *db.New(conn),
		ClerkJwks:        clerkJwks,
		ClerkUser:        clerkUser,
		ClerkSignInToken: clerkSIT,
	}
}

func (s *AuthService) GetAuthFromContext(ctx context.Context) (*Auth, error) {
	ginCtx, okCtx := ctx.(*gin.Context)
	if okCtx {
		cachedAuth, exists := ginCtx.Get(AuthContextKey)
		if exists {
			if auth, ok := cachedAuth.(*Auth); ok {
				return auth, nil
			}

			return nil, fmt.Errorf("cached auth in context is not of type *auth.Auth")
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

			return &Auth{ClerkUserID: claims.Subject}, nil
		}
	}

	return nil, nil
}

func (s *AuthService) GetUserFromAuth(ctx context.Context, auth *Auth) (db.User, error) {
	if auth == nil {
		return db.User{}, nil
	}

	dbuser, err := s.repository.GetUserByClerkUserID(ctx, auth.ClerkUserID)
	if err == nil {
		return dbuser, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return db.User{}, err
	}

	authUser, err := s.ClerkUser.Get(ctx, auth.ClerkUserID)
	if err != nil {
		return db.User{}, err
	}

	name := authUser.ID
	if authUser.Username != nil && strings.TrimSpace(*authUser.Username) != "" {
		name = strings.TrimSpace(*authUser.Username)
	}

	displayName := ""

	hasFirstName := authUser.FirstName != nil && strings.TrimSpace(*authUser.FirstName) != ""
	hasLastName := authUser.LastName != nil && strings.TrimSpace(*authUser.LastName) != ""

	if hasFirstName {
		displayName += strings.TrimSpace(*authUser.FirstName)
	}
	if hasFirstName && hasLastName {
		displayName += " "
	}
	if hasLastName {
		displayName += strings.TrimSpace(*authUser.LastName)
	}

	if !hasFirstName && !hasLastName {
		displayName = name
	}

	newUser, err := s.repository.CreateUser(ctx, db.CreateUserParams{
		ClerkUserID: authUser.ID,
		Name:        name,
		DisplayName: &displayName,
		AvatarUrl:   authUser.ImageURL,
	})

	if err != nil {
		return db.User{}, err
	}

	return newUser, nil
}

func (s *AuthService) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := s.GetAuthFromContext(c)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		if auth != nil {
			c.Set(AuthContextKey, auth)
		}
		c.Next()
	}
}

func (s *AuthService) CreateClerkSignInToken(auth *Auth) (string, error) {
	expiresInSeconds := int64(20)
	signInToken, err := s.ClerkSignInToken.Create(context.Background(), &signintoken.CreateParams{
		UserID:           &auth.ClerkUserID,
		ExpiresInSeconds: &expiresInSeconds,
	})
	if err != nil {
		return "", err
	}

	return signInToken.Token, nil
}
