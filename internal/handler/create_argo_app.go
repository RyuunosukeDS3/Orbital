package handler

import (
	"net/http"
	"orbital/internal/argocd"
	"orbital/internal/config"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	AppName string `json:"appName"`
}

func UpsertArgoApp(c *gin.Context) {
	var requestBody RequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call ArgoCD Upsert logic
	if err := argocd.UpsertApplication(requestBody.AppName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create ArgoCD application",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "ArgoCD application created or updated successfully",
		"appName":     requestBody.AppName,
		"argocdUrl":   config.AppConfig.ArgoCDURL,
	})
}
