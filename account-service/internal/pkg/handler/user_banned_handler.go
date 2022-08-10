package handler

import (
	"api/internal/pkg/domain/domain_model/dto"
	"api/internal/pkg/usecase"
	"api/pkg/infrastucture/db"
	"api/pkg/share/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserBannedHandler struct {
	userBannedUsecase usecase.UserBannedUsecase
}

func NewuserBannedHandler(db db.Database) *UserBannedHandler {
	u := usecase.NewuserBannedUsecase(db)
	return &UserBannedHandler{
		userBannedUsecase: u,
	}
}

func (h *UserBannedHandler) BanUser(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	type UserIdRequest struct {
		UserId int `json:"user_id" form:"user_id" binding:"required"`
	}
	req := UserIdRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	err = h.userBannedUsecase.BanUser(user, req.UserId)
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
		Result: nil,
	}
	c.JSON(http.StatusOK, data)
}
func (h *UserBannedHandler) UnBanUser(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	type UserIdRequest struct {
		UserId int `json:"user_id" form:"user_id" binding:"required"`
	}
	req := UserIdRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}
	err = h.userBannedUsecase.UnBanUser(user, req.UserId)
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
		Result: nil,
	}
	c.JSON(http.StatusOK, data)
}
func (h *UserBannedHandler) GetUserBannedList(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	res, err := h.userBannedUsecase.GetUserBannedList(user)
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
