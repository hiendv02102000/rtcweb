package dto

type GetRoomInfoRequest struct {
	IdRoom string `json:"id_room" form:"id_room" binding:"required"`
}
