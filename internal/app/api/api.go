package api

import (
	"context"
	"net/http"
	"time"

	scalargo "github.com/bdpiprava/scalar-go"
	catena "github.com/floffah/catena"
	"github.com/floffah/catena/api"
	"github.com/floffah/catena/internal/app/gitserver"
	"github.com/floffah/catena/internal/pkg/auth"
	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/gitauth"
	"github.com/floffah/catena/internal/pkg/gitstore"
	"github.com/floffah/catena/internal/pkg/httputil"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml ../../../api/api.v1.openapi.yaml

type Server struct {
	repository db.Queries
	dbConn     db.TxDB
	auth       auth.Provider
	git        gitstore.Store
	gitTokens  gitauth.TokenIssuer
}

type ServerDeps struct {
	DB                    db.TxDB
	Auth                  auth.Provider
	Git                   gitstore.Store
	GitTokens             gitauth.TokenIssuer
	GitCredentialVerifier gitauth.GitCredentialVerifier
	CorsAllowedOrigins    []string
}

func NewServer(
	conn *pgxpool.Pool,
	authService *auth.Service,
	gitService gitstore.Store,
	corsAllowedOrigins []string,
	port string,
) (*http.Server, error) {
	gitAuthService := gitauth.NewService(conn)
	r := NewRouter(ServerDeps{
		DB:                    conn,
		Auth:                  authService,
		Git:                   gitService,
		GitTokens:             gitAuthService,
		GitCredentialVerifier: gitAuthService,
		CorsAllowedOrigins:    corsAllowedOrigins,
	})

	s := &http.Server{
		Handler:           r,
		Addr:              "0.0.0.0:" + port,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s, nil
}

func NewRouter(deps ServerDeps) *gin.Engine {
	server := Server{
		repository: *db.New(deps.DB),
		dbConn:     deps.DB,
		auth:       deps.Auth,
		git:        deps.Git,
		gitTokens:  deps.GitTokens,
	}
	strictServer := NewStrictHandler(&server, []StrictMiddlewareFunc{})

	r := gin.Default()
	r.Use(httputil.ServerErrorLogger())
	r.Use(httputil.CorsMiddleware(httputil.CorsConfig{
		AllowedOrigins: deps.CorsAllowedOrigins,
	}))
	r.Use(deps.Auth.Middleware())
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
	gitHandler := gitserver.NewHandler(deps.DB, deps.Git, deps.GitCredentialVerifier)
	r.NoRoute(gitHandler.Handle)

	return r
}

func (s *Server) Healthz(ctx context.Context, request HealthzRequestObject) (HealthzResponseObject, error) {
	return Healthz200JSONResponse{Status: new("ok")}, nil
}

func (s *Server) Version(ctx context.Context, request VersionRequestObject) (VersionResponseObject, error) {
	return Version200JSONResponse{
		Commit:  catena.Commit,
		Version: catena.Version,
	}, nil
}
