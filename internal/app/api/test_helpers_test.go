package api

import (
	"context"
	"os"
	"testing"

	"github.com/floffah/catena/internal/pkg/auth"
	"github.com/floffah/catena/internal/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const testBearerToken = "test-user"

type testAuthProvider struct {
	user    db.User
	authErr error
	userErr error
}

func (p testAuthProvider) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "Bearer "+testBearerToken {
			c.Set(auth.AuthContextKey, &auth.Auth{ClerkUserID: p.user.ClerkUserID})
		}

		c.Next()
	}
}

func (p testAuthProvider) GetAuthFromContext(ctx context.Context) (*auth.Auth, error) {
	if p.authErr != nil {
		return nil, p.authErr
	}

	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil, nil
	}

	value, exists := ginCtx.Get(auth.AuthContextKey)
	if !exists {
		return nil, nil
	}

	authUser, ok := value.(*auth.Auth)
	if !ok {
		return nil, nil
	}

	return authUser, nil
}

func (p testAuthProvider) GetUserFromAuth(_ context.Context, authUser *auth.Auth) (db.User, error) {
	if p.userErr != nil {
		return db.User{}, p.userErr
	}

	if authUser == nil || authUser.ClerkUserID != p.user.ClerkUserID {
		return db.User{}, nil
	}

	return p.user, nil
}

func (p testAuthProvider) CreateClerkSignInToken(*auth.Auth) (string, error) {
	return "test-sign-in-token", nil
}

type failDB struct {
	t *testing.T
}

func (f failDB) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	f.t.Fatal("database Exec should not be called")
	return pgconn.CommandTag{}, nil
}

func (f failDB) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	f.t.Fatal("database Query should not be called")
	return nil, nil
}

func (f failDB) QueryRow(context.Context, string, ...interface{}) pgx.Row {
	f.t.Fatal("database QueryRow should not be called")
	return nil
}

func (f failDB) BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error) {
	f.t.Fatal("database BeginTx should not be called")
	return nil, nil
}

func isDir(t *testing.T, path string) bool {
	t.Helper()

	info, err := os.Stat(path)
	assertNoError(t, err)
	return info.IsDir()
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
