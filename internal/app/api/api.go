package api

import (
	"context"
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
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

//go:generate go tool oapi-codegen -config cfg.yaml ../../../api/api.v1.openapi.yaml

type Server struct {
	repository db.Queries
	dbConn     *pgxpool.Pool
	auth       auth.AuthService
	git        gitstore.Store
	gitAuth    gitauth.Service
}

func NewServer(
	conn *pgxpool.Pool,
	authService auth.AuthService,
	gitService gitstore.Store,
	corsAllowedOrigins []string,
) (*http.Server, error) {
	gitAuthService := gitauth.NewService(conn)
	server := Server{
		repository: *db.New(conn),
		dbConn:     conn,
		auth:       authService,
		git:        gitService,
		gitAuth:    gitAuthService,
	}
	strictServer := NewStrictHandler(&server, []StrictMiddlewareFunc{})

	r := gin.Default()
	r.Use(httputil.CorsMiddleware(httputil.CorsConfig{
		AllowedOrigins: corsAllowedOrigins,
	}))
	r.Use(authService.Middleware())
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
	gitHandler := gitserver.NewHandler(conn, gitService, gitAuthService)
	r.NoRoute(gitHandler.Handle)

	s := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}

	return s, nil
}

func (s *Server) Healthz(ctx context.Context, request HealthzRequestObject) (HealthzResponseObject, error) {
	ok := "ok"
	return Healthz200JSONResponse{Status: &ok}, nil
}
