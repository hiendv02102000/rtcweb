package repository

import (
	"api/internal/pkg/domain/domain_model/entity"
	"api/pkg/infrastucture/db"
)

type UserBannedRepository struct {
	DB db.Database
}

func (u *UserBannedRepository) FirstUserBanned(condition entity.UserBanned) (entity.UserBanned, error) {
	user := entity.UserBanned{}
	err := u.DB.First(&user, condition)
	return user, err
}
func (u *UserBannedRepository) FindUserBannedList(condition entity.UserBanned) (user []entity.UserBanned, err error) {
	err = u.DB.Find(&user, condition)
	return
}

func (u *UserBannedRepository) CreateBannedUser(user entity.UserBanned) (entity.UserBanned, error) {
	err := u.DB.Create(&user)
	return user, err
}
func (u *UserBannedRepository) DeleteUser(user entity.UserBanned) error {
	err := u.DB.Delete(&user)
	return err
}

func NewUserBannedRepository(db db.Database) *UserBannedRepository {
	return &UserBannedRepository{
		DB: db,
	}
}
