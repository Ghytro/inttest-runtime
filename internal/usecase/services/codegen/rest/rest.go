package rest

import (
	"inttest-runtime/internal/config"
	"inttest-runtime/internal/useCase/services/configparser"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type RestCodeGenerator struct {
	ApiPrefix string
	Handlers  []RestRoute
}

type RestRoute struct {
	Route     string
	Method    string
	Params    configparser.ParsedRoute
	Behaviour config.RestHandlerBehavior
}

func (r RestRoute) genParamsDefinition() {

}

func (r RestRoute) genMethodLogic() {

}

func (g RestCodeGenerator) Generate() (string, error) {
	type genMethod struct {
		name      string
		restRoute RestRoute
	}

	methods := lo.Map(g.Handlers, func(hand RestRoute, _ int) genMethod {
		return genMethod{
			name:      g.genMethodName(),
			restRoute: hand,
		}
	})

	var code strings.Builder
	code.WriteString(header)
}

func (g RestCodeGenerator) genMethodName() string {
	suffix := uuid.New()
	return "M_" + strings.ReplaceAll(suffix.String(), "-", "")
}
