package repository

import (
	"api/internal/pkg/domain/domain_model/entity"
	"api/pkg/infrastucture/db"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm/clause"
)

type RoomRepository struct {
	DB db.Database
}

func (u *RoomRepository) FirstRoom(condition entity.Room) (entity.Room, error) {
	Room := entity.Room{}
	err := u.DB.First(&Room, condition)
	return Room, err
}
func (u *RoomRepository) FindRoomList(page int, size int, condition entity.Room) (room []entity.Room, err error) {
	offset := (page - 1) * size
	err = u.DB.DB.Where("title LIKE ? AND is_stream = 1", "%"+condition.Title+"%").Offset(offset).Limit(size).Preload(clause.Associations).Scan(&room).Error
	return
}

func (u *RoomRepository) CreateRoom(room entity.Room) (entity.Room, error) {
	room.ID = uuid.Must(uuid.NewV4(), nil).String()
	err := u.DB.Create(&room)
	return room, err
}
func (u *RoomRepository) UpdateRoom(Room, oldRoom entity.Room) (entity.Room, error) {
	err := u.DB.Update(&entity.Room{}, &oldRoom, &Room)
	return Room, err
}

func NewRoomRepository(db db.Database) *RoomRepository {
	return &RoomRepository{
		DB: db,
	}
}
