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

func (s *AuthService) GetAuthFromContext(ctx context.Context) (*clerk.User, error) {
	ginCtx, okCtx := ctx.(*gin.Context)
	if okCtx {
		cachedUser, exists := ginCtx.Get(AuthContextKey)
		if exists {
			if authuser, ok := cachedUser.(*clerk.User); ok {
				return authuser, nil
			}

			return nil, fmt.Errorf("cached auth in context is not of type *clerk.User")
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

func (s *AuthService) GetUserFromAuth(ctx context.Context, authUser *clerk.User) (db.User, error) {
	if authUser == nil {
		return db.User{}, nil
	}

	dbuser, err := s.repository.GetUserByClerkUserID(ctx, authUser.ID)
	if err == nil {
		return dbuser, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
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
		authuser, err := s.GetAuthFromContext(c)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		if authuser != nil {
			c.Set(AuthContextKey, authuser)
		}
		c.Next()
	}
}

func (s *AuthService) CreateClerkSignInToken(authUser *clerk.User) (string, error) {
	expiresInSeconds := int64(20)
	signInToken, err := s.ClerkSignInToken.Create(context.Background(), &signintoken.CreateParams{
		UserID:           &authUser.ID,
		ExpiresInSeconds: &expiresInSeconds,
	})
	if err != nil {
		return "", err
	}

	return signInToken.Token, nil
}
