package request

type ProjectDTO struct {
	Name        string `json:"name"`
	TechStack   string `json:"tech_stack"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Meta        string `json:"meta"`
}

type ExperienceDTO struct {
	CompanyName string `json:"company_name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	WorkDone    string `json:"work_done"`
}
