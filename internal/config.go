package internal

import (
	"github.com/caarlos0/env/v7"
)

type AppConfig struct {
	RegistryName   string   `env:"REGISTRY_NAME"`
	RegistryDomain string   `env:"REGISTRY_DOMAIN"`
	AuthTokens     []string `env:"AUTH_TOKENS" envSeparator:","`
}

var appConfig AppConfig

func init() {
	if err := env.Parse(&appConfig); err != nil {
		panic(err)
	}
}
