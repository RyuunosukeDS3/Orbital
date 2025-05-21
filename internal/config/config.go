package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	ArgoCDURL   string
	ArgoCDToken string
	VrisingIp   string
	VrisingPort string
	MacAddress  string
}

var AppConfig *Config

func LoadConfig() error {
	required := []string{
		"ARGOCD_URL",
		"ARGOCD_TOKEN",
		"VRISING_IP",
		"VRISING_PORT",
		"MAC_ADDRESS",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
	}

	AppConfig = &Config{
		ArgoCDURL:   os.Getenv("ARGOCD_URL"),
		ArgoCDToken: os.Getenv("ARGOCD_TOKEN"),
		VrisingIp:   os.Getenv("VRISING_IP"),
		VrisingPort: os.Getenv("VRISING_PORT"),
		MacAddress:  os.Getenv("MAC_ADDRESS"),
	}

	log.Println("Configuration loaded successfully.")
	return nil
}
