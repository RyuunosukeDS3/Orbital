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
	wol "orbital/internal/wake_on_lan"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIncreaseArgoAppReplicaCountSuccess(t *testing.T) {
	config.AppConfig = &config.Config{
		ArgoCDURL:  "argocd.example.com:443",
		VrisingIp:  "127.0.0.1",
		VrisingPort: "9876",
		MacAddress: "00:11:22:33:44:55",
	}

	origSleep := handler.Sleep
	handler.Sleep = func(d time.Duration) {}
	defer func() { handler.Sleep = origSleep }()

	origWake := wol.WakeOnLan
	wol.WakeOnLan = func(mac string) error {
		assert.Equal(t, "00:11:22:33:44:55", mac)
		return nil
	}
	defer func() { wol.WakeOnLan = origWake }()

	origSetReplica := argocd.SetReplicaCount
	argocd.SetReplicaCount = func(appName, count string) error {
		assert.Equal(t, "test-app", appName)
		assert.Equal(t, "1", count)
		return nil
	}
	defer func() { argocd.SetReplicaCount = origSetReplica }()

	router := gin.Default()
	router.POST("/upsert", handler.IncreaseArgoAppReplicaCount)

	payload := map[string]string{"appName": "test-app"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/upsert", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ArgoCD application created or updated successfully")
}

func TestIncreaseArgoAppReplicaCountBadRequest(t *testing.T) {
	router := gin.Default()
	router.POST("/upsert", handler.IncreaseArgoAppReplicaCount)

	req, _ := http.NewRequest("POST", "/upsert", bytes.NewBuffer([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestIncreaseArgoAppReplicaCountWakeOnLanFailure(t *testing.T) {
	config.AppConfig = &config.Config{
		ArgoCDURL:  "https://argocd.example.com",
		VrisingIp:  "127.0.0.1",
		VrisingPort: "9876",
		MacAddress: "00:11:22:33:44:55",
	}

	origWake := wol.WakeOnLan
	wol.WakeOnLan = func(mac string) error {
		return fmt.Errorf("simulated WOL failure")
	}
	defer func() { wol.WakeOnLan = origWake }()

	payload := map[string]string{"appName": "test-app"}
	body, _ := json.Marshal(payload)

	router := gin.Default()
	router.POST("/upsert", handler.IncreaseArgoAppReplicaCount)

	req, _ := http.NewRequest("POST", "/upsert", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to send Wake-on-LAN packet")
}

func TestIncreaseArgoAppReplicaCountSetReplicaCountFailure(t *testing.T) {
	config.AppConfig = &config.Config{
		ArgoCDURL:  "https://argocd.example.com",
		VrisingIp:  "127.0.0.1",
		VrisingPort: "9876",
		MacAddress: "00:11:22:33:44:55",
	}

	origWake := wol.WakeOnLan
	wol.WakeOnLan = func(mac string) error {
		return nil
	}
	defer func() { wol.WakeOnLan = origWake }()

	origSetReplica := argocd.SetReplicaCount
	argocd.SetReplicaCount = func(appName, count string) error {
		return fmt.Errorf("simulated failure")
	}
	defer func() { argocd.SetReplicaCount = origSetReplica }()

	payload := map[string]string{"appName": "bad-app"}
	body, _ := json.Marshal(payload)

	router := gin.Default()
	router.POST("/upsert", handler.IncreaseArgoAppReplicaCount)

	req, _ := http.NewRequest("POST", "/upsert", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create ArgoCD application")
}
