package resume

import (
	"fmt"
	"os"
	"strings"

	pdfLib "github.com/ledongthuc/pdf"
	"github.com/yourusername/autoreach-backend/pkg/logger"
)

// ExtractText reads a PDF file and extracts all text content.
func ExtractText(filePath string) (string, error) {
	f, reader, err := pdfLib.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	var sb strings.Builder
	totalPages := reader.NumPage()

	if totalPages == 0 {
		return "", fmt.Errorf("PDF has no pages")
	}

	logger.Info(fmt.Sprintf("Extracting text from %d pages", totalPages))

	for i := 1; i <= totalPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to extract text from page %d: %s", i, err.Error()))
			continue
		}

		sb.WriteString(text)
		sb.WriteString("\n")
	}

	result := sb.String()
	if strings.TrimSpace(result) == "" {
		return "", fmt.Errorf("no text content found in PDF")
	}

	logger.Info(fmt.Sprintf("Extracted %d characters of text", len(result)))
	return result, nil
}

// DownloadAndExtract downloads a file from S3 to a temp file, extracts text, then cleans up.
func DownloadAndExtract(s3Body []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", "resume-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(s3Body); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	return ExtractText(tmpFile.Name())
}
