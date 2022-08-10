package dto

type ChangePassWordRequest struct {
	OldPassword string `json:"old_password" form:"old_password" binding:"password"`
	NewPassword string `json:"new_password" form:"new_password" binding:"password"`
}
