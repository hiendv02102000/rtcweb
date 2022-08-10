package dto

import "api/internal/pkg/domain/domain_model/entity"

type GetRoomListRequest struct {
	Title string `json:"title" form:"title" binding:"omitemty"`
	Page  int    `json:"page" form:"page" binding:"required,min=1"`
	Size  int    `json:"size" form:"size" binding:"required,min=1"`
}
type GetRoomListResponse struct {
	Total    int           `json:"total"`
	RoomList []entity.Room `json:"room_list"`
}
