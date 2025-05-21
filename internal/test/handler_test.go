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

const route string = "/upsert" 
const contentType string = "Content-Type"
const appJson string = "application/json"

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST(route, handler.IncreaseArgoAppReplicaCount)
	return r
}

func TestIncreaseArgoAppReplicaCountSuccess(t *testing.T) {
	config.AppConfig = &config.Config{ArgoCDURL: "argocd.example.com:443"}

	// override the dependency with mock
	argocd.IncreaseReplicaCount = func(appName string) error {
		assert.Equal(t, "test-app", appName)
		return nil
	}

	router := setupRouter()
	
	payload := map[string]string{"appName": "test-app"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", route, bytes.NewBuffer(body))
	req.Header.Set(contentType, appJson)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ArgoCD application created or updated successfully")
}


func TestIncreaseArgoAppReplicaCountBadRequest(t *testing.T) {
	req, _ := http.NewRequest("POST", route, bytes.NewBuffer([]byte("not json")))
	req.Header.Set(contentType, appJson)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestIncreaseArgoAppReplicaCountFailure(t *testing.T) {
	config.AppConfig = &config.Config{ArgoCDURL: "https://argocd.example.com"}

	argocd.IncreaseReplicaCount = func(appName string) error {
		return fmt.Errorf("simulated failure")
	}

	payload := map[string]string{"appName": "bad-app"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", route, bytes.NewBuffer(body))
	req.Header.Set(contentType, appJson)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create ArgoCD application")
}

