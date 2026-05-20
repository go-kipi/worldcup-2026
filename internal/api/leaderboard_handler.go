package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kipi/worldcup-2026/internal/service"
)

type LeaderboardHandler struct {
	leaderboardService *service.LeaderboardService
}

func NewLeaderboardHandler(leaderboardService *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{leaderboardService: leaderboardService}
}

func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	users, err := h.leaderboardService.GetLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load leaderboard"})
		return
	}

	c.JSON(http.StatusOK, users)
}
