package config

import (
	"encoding/json"
	"os"
)

func FromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := new(Config)
	if err := json.NewDecoder(f).Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}
