package configs

import _ "embed"

var (
	//go:embed init-route-base.yaml
	InitRouteConfig []byte
)
