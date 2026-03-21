package main

import (
	"github.com/yourusername/autoreach-backend/internal/auth"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/profile"
	"github.com/yourusername/autoreach-backend/internal/router"
	"github.com/yourusername/autoreach-backend/pkg/logger"
)


func main() {
    config.AppConfig = config.LoadConfig()
    logger.Info("Starting AutoReach backend on port " + config.AppConfig.Port)

    config.ConnectDB(config.AppConfig)

    // Run migrations
    config.DB.AutoMigrate(&auth.User{}, &profile.Profile{})

    r := router.SetupRouter()
    r.Run(":" + config.AppConfig.Port)
}
