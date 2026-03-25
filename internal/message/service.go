package message

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/yourusername/autoreach-backend/internal/config"
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

	var skills []string
	if prof != nil {
		skills = prof.Skills
	}

	data := map[string]interface{}{
		"Recruiter": req.Recruiter,
		"Role":      req.Role,
		"Company":   req.Company,
		"Tone":      req.Tone,
		"Context":   req.Context,
		"Skills":    strings.Join(skills, ", "),
	}

	var msgBuf bytes.Buffer
	if err := tmpl.Execute(&msgBuf, data); err != nil {
		return "", err
	}

	generatedText := msgBuf.String()

	// Persist to history
	msgRecord := Message{
		UserID:    userID,
		Company:   req.Company,
		Role:      req.Role,
		Recruiter: req.Recruiter,
		Tone:      req.Tone,
		Context:   req.Context,
		Text:      generatedText,
		Status:    "GENERATED",
	}

	config.DB.Create(&msgRecord)

	return generatedText, nil
}

func GetMessageHistory(userID string) ([]Message, error) {
	var messages []Message
	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
	
