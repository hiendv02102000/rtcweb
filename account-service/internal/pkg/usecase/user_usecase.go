package usecase

import (
	"api/internal/pkg/domain/domain_model/dto"
	"api/internal/pkg/domain/domain_model/entity"
	"api/internal/pkg/domain/repository"
	"context"
	"fmt"
	"strconv"
	"sync"

	"api/pkg/infrastucture/db"
	"api/pkg/share/middleware"
	"api/pkg/share/utils"
	"errors"
	"io"
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserUsecase interface {
	GetProfile(entity.Users) dto.GetProfileResponse
	ChangePassWord(dto.ChangePassWordRequest, entity.Users) (string, error)
	UpdateProfile(dto.UpdateProfileRequest, entity.Users, io.Reader) (dto.UpdateProfileResponse, error)
	CreateUser(dto.CreateUserRequest) (dto.CreateUserResponse, error)
	Login(dto.LoginRequest) (string, error)
	GetAllUserInRoom(page int64, size int64, roomId string) (dto.GetAllUserInRoomResponse, error)
}

type customerUsecase struct {
	userRepo repository.UserRepository
}

func (u *customerUsecase) GetProfile(user entity.Users) dto.GetProfileResponse {
	return dto.GetProfileResponse{
		Email:     user.Email,
		AvatarURL: user.AvatarUrl,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
func (u *customerUsecase) ChangePassWord(req dto.ChangePassWordRequest, user entity.Users) (string, error) {
	req.OldPassword = utils.EncryptPassword(req.OldPassword)
	if user.Password != req.OldPassword {
		return "", errors.New("password not matched")
	}

	timeNow := time.Now()

	timeExpiredAt := timeNow.Add(time.Hour * 48)
	// generate uuid
	uuid := uuid.Must(uuid.NewV4(), nil)
	tokenString, err := middleware.GenerateJWTToken(middleware.JWTParam{
		UUID:       uuid,
		Authorized: true,
		ExpriedAt:  timeExpiredAt,
	})

	if err != nil {
		return "", err
	}

	newUser := entity.Users{
		Token:          &tokenString,
		TokenExpiredAt: &timeExpiredAt,
		Password:       utils.EncryptPassword(req.NewPassword),
	}
	_, err = u.userRepo.UpdateUser(newUser, user)
	return tokenString, err
}
func (u *customerUsecase) UpdateProfile(req dto.UpdateProfileRequest, user entity.Users, image io.Reader) (dto.UpdateProfileResponse, error) {
	newUser := user
	newUser.FirstName = req.FirstName
	newUser.LastName = req.LastName
	if image != nil {
		url, errUpload := utils.UploadFile(image, user.Email+"avatar", []string{"avt_user"})
		if errUpload != nil {
			return dto.UpdateProfileResponse{}, errUpload
		}
		newUser.AvatarUrl = &url
	}
	user, err := u.userRepo.UpdateUser(newUser, user)
	if err != nil {
		return dto.UpdateProfileResponse{}, err
	}

	return dto.UpdateProfileResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		AvatarURL: user.AvatarUrl,
	}, nil
}
func (u *customerUsecase) CreateUser(req dto.CreateUserRequest) (dto.CreateUserResponse, error) {

	user, err := u.userRepo.FirstUser(entity.Users{
		Email: req.Email,
	})
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	if user.ID != 0 {
		err = errors.New(utils.USER_EXIST_ERROR)
		return dto.CreateUserResponse{}, err
	}
	user, err = u.userRepo.CreateUser(entity.Users{
		Email:     req.Email,
		Password:  utils.EncryptPassword(req.Password),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      entity.CustomerRole,
	})
	if err != nil {
		return dto.CreateUserResponse{}, err
	}
	return dto.CreateUserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
	}, nil
}
func (u *customerUsecase) Login(loginReq dto.LoginRequest) (string, error) {

	loginReq.Password = utils.EncryptPassword(loginReq.Password)

	user, err := u.userRepo.FirstUser(entity.Users{
		Email:    loginReq.Email,
		Password: loginReq.Password,
		Role:     entity.CustomerRole,
	})
	if err != nil {
		return "", err
	}

	if user.ID == 0 {
		err = errors.New(utils.LOGIN_FAIL_ERROR)
		return "", err
	}

	timeNow := time.Now()
	if user.TokenExpiredAt != nil && timeNow.After((*user.TokenExpiredAt).Add(time.Hour*2)) {
		return *user.Token, nil
	}
	timeExpiredAt := timeNow.Add(time.Hour * 48)
	// generate uuid
	uuid := uuid.Must(uuid.NewV4(), nil)
	tokenString, err := middleware.GenerateJWTToken(middleware.JWTParam{
		UUID:       uuid,
		Authorized: true,
		ExpriedAt:  timeExpiredAt,
	})

	if err != nil {
		return "", err
	}

	newUser := entity.Users{
		Token:          &tokenString,
		TokenExpiredAt: &timeExpiredAt,
	}
	_, err = u.userRepo.UpdateUser(newUser, user)
	return tokenString, err
}
func (u *customerUsecase) GetAllUserInRoom(page int64, size int64, roomId string) (dto.GetAllUserInRoomResponse, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 1
	}
	offset := (page - 1) * size
	listUserID, _ := db.RedisPool.ZRange(context.Background(), "room-user-"+roomId, offset, offset+size).Result()
	listUser := []entity.Users{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, idUser := range listUserID {
		go func(w *sync.WaitGroup, m *sync.Mutex, id string) {
			defer w.Done()
			i, erri := strconv.Atoi(id)
			if erri != nil {
				fmt.Println(erri)
				return
			}
			user, err := u.userRepo.FirstUser(entity.Users{ID: i})
			if user.ID == 0 || err != nil {
				return
			}
			m.Lock()
			defer m.Unlock()
			listUser = append(listUser, user)
		}(&wg, &mu, idUser)
	}
	wg.Wait()
	return dto.GetAllUserInRoomResponse{ListUser: listUser}, nil
}
func NewCustomerUsecase(db db.Database) *customerUsecase {
	repo := repository.NewUserRepository(db)
	return &customerUsecase{
		userRepo: *repo,
	}
}
