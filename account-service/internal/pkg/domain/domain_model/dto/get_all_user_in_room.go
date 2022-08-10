package dto

import "api/internal/pkg/domain/domain_model/entity"

type GetAllUserInRoomResponse struct {
	ListUser []entity.Users `json:"list_user"`
}
