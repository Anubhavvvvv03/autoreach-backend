package response

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"fullName"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
}

type SignupResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
