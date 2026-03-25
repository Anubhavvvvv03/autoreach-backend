package request

type UpsertProfileRequest struct {
	FullName    string          `json:"fullName"`
	Title       string          `json:"title"`
	Bio         string          `json:"bio"`
	Location    string          `json:"location"`
	SocialLinks SocialLinksDTO `json:"socialLinks"`
	Skills      []string        `json:"skills"`
	Experience  []ExperienceDTO `json:"experience"`
	Projects    []ProjectDTO    `json:"projects"`
	ResumeRaw   string          `json:"resume_raw"`
	Meta        string          `json:"meta"`
}
