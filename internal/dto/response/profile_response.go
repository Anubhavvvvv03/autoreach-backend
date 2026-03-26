package response

import "github.com/yourusername/autoreach-backend/internal/dto/request"

type ProfileResponse struct {
	ID          string                  `json:"id"`
	UserID      string                  `json:"user_id"`
	FullName    string                  `json:"fullName"`
	Title       string                  `json:"title"`
	Bio         string                  `json:"bio"`
	Location    string                  `json:"location"`
	SocialLinks request.SocialLinksDTO  `json:"socialLinks"`
	Skills      []string                `json:"skills"`
	Experience  []request.ExperienceDTO `json:"experience"`
	Projects    []request.ProjectDTO    `json:"projects"`
	ResumeRaw   string                  `json:"resume_raw"`
	Meta        string                  `json:"meta"`
}
