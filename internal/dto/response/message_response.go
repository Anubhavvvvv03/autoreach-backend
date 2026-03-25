package response

type MessageResponse struct {
	MessageID   string `json:"messageId"`
	Text        string `json:"text"`
	GeneratedAt string `json:"generatedAt"`
}
