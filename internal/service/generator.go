package service

import (
	"bytes"
	"strings"
	"text/template"

	_ "embed"

	"github.com/yourusername/autoreach-backend/internal/model"
)

// Embed the template into the binary
//go:embed templates/outreach.tmpl
var outreachTemplate string

func GenerateMessage(req model.MessageRequest) (string, error) {
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
