package resume

import (
	"github.com/yourusername/autoreach-backend/internal/dto/request"
)

// ParseResume is a placeholder for actual resume parsing logic.
// In the future, this will use an AI service or a parsing library.
func ParseResume(rawData string) (*request.UpsertProfileRequest, error) {
	// Mock parsing logic
	return &request.UpsertProfileRequest{
		ResumeRaw: rawData,
		Skills:    []string{"Go", "PostgreSQL", "Gin"},
		Experience: []request.ExperienceDTO{
			{
				CompanyName: "Parsed Tech",
				WorkDone:    "Developed scalable backends",
			},
		},
		Projects: []request.ProjectDTO{
			{
				Name:      "Parsed Project",
				TechStack: "Go, GORM",
			},
		},
		Meta: "Auto-parsed from resume",
	}, nil
}
