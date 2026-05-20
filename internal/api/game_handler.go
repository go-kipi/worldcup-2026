package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kipi/worldcup-2026/internal/service"
)

type GameHandler struct {
	gameService *service.GameService
}

func NewGameHandler(gameService *service.GameService) *GameHandler {
	return &GameHandler{gameService: gameService}
}

func (h *GameHandler) GetAppState(c *gin.Context) {
	userIdVal, exists := c.Get("user_id")
	var userID uint
	if exists {
		userID = userIdVal.(uint)
	}

	state, err := h.gameService.GetAppState(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load app state"})
		return
	}

	c.JSON(http.StatusOK, state)
}

func (h *GameHandler) UpdateMatch(c *gin.Context) {
	if role, _ := c.Get("role"); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admins only"})
		return
	}

	var req struct {
		ID        uint `json:"id"`
		HomeScore *int `json:"home_score"`
		AwayScore *int `json:"away_score"`
		Finished  bool `json:"finished"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.gameService.UpdateMatch(req.ID, req.HomeScore, req.AwayScore, req.Finished); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match updated"})
}

func (h *GameHandler) UpdateKnockoutSlot(c *gin.Context) {
	if role, _ := c.Get("role"); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admins only"})
		return
	}

	var req struct {
		ID    uint                   `json:"id"`
		Patch map[string]interface{} `json:"patch"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.gameService.UpdateKnockoutSlot(req.ID, req.Patch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slot updated"})
}

func (h *GameHandler) UpdateSettings(c *gin.Context) {
	if role, _ := c.Get("role"); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admins only"})
		return
	}

	var req struct {
		LockAt string `json:"lock_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.gameService.UpdateSettings(req.LockAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}
