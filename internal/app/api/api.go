package api

import (
	"context"
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/floffah/catena/api"
	"github.com/gin-gonic/gin"
)

//go:generate go tool oapi-codegen -config cfg.yaml ../../../api/api.v1.openapi.yml

type Server struct {
}

func NewServer() *http.Server {
	server := Server{}
	strictServer := NewStrictHandler(&server, []StrictMiddlewareFunc{})

	r := gin.Default()
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

	return s
}

func (s *Server) Healthz(ctx context.Context, request HealthzRequestObject) (HealthzResponseObject, error) {
	ok := "ok"
	return Healthz200JSONResponse{Status: &ok}, nil
}
