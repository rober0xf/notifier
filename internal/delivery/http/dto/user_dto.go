package dto

import "github.com/rober0xf/notifier/internal/domain/entity"

type UserPayload struct {
	ID       int    `json:"id"       example:"1"`
	Username string `json:"username" example:"rober"`
	Email    string `json:"email"    example:"rober@example.com"`
	Active   bool   `json:"active"   example:"false"`
}

type MeResponse struct {
	UserID int    `json:"user_id" example:"1"`
	Email  string `json:"email"   example:"rober@example.com"`
	Role   string `json:"role"    example:"user"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" example:"rober"`
	Email    *string `json:"email,omitempty" example:"rober@example.com"`
	Password *string `json:"password,omitempty" example:"newpassword123#!"`
}

func ToUserResponse(user *entity.User) UserPayload {
	return UserPayload{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Active:   user.IsActive,
	}
}
