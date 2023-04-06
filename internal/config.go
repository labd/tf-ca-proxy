package internal

import (
	"encoding/base64"

	"github.com/caarlos0/env/v7"
)

type AppConfig struct {
	RegistryName   string   `env:"REGISTRY_NAME,required"`
	RegistryDomain string   `env:"REGISTRY_DOMAIN,required"`
	AuthTokens     []string `env:"AUTH_TOKENS" envSeparator:","`
	SecretKey      string   `env:"SECRET_KEY"`
}

var appConfig AppConfig

func init() {
	if err := env.Parse(&appConfig); err != nil {
		panic(err)
	}

	if appConfig.SecretKey == "" {
		val, err := generateRandomBytes(32)
		if err != nil {
			panic(err)
		}
		appConfig.SecretKey = base64.URLEncoding.EncodeToString(val)
	}
}
