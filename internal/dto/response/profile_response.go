package response

import "github.com/yourusername/autoreach-backend/internal/dto/request"

type ProfileResponse struct {
	ID         string                  `json:"id"`
	UserID     string                  `json:"user_id"`
	Skills     []string                `json:"skills"`
	Experience []request.ExperienceDTO `json:"experience"`
	Projects   []request.ProjectDTO    `json:"projects"`
	ResumeRaw  string                  `json:"resume_raw"`
	Meta       string                  `json:"meta"`
}
