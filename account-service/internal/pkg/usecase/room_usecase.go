package usecase

import (
	"api/internal/pkg/domain/domain_model/dto"
	"api/internal/pkg/domain/domain_model/entity"
	"api/internal/pkg/domain/repository"
	"context"
	"errors"

	"api/pkg/infrastucture/db"
)

type RoomUsecase interface {
	GetRoomList(req dto.GetRoomListRequest) (dto.GetRoomListResponse, error)
	GetRoomInfo(roomId string) (entity.Room, error)
	StartRoom(entity.Users, dto.StartRoomRequest) (entity.Room, error)
	EndRoom(entity.Users) (entity.Room, error)
}

type roomUsecase struct {
	roomRepo repository.RoomRepository
}

func (u *roomUsecase) GetRoomList(req dto.GetRoomListRequest) (dto.GetRoomListResponse, error) {
	rooms, err := u.roomRepo.FindRoomList(req.Page, req.Size, entity.Room{
		Title: req.Title,
	})
	return dto.GetRoomListResponse{
		Total:    len(rooms),
		RoomList: rooms,
	}, err
}
func (u *roomUsecase) GetRoomInfo(roomId string) (entity.Room, error) {
	return u.roomRepo.FirstRoom(entity.Room{ID: roomId})
}
func (u *roomUsecase) StartRoom(user entity.Users, req dto.StartRoomRequest) (entity.Room, error) {
	room, err := u.roomRepo.FirstRoom(entity.Room{StreamerID: user.ID})
	if err != nil {
		return entity.Room{}, err
	}
	if room.IsStream {
		return entity.Room{}, errors.New("room is live streaming")
	}
	if len(room.ID) == 0 {
		room, err = u.roomRepo.CreateRoom(entity.Room{
			Title:      req.Title,
			StreamerID: user.ID,
			IsStream:   true,
		})
		if err != nil {
			return entity.Room{}, err
		}
	} else {
		newRoom := room
		newRoom.Title = req.Title
		newRoom.IsStream = true
		room, err = u.roomRepo.UpdateRoom(newRoom, entity.Room{
			ID: room.ID,
		})
		if err != nil {
			return entity.Room{}, err
		}

	}
	db.RedisPool.Set(context.Background(), "room-user-"+room.ID, nil, 0)
	db.RedisPool.Set(context.Background(), "room-chat-"+room.ID, nil, 0)
	return room, nil
}
func (u *roomUsecase) EndRoom(user entity.Users) (entity.Room, error) {
	room, err := u.roomRepo.FirstRoom(entity.Room{StreamerID: user.ID})
	if err != nil {
		return entity.Room{}, err
	}
	if room.ID == "" || !room.IsStream {
		return entity.Room{}, errors.New("room is stopping")
	}

	newRoom := room
	newRoom.IsStream = false
	room, err = u.roomRepo.UpdateRoom(newRoom, entity.Room{
		ID: room.ID,
	})
	if err != nil {
		return entity.Room{}, err
	}

	db.RedisPool.Del(context.Background(), "room-user-"+room.ID)
	db.RedisPool.Del(context.Background(), "room-chat-"+room.ID)
	return room, nil
}

func NewRoomUsecase(db db.Database) RoomUsecase {
	repo := repository.NewRoomRepository(db)
	return &roomUsecase{
		roomRepo: *repo,
	}
}
