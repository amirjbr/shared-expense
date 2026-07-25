package config

import (
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func LoadConfig() (*Config, error) {
	var (
		k   = koanf.New(".")
		err error
	)
	conf := &Config{}

	//TODO : read from yaml then read from env
	if err = k.Load(env.Provider("MYVAR_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "MYVAR_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, err
	}

	if err = k.Unmarshal("", conf); err != nil {
		return nil, err
	}

	return conf, nil
}
