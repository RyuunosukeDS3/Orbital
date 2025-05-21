package handler

import (
	"net/http"
	"orbital/internal/argocd"
	"orbital/internal/config"
	"orbital/internal/monitoring"
	wol "orbital/internal/wake_on_lan"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	AppName string `json:"appName"`
}

var Sleep = time.Sleep

func IncreaseArgoAppReplicaCount(c *gin.Context) {
	var requestBody RequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := wol.WakeOnLan(config.AppConfig.MacAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send Wake-on-LAN packet",
			"details": err.Error(),
		})
		return
	}

	if err := argocd.SetReplicaCount(requestBody.AppName, "1"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create ArgoCD application",
			"details": err.Error(),
		})
		return
	}

	go func() {
		Sleep(5 * time.Minute)
		monitoring.PlayersOffline(config.AppConfig.VrisingIp, config.AppConfig.VrisingPort, requestBody.AppName)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":   "ArgoCD application created or updated successfully and monitoring started",
		"appName":   requestBody.AppName,
		"argocdUrl": config.AppConfig.ArgoCDURL,
	})
}
