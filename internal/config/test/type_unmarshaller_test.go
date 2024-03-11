package test

import (
	"encoding/json"
	"fmt"
	"inttest-runtime/internal/config"
	"os"
	"testing"
)

func TestUnmarshalConfig(t *testing.T) {
	conf, err := parseConfig()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*conf)
}

func TestParseRestRoute(t *testing.T) {
	conf, err := parseConfig()
	if err != nil {
		t.Fatal(err)
	}
	routeParams := conf.RpcServices[0].RpcServiceUnion.RestService.Routes[0].Route.Params()
	fmt.Println(routeParams)
}

func parseConfig() (*config.Config, error) {
	payload, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var conf config.Config
	if err := json.Unmarshal(payload, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
