package dto

import "time"

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
