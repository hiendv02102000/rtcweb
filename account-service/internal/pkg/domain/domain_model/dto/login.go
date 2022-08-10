package dto

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"email"`
	Password string `json:"password" form:"password" binding:"password"`
}
