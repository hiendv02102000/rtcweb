package handler

import (
	"api/internal/pkg/domain/domain_model/dto"
	"api/internal/pkg/usecase"
	"api/pkg/infrastucture/db"
	"api/pkg/share/middleware"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewCustomerHandler(db db.Database) *UserHandler {
	u := usecase.NewCustomerUsecase(db)
	return &UserHandler{
		userUsecase: u,
	}
}

func (h *UserHandler) Login(c *gin.Context) {

	req := dto.LoginRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	tokenString, err := h.userUsecase.Login(req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	data := dto.BaseResponse{
		Status: http.StatusOK,
		Result: tokenString,
	}
	c.JSON(http.StatusOK, data)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	req := dto.CreateUserRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	res, err := h.userUsecase.CreateUser(req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	data := dto.BaseResponse{
		Status: http.StatusOK,
		Result: res,
	}
	c.JSON(http.StatusOK, data)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	req := dto.UpdateProfileRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	user := middleware.GetUserFromContext(c)
	file, _ := c.FormFile("avatar")
	var ioFile multipart.File
	ioFile = nil
	if file != nil {
		ioFile, err = file.Open()
		if err != nil {

			data := dto.BaseResponse{
				Status: http.StatusBadRequest,
				Error:  err.Error(),
			}
			c.JSON(http.StatusBadRequest, data)
			return
		}
	}

	res, err := h.userUsecase.UpdateProfile(req, user, ioFile)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	data := dto.BaseResponse{
		Status: http.StatusOK,
		Result: res,
	}
	c.JSON(http.StatusOK, data)
}

func (h *UserHandler) ChangePassWord(c *gin.Context) {
	req := dto.ChangePassWordRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	user := middleware.GetUserFromContext(c)

	res, err := h.userUsecase.ChangePassWord(req, user)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	data := dto.BaseResponse{
		Status: http.StatusOK,
		Result: res,
	}
	c.JSON(http.StatusOK, data)
}
func (h *UserHandler) GetProfile(c *gin.Context) {

	user := middleware.GetUserFromContext(c)
	res := h.userUsecase.GetProfile(user)
	data := dto.BaseResponse{
		Status: http.StatusOK,
		Result: res,
	}
	c.JSON(http.StatusOK, data)
}
func (h *UserHandler) GetUsersInRoom(c *gin.Context) {
	pageS := c.Query("page")
	sizeS := c.Query("size")
	page, errP := strconv.Atoi(pageS)
	size, errS := strconv.Atoi(sizeS)
	if errP != nil || errS != nil {

		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  "page size is int",
		}
		c.JSON(http.StatusBadRequest, data)
		return

	}
	roomId := c.Query("room_id")
	res, err := h.userUsecase.GetAllUserInRoom((int64)(page), (int64)(size), roomId)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	data := dto.BaseResponse{
		Status: http.StatusOK,
		Result: res,
	}
	c.JSON(http.StatusOK, data)
}
func (h *UserHandler) CheckLogin(c *gin.Context) {
	data := dto.BaseResponse{
		Status: http.StatusOK,
		Result: "login",
	}
	c.JSON(http.StatusOK, data)
}
