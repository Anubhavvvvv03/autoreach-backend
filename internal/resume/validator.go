package resume

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"
)

// ParsedResume represents the validated, structured output from the LLM.
type ParsedResume struct {
	Name       string       `json:"name"`
	Email      string       `json:"email"`
	Skills     []string     `json:"skills"`
	Experience []ExpParsed  `json:"experience"`
	Projects   []ProjParsed `json:"projects"`
	Education  []EduParsed  `json:"education"`
}

type ExpParsed struct {
	CompanyName string `json:"company_name"`
	Role        string `json:"role"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	WorkDone    string `json:"work_done"`
}

type ProjParsed struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TechStack   string `json:"tech_stack"`
	URL         string `json:"url"`
}

type EduParsed struct {
	Degree      string `json:"degree"`
	Institution string `json:"institution"`
	Year        string `json:"year"`
}

// ValidateParsedResume takes raw JSON from the LLM and validates + sanitizes it.
func ValidateParsedResume(rawJSON string) (*ParsedResume, error) {
	// Strip markdown code fences if present
	rawJSON = strings.TrimSpace(rawJSON)
	rawJSON = strings.TrimPrefix(rawJSON, "```json")
	rawJSON = strings.TrimPrefix(rawJSON, "```")
	rawJSON = strings.TrimSuffix(rawJSON, "```")
	rawJSON = strings.TrimSpace(rawJSON)

	var parsed ParsedResume
	if err := json.Unmarshal([]byte(rawJSON), &parsed); err != nil {
		return nil, fmt.Errorf("invalid JSON from LLM: %w", err)
	}

	// Trim all strings
	parsed.Name = strings.TrimSpace(parsed.Name)
	parsed.Email = strings.TrimSpace(parsed.Email)

	// Validate email format (if provided)
	if parsed.Email != "" {
		if _, err := mail.ParseAddress(parsed.Email); err != nil {
			parsed.Email = "" // Clear invalid email silently
		}
	}

	// Clean skills
	var cleanSkills []string
	for _, s := range parsed.Skills {
		s = strings.TrimSpace(s)
		if s != "" {
			cleanSkills = append(cleanSkills, s)
		}
	}
	parsed.Skills = cleanSkills

	// Clean experience
	var cleanExp []ExpParsed
	for _, e := range parsed.Experience {
		e.CompanyName = strings.TrimSpace(e.CompanyName)
		e.Role = strings.TrimSpace(e.Role)
		e.WorkDone = strings.TrimSpace(e.WorkDone)
		e.StartDate = strings.TrimSpace(e.StartDate)
		e.EndDate = strings.TrimSpace(e.EndDate)
		if e.CompanyName != "" || e.Role != "" {
			cleanExp = append(cleanExp, e)
		}
	}
	parsed.Experience = cleanExp

	// Clean projects
	var cleanProj []ProjParsed
	for _, p := range parsed.Projects {
		p.Name = strings.TrimSpace(p.Name)
		p.Description = strings.TrimSpace(p.Description)
		p.TechStack = strings.TrimSpace(p.TechStack)
		p.URL = strings.TrimSpace(p.URL)
		if p.Name != "" {
			cleanProj = append(cleanProj, p)
		}
	}
	parsed.Projects = cleanProj

	// Clean education
	var cleanEdu []EduParsed
	for _, e := range parsed.Education {
		e.Degree = strings.TrimSpace(e.Degree)
		e.Institution = strings.TrimSpace(e.Institution)
		e.Year = strings.TrimSpace(e.Year)
		if e.Degree != "" || e.Institution != "" {
			cleanEdu = append(cleanEdu, e)
		}
	}
	parsed.Education = cleanEdu

	return &parsed, nil
}
