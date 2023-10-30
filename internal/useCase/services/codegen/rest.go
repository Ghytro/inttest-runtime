package codegen

import (
	"inttest-runtime/internal/config"
	"inttest-runtime/internal/useCase/services/configparser"
	"strings"
)

type RestCodeGenerator struct {
	ApiPrefix string
	Handlers  []RestRoute
}

type RestRoute struct {
	Route     string
	Params    configparser.ParsedRoute
	Behaviour config.RestHandlerBehavior
}

func (g RestCodeGenerator) Generate() (string, error) {
	var code strings.Builder
}
