package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kipi/worldcup-2026/internal/api"
	"github.com/go-kipi/worldcup-2026/internal/config"
	"github.com/go-kipi/worldcup-2026/internal/db"
	"github.com/go-kipi/worldcup-2026/internal/routes"
	"github.com/go-kipi/worldcup-2026/internal/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newServer(lc fx.Lifecycle, cfg *config.Config, router *gin.Engine) *http.Server {
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Starting HTTP server on port %s", cfg.Port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Stopping HTTP server")
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func newLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return zap.NewDevelopment()
	}
	return logger, nil
}

func main() {
	app := fx.New(
		fx.Provide(
			config.LoadConfig,
			db.NewDatabase,
			db.NewMongoDatabase,
			newLogger,

			service.NewAuthService,
			service.NewGameService,
			service.NewPredictionService,
			service.NewLeaderboardService,
			service.NewEmailService,

			api.NewAuthHandler,
			api.NewGameHandler,
			api.NewPredictionHandler,
			api.NewLeaderboardHandler,

			routes.NewRouter,
			newServer,
		),
		fx.Invoke(func(*http.Server) {}),
	)

	app.Run()
}
