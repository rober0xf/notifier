package dto

type RegisterPayload struct {
	Username string `json:"username" binding:"required,min=4,max=32"`
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
