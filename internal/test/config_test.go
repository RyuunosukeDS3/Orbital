package test

import (
	"testing"

	"orbital/internal/config"

	"github.com/stretchr/testify/assert"
)

const argocdURL = "argocd.example.com:443"
const dummyToken = "dummy-token"

func TestLoadConfigSuccess(t *testing.T) {
	t.Setenv("ARGOCD_URL", argocdURL)
	t.Setenv("ARGOCD_TOKEN", dummyToken)

	err := config.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config.AppConfig)
	assert.Equal(t, argocdURL, config.AppConfig.ArgoCDURL)
	assert.Equal(t, dummyToken, config.AppConfig.ArgoCDToken)
}

func TestLoadConfigMissingEnvVar(t *testing.T) {
	t.Setenv("ARGOCD_URL", "")
	t.Setenv("ARGOCD_TOKEN", dummyToken)

	err := config.LoadConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required environment variable: ARGOCD_URL")
}

