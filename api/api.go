package api

import _ "embed"

//go:embed api.v1.openapi.yml
var V1ApiSpec []byte
