package auth

type Session struct {
	User SessionUser `json:"user"`
}

type SessionUser struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Image       string `json:"image"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}
