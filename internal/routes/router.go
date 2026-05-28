package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kipi/worldcup-2026/internal/api"
	"github.com/go-kipi/worldcup-2026/internal/config"
	"github.com/go-kipi/worldcup-2026/internal/pkg/jwtutil"
)

// NewRouter sets up the Gin router and registers all routes.
func NewRouter(
	cfg *config.Config,
	authHandler *api.AuthHandler,
	gameHandler *api.GameHandler,
	predictionHandler *api.PredictionHandler,
	leaderboardHandler *api.LeaderboardHandler,
) *gin.Engine {
	r := gin.Default()

	// CORS middleware... (existing code)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		auth := apiGroup.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/verify", authHandler.Verify)
		}

		apiGroup.GET("/leaderboard", leaderboardHandler.GetLeaderboard)

		protectedGroup := apiGroup.Group("/")
		protectedGroup.Use(AuthMiddleware(cfg))
		{
			protectedGroup.GET("/app-state", gameHandler.GetAppState)

			// Admin routes
			protectedGroup.PUT("/admin/match", gameHandler.UpdateMatch)
			protectedGroup.PUT("/admin/slot", gameHandler.UpdateKnockoutSlot)
			protectedGroup.PUT("/admin/settings", gameHandler.UpdateSettings)

			// predictions endpoints remain the same, though the payload shape will be refactored
			protectedGroup.GET("/predictions", predictionHandler.GetPredictions)
			protectedGroup.POST("/predictions", predictionHandler.SavePredictions)
		}
	}

	// Serve static assets
	r.Static("/assets", "./dist/assets")

	// Serve favicon
	r.StaticFile("/favicon.ico", "./dist/favicon.ico")

	// Serve root index.html
	r.StaticFile("/", "./dist/index.html")

	// SPA fallback
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	return r
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtutil.ValidateToken(tokenString, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		c.Next()
	}
}
