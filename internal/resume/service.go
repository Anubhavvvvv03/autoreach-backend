package resume

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/pkg/logger"
	"github.com/yourusername/autoreach-backend/pkg/storage"
)

// ProcessResume orchestrates the full resume parsing pipeline.
func ProcessResume(ctx context.Context, userID, resumeID, s3Key string) (*ParsedResume, error) {
	logger.Info(fmt.Sprintf("Starting resume processing for user=%s, resume=%s", userID, resumeID))

	// Download PDF from S3
	logger.Info("Downloading PDF from S3...")
	reader, err := storage.Client.Download(ctx, s3Key)
	if err != nil {
		markFailed(resumeID, "Failed to download from S3: "+err.Error())
		return nil, fmt.Errorf("S3 download failed: %w", err)
	}
	defer reader.Close()

	pdfBytes, err := io.ReadAll(reader)
	if err != nil {
		markFailed(resumeID, "Failed to read S3 response: "+err.Error())
		return nil, fmt.Errorf("failed to read S3 response: %w", err)
	}

	// Extract text from PDF
	logger.Info("Extracting text from PDF...")
	rawText, err := DownloadAndExtract(pdfBytes)
	if err != nil {
		markFailed(resumeID, "Text extraction failed: "+err.Error())
		return nil, fmt.Errorf("text extraction failed: %w", err)
	}

	if rawText == "" {
		markFailed(resumeID, "No text extracted from PDF")
		return nil, fmt.Errorf("no text extracted from PDF")
	}

	// Clean text
	logger.Info("Cleaning extracted text...")
	cleanedText := CleanText(rawText)
	logger.Info(fmt.Sprintf("Cleaned text: %d chars", len(cleanedText)))

	// Call OpenAI
	logger.Info("Calling OpenAI for structured extraction...")
	llmResponse, err := CallOpenAI(ctx, cleanedText)
	if err != nil {
		markFailed(resumeID, "LLM failed: "+err.Error())
		return nil, fmt.Errorf("LLM failed: %w", err)
	}

	// Validate and parse JSON
	logger.Info("Validating LLM response...")
	parsed, err := ValidateParsedResume(llmResponse)
	if err != nil {
		markFailed(resumeID, "Validation failed: "+err.Error())
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Store parsed data and update status
	parsedJSON, _ := json.Marshal(parsed)
	if err := UpdateResumeStatus(config.DB, resumeID, StatusSuccess, string(parsedJSON), ""); err != nil {
		logger.Error("Failed to update resume status", err)
	}

	logger.Info(fmt.Sprintf("Resume processing complete for resume=%s", resumeID))
	return parsed, nil
}

func markFailed(resumeID, reason string) {
	logger.Warn(fmt.Sprintf("Resume %s failed: %s", resumeID, reason))
	UpdateResumeStatus(config.DB, resumeID, StatusFailed, "", reason)
}
