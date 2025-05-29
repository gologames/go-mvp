package server

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types,spec,chi-server -package server -o server.gen.go -import-mapping "./networkconfig.schemas.yaml:github.com/gologames/go-mvp/internal/openapi" ../../api/v1/networkconfig.openapi.yaml
