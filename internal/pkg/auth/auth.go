package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
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
		if !exists {
			return nil, nil
		}

		auth, ok := cachedAuth.(*Auth)
		if !ok {
			return nil, fmt.Errorf("cached auth in context is not of type *auth.Auth")
		}

		return auth, nil
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

	primaryEmail := "catenauser+" + authUser.ID + "@oncatena.com"

	if authUser.PrimaryEmailAddressID != nil {
		for _, email := range authUser.EmailAddresses {
			if email.ID == *authUser.PrimaryEmailAddressID {
				primaryEmail = email.EmailAddress
				break
			}
		}
	}

	newUser, err := s.repository.CreateUser(ctx, db.CreateUserParams{
		ClerkUserID: authUser.ID,
		Name:        name,
		DisplayName: &displayName,
		AvatarUrl:   authUser.ImageURL,
		Email:       primaryEmail,
	})

	if err != nil {
		return db.User{}, err
	}

	return newUser, nil
}

func (s *AuthService) Middleware() gin.HandlerFunc {
	clerkMiddleware := clerkhttp.WithHeaderAuthorization(
		clerkhttp.JWKSClient(s.ClerkJwks),
		clerkhttp.AuthorizationJWTExtractor(func(r *http.Request) string {
			header := strings.TrimSpace(r.Header.Get("Authorization"))
			if !strings.HasPrefix(header, "Bearer ") {
				return ""
			}

			return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		}),
	)

	return func(c *gin.Context) {
		hasBearerAuth := strings.HasPrefix(strings.TrimSpace(c.GetHeader("Authorization")), "Bearer ")
		nextCalled := false

		clerkMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			nextCalled = true
			c.Request = r

			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if hasBearerAuth && (!ok || claims == nil) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if ok && claims != nil {
				c.Set(AuthContextKey, &Auth{ClerkUserID: claims.Subject})
			}

			c.Next()
		})).ServeHTTP(c.Writer, c.Request)

		if !nextCalled {
			c.Abort()
		}
	}
}

func (s *AuthService) CreateClerkSignInToken(auth *Auth) (string, error) {
	signInToken, err := s.ClerkSignInToken.Create(context.Background(), &signintoken.CreateParams{
		UserID:           &auth.ClerkUserID,
		ExpiresInSeconds: new(int64(20)),
	})
	if err != nil {
		return "", err
	}

	return signInToken.Token, nil
}
