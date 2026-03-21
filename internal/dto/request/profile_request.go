package request

type UpsertProfileRequest struct {
	Skills     []string        `json:"skills"`
	Experience []ExperienceDTO `json:"experience"`
	Projects   []ProjectDTO   `json:"projects"`
	ResumeRaw  string          `json:"resume_raw"`
	Meta       string          `json:"meta"`
}
