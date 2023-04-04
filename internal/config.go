package internal

import (
	"github.com/caarlos0/env/v7"
)

type AppConfig struct {
	RegistryName   string `env:"REGISTRY_NAME"`
	RegistryDomain string `env:"REGISTRY_DOMAIN"`
}

var appConfig AppConfig

func init() {
	if err := env.Parse(&appConfig); err != nil {
		panic(err)
	}
}
