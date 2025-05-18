package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"orbital/internal/argocd"
	"orbital/internal/config"
	"orbital/internal/handler"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mock implementation of argocd.UpsertApplication
var mockUpsertApplication func(appName string) error

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/upsert", handler.UpsertArgoApp)
	return r
}

func TestUpsertArgoAppSuccess(t *testing.T) {
	config.AppConfig = &config.Config{ArgoCDURL: "argocd.example.com:443"}

	// override the dependency with mock
	argocd.UpsertApplication = func(appName string) error {
		assert.Equal(t, "test-app", appName)
		return nil
	}

	router := setupRouter()
	
	payload := map[string]string{"appName": "test-app"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/upsert", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ArgoCD application created or updated successfully")
}


func TestUpsertArgoAppBadRequest(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upsert", bytes.NewBuffer([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestUpsertArgoAppFailure(t *testing.T) {
	config.AppConfig = &config.Config{ArgoCDURL: "https://argocd.example.com"}

	argocd.UpsertApplication = func(appName string) error {
		return fmt.Errorf("simulated failure")
	}

	payload := map[string]string{"appName": "bad-app"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/upsert", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create ArgoCD application")
}

