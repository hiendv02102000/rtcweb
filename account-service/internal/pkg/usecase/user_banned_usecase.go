package usecase

import (
	"api/internal/pkg/domain/domain_model/entity"
	"api/internal/pkg/domain/repository"

	"api/pkg/infrastucture/db"
)

type UserBannedUsecase interface {
	BanUser(streamer entity.Users, userId int) error
	UnBanUser(streamer entity.Users, userId int) error
	GetUserBannedList(streamer entity.Users) ([]entity.Users, error)
}

type userBannedUsecase struct {
	userRepo repository.UserBannedRepository
}

func (u *userBannedUsecase) BanUser(streamer entity.Users, userId int) error {
	userBanned, err := u.userRepo.FirstUserBanned(entity.UserBanned{
		UserID:     userId,
		StreamerID: streamer.ID,
	})
	if err != nil {
		return err
	}
	if userBanned.ID != 0 {
		return nil
	}
	_, err = u.userRepo.CreateBannedUser(entity.UserBanned{
		UserID:     userId,
		StreamerID: streamer.ID,
	})
	return err
}
func (u *userBannedUsecase) UnBanUser(streamer entity.Users, userId int) error {
	userBanned, err := u.userRepo.FirstUserBanned(entity.UserBanned{
		UserID:     userId,
		StreamerID: streamer.ID,
	})
	if err != nil {
		return err
	}
	if userBanned.ID == 0 {
		return nil
	}
	err = u.userRepo.DeleteUser(entity.UserBanned{
		UserID:     userId,
		StreamerID: streamer.ID,
	})
	return err
}
func (u *userBannedUsecase) GetUserBannedList(streamer entity.Users) ([]entity.Users, error) {
	userBanneds, err := u.userRepo.FindUserBannedList(entity.UserBanned{
		StreamerID: streamer.ID,
	})
	listUser := []entity.Users{}
	if err != nil {
		return listUser, err
	}

	for _, userBanned := range userBanneds {
		listUser = append(listUser, userBanned.User)
	}

	return listUser, nil
}
func NewuserBannedUsecase(db db.Database) *userBannedUsecase {
	repo := repository.NewUserBannedRepository(db)
	return &userBannedUsecase{
		userRepo: *repo,
	}
}
