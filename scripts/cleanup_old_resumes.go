package main

import (
	"context"
	"fmt"

	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/resume"
	"github.com/yourusername/autoreach-backend/pkg/logger"
	"github.com/yourusername/autoreach-backend/pkg/storage"
)

const retentionDays = 90

func main() {
	cfg := config.LoadConfig()
	config.ConnectDB(cfg)

	if err := storage.InitS3(cfg.AWSRegion, cfg.AWSBucket, cfg.AWSAccessKey, cfg.AWSSecretKey); err != nil {
		logger.Fatal("Failed to initialize S3", err)
	}

	logger.Info(fmt.Sprintf("Running cleanup for resumes older than %d days...", retentionDays))

	resumes, err := resume.GetExpiredResumes(config.DB, retentionDays)
	if err != nil {
		logger.Fatal("Failed to query expired resumes", err)
	}

	logger.Info(fmt.Sprintf("Found %d expired resumes", len(resumes)))

	ctx := context.Background()
	deleted := 0

	for _, r := range resumes {
		// Delete from S3
		if err := storage.Client.Delete(ctx, r.S3Key); err != nil {
			logger.Error(fmt.Sprintf("Failed to delete S3 object %s", r.S3Key), err)
			continue
		}

		// Soft-delete DB record
		if err := config.DB.Delete(&r).Error; err != nil {
			logger.Error(fmt.Sprintf("Failed to delete DB record %s", r.ID), err)
			continue
		}

		deleted++
		logger.Info(fmt.Sprintf("Deleted resume: id=%s, s3_key=%s", r.ID, r.S3Key))
	}

	logger.Info(fmt.Sprintf("Cleanup complete. Deleted %d/%d resumes.", deleted, len(resumes)))
}
