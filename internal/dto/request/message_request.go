package request

type MessageRequest struct {
	JobTitle string   `json:"jobTitle" binding:"required"`
	Poster   string   `json:"poster" binding:"required"`
	Skills   []string `json:"skills"` // Optional, can pull from profile
}
