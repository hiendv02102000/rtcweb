package dto

// CreateUserRequest struct
type CreateUserRequest struct {
	Password  string `json:"password" form:"password" binding:"required,password"`
	Email     string `json:"email" form:"email" binding:"email"`
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
}

// CreateUserResponse struct
type CreateUserResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}
