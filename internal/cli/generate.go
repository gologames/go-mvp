package cli

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types,client -package cli -o cli.gen.go -import-mapping "./networkconfig.schemas.yml:github.com/gologames/go-mvp/internal/openapi" ../../api/v1/networkconfig.openapi.yml
