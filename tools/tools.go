package main

import (
	_ "github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ../openapi.json
