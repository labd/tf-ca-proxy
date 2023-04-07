package internal

import (
	"encoding/base64"
	"os"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	RegistryName   string   `env:"REGISTRY_NAME,required"`
	RegistryDomain string   `env:"REGISTRY_DOMAIN,required"`
	AuthTokens     []string `env:"AUTH_TOKENS" envSeparator:","`
	SecretKey      string   `env:"SECRET_KEY"`
}

var appConfig AppConfig

func init() {

	if _, err := os.Stat(".env"); err == nil {
		godotenv.Load(".env")
	}

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
