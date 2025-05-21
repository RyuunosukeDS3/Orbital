package test

import (
	"testing"

	"orbital/internal/config"

	"github.com/stretchr/testify/assert"
)

const argocdURL = "argocd.example.com:443"
const dummyToken = "dummy-token"
const dummyVrisingIp = "192.168.1.100"
const dummyMac = "12:AB:34:CD:56:EF"

func TestLoadConfigSuccess(t *testing.T) {
	t.Setenv("ARGOCD_URL", argocdURL)
	t.Setenv("ARGOCD_TOKEN", dummyToken)
	t.Setenv("VRISING_IP", dummyVrisingIp)
	t.Setenv("VRISING_PORT", "12345")
	t.Setenv("MAC_ADDRESS", dummyMac)

	err := config.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config.AppConfig)
	assert.Equal(t, argocdURL, config.AppConfig.ArgoCDURL)
	assert.Equal(t, dummyToken, config.AppConfig.ArgoCDToken)
	assert.Equal(t, dummyVrisingIp, config.AppConfig.VrisingIp)
	assert.Equal(t, "12345", config.AppConfig.VrisingPort)
	assert.Equal(t, dummyMac, config.AppConfig.MacAddress)
}

func TestLoadConfigMissingVrisingIp(t *testing.T) {
	t.Setenv("ARGOCD_URL", argocdURL)
	t.Setenv("ARGOCD_TOKEN", dummyToken)
	t.Setenv("VRISING_IP", "")
	t.Setenv("VRISING_PORT", "12345")
	t.Setenv("MAC_ADDRESS", dummyMac)

	err := config.LoadConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required environment variable: VRISING_IP")
}

func TestLoadConfigMissingVrisingPort(t *testing.T) {
	t.Setenv("ARGOCD_URL", argocdURL)
	t.Setenv("ARGOCD_TOKEN", dummyToken)
	t.Setenv("VRISING_IP", dummyVrisingIp)
	t.Setenv("VRISING_PORT", "")
	t.Setenv("MAC_ADDRESS", dummyMac)

	err := config.LoadConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required environment variable: VRISING_PORT")
}


func TestLoadConfigMissingEnvVar(t *testing.T) {
	t.Setenv("ARGOCD_URL", "")
	t.Setenv("ARGOCD_TOKEN", dummyToken)

	err := config.LoadConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required environment variable: ARGOCD_URL")
}

