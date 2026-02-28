package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=256"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type LoginResponse struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
	Email string `json:"email"`
}
