package api

import (
	"context"
	"net/http"
	"os"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/floffah/catena/api"
	"github.com/floffah/catena/internal/pkg/auth"
	"github.com/floffah/catena/internal/pkg/db"
	http2 "github.com/floffah/catena/internal/pkg/http"
	"github.com/gin-gonic/gin"
)

//go:generate go tool oapi-codegen -config cfg.yaml ../../../api/api.v1.openapi.yaml

type Server struct {
	repository db.Queries
	clerkUser  *user.Client
	clerkJwks  *jwks.Client
}

func NewServer() (*http.Server, error) {
	ctx := context.Background()

	dbConn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	clerkConf := &clerk.ClientConfig{}
	key := os.Getenv("CLERK_SECRET_KEY")
	clerkConf.Key = &key
	clerkUser := user.NewClient(clerkConf)
	clerkJwks := jwks.NewClient(clerkConf)

	server := Server{
		repository: *db.New(dbConn),
		clerkUser:  clerkUser,
		clerkJwks:  clerkJwks,
	}
	strictServer := NewStrictHandler(&server, []StrictMiddlewareFunc{})

	r := gin.Default()
	r.Use(http2.CorsMiddleware())
	r.Handle("GET", "/docs", func(c *gin.Context) {
		html, err := scalargo.NewV2(
			scalargo.WithSpecBytes(api.V1ApiSpec),
			scalargo.WithTheme(scalargo.ThemeBluePlanet),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("Catena API"),
			),
		)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.Header("Content-Type", "text/html")
		c.String(200, html)
	})

	RegisterHandlers(r, strictServer)

	s := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}

	return s, nil
}

func (s *Server) Healthz(ctx context.Context, request HealthzRequestObject) (HealthzResponseObject, error) {
	ok := "ok"

	user, _ := auth.GetUserFromContext(ctx, s.clerkJwks, s.clerkUser)

	if user != nil {
		ok = "ok - authenticated user: " + user.ID
	}

	return Healthz200JSONResponse{Status: &ok}, nil
}
