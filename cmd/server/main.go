package main

import (
	"github.com/yourusername/autoreach-backend/internal/auth"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/profile"
	"github.com/yourusername/autoreach-backend/internal/resume"
	"github.com/yourusername/autoreach-backend/internal/router"
	"github.com/yourusername/autoreach-backend/pkg/logger"
	"github.com/yourusername/autoreach-backend/pkg/storage"
)


func main() {
    config.AppConfig = config.LoadConfig()
    logger.Info("Starting AutoReach backend on port " + config.AppConfig.Port)

    config.ConnectDB(config.AppConfig)

    // Initialize S3
    if err := storage.InitS3(
        config.AppConfig.AWSRegion,
        config.AppConfig.AWSBucket,
        config.AppConfig.AWSAccessKey,
        config.AppConfig.AWSSecretKey,
    ); err != nil {
        logger.Warn("S3 initialization failed: " + err.Error())
    }

    
    // Run migrations
    config.DB.AutoMigrate(&auth.User{}, &profile.Profile{}, &resume.ResumeFile{})

    r := router.SetupRouter()
    r.Run(":" + config.AppConfig.Port)
}
