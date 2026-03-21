package message

import (
	"bytes"
	"strings"
	"text/template"
	"github.com/yourusername/autoreach-backend/internal/dto/request"
	"github.com/yourusername/autoreach-backend/internal/profile"

	_ "embed"
)

// Embed the template into the binary
//go:embed templates/outreach.tmpl
var outreachTemplate string

func GenerateMessageService(req request.MessageRequest, userID string) (string, error) {
	tmpl, err := template.New("outreach").Parse(outreachTemplate)
	if err != nil {
		return "", err
	}

	// Fetch profile for additional context if available
	prof, _ := profile.GetProfileByUserID(userID)

	poster := req.Poster
	job := req.JobTitle
	var skills []string

	if len(req.Skills) > 0 {
		skills = req.Skills
	} else if prof != nil {
		skills = prof.Skills
	}

	data := map[string]interface{}{
		"Poster": poster,
		"Job":    job,
		"Skills": strings.Join(skills, ", "),
	}

	var msg bytes.Buffer
	if err := tmpl.Execute(&msg, data); err != nil {
		return "", err
	}

	return msg.String(), nil
}
