package message

import (
	"bytes"
	"strings"
	"text/template"
	"github.com/yourusername/autoreach-backend/internal/dto/request"

	_ "embed"
)

// Embed the template into the binary
//go:embed templates/outreach.tmpl
var outreachTemplate string

func GenerateMessageService(req request.MessageRequest) (string, error) {
	tmpl, err := template.New("outreach").Parse(outreachTemplate)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"Poster": req.Poster,
		"Job":    req.JobTitle,
		"Skills": strings.Join(req.Skills, ", "),
	}

	var msg bytes.Buffer
	if err := tmpl.Execute(&msg, data); err != nil {
		return "", err
	}

	return msg.String(), nil
}
