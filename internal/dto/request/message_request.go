package request

type MessageRequest struct {
	Company   string `json:"company" binding:"required"`
	Role      string `json:"role" binding:"required"`
	Recruiter string `json:"recruiter" binding:"required"`
	Tone      string `json:"tone"`
	Context   string `json:"context"`
}
