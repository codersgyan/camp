package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	HttpPort string
}

func SetupENV() (cfg *AppConfig, err error) {

	godotenv.Load()

	httpPort := os.Getenv("HTTP_PORT")

	if len(httpPort) < 1 {
		return &AppConfig{}, errors.New("env file is not configured properly")
	}

	return &AppConfig{
		HttpPort: httpPort,
	}, nil
}
