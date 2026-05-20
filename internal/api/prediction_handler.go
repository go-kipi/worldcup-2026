package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kipi/worldcup-2026/internal/service"
)

type PredictionHandler struct {
	predictionService *service.PredictionService
}

func NewPredictionHandler(predictionService *service.PredictionService) *PredictionHandler {
	return &PredictionHandler{predictionService: predictionService}
}

func (h *PredictionHandler) GetPredictions(c *gin.Context) {
	userIdVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIdVal.(uint)
	matches, knockouts, err := h.predictionService.GetUserPredictions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load predictions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"match_predictions":    matches,
		"knockout_predictions": knockouts,
	})
}

type savePredictionsRequest struct {
	MatchPredictions    []service.MatchPredictionInput    `json:"match_predictions"`
	KnockoutPredictions []service.KnockoutPredictionInput `json:"knockout_predictions"`
}

func (h *PredictionHandler) SavePredictions(c *gin.Context) {
	userIdVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIdVal.(uint)

	var req savePredictionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.predictionService.SavePredictions(userID, req.MatchPredictions, req.KnockoutPredictions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Predictions saved successfully"})
}
