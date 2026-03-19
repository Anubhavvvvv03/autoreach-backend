package main

import (
	"github.com/yourusername/autoreach-backend/internal/auth"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/router"
	"github.com/yourusername/autoreach-backend/pkg/logger"
)


func main() {
    cfg := config.LoadConfig()
    logger.Info("Starting AutoReach backend on port " + cfg.Port)

    config.ConnectDB(cfg)

    // Run migrations
    config.DB.AutoMigrate(&auth.User{})

    r := router.SetupRouter()
    r.Run(":" + cfg.Port)
}
