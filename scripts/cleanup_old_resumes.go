package main

import (
	"context"
	"fmt"
	"sync"

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

	if len(resumes) == 0 {
		logger.Info("No resumes to clean up.")
		return
	}

	ctx := context.Background()
	numWorkers := 5
	jobs := make(chan resume.ResumeFile, len(resumes))
	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for r := range jobs {
				logger.Info(fmt.Sprintf("[Worker %d] Deleting resume: id=%s, s3_key=%s", workerID, r.ID, r.S3Key))

				// Delete from S3
				if err := storage.Client.Delete(ctx, r.S3Key); err != nil {
					logger.Error(fmt.Sprintf("[Worker %d] Failed to delete S3 object %s", workerID, r.S3Key), err)
					continue
				}

				// Soft-delete DB record
				if err := config.DB.Delete(&r).Error; err != nil {
					logger.Error(fmt.Sprintf("[Worker %d] Failed to delete DB record %s", workerID, r.ID), err)
					continue
				}
				logger.Info(fmt.Sprintf("[Worker %d] Successfully deleted %s", workerID, r.ID))
			}
		}(w)
	}

	// Send jobs
	for _, r := range resumes {
		jobs <- r
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()

	logger.Info("Cleanup process finished.")
}
