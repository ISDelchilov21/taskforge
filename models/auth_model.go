package models

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
}
