package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	ArgoCDURL   string
	ArgoCDToken string
}

var AppConfig *Config

func LoadConfig() error {
	// Required environment variables
	required := []string{
		"ARGOCD_URL",
		"ARGOCD_TOKEN",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
	}

	AppConfig = &Config{
		ArgoCDURL:   os.Getenv("ARGOCD_URL"),
		ArgoCDToken: os.Getenv("ARGOCD_TOKEN"),
	}

	log.Println("Configuration loaded successfully.")
	return nil
}

