package handler

import (
	"api/internal/pkg/domain/domain_model/dto"
	"api/internal/pkg/usecase"
	"api/pkg/infrastucture/db"
	"api/pkg/share/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomUsecase usecase.RoomUsecase
}

func NewRoomHandler(db db.Database) *RoomHandler {
	u := usecase.NewRoomUsecase(db)
	return &RoomHandler{
		roomUsecase: u,
	}
}

func (h *RoomHandler) GetRoomList(c *gin.Context) {
	req := dto.GetRoomListRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}

	res, err := h.roomUsecase.GetRoomList(req)
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
func (h *RoomHandler) GetRoomInfo(c *gin.Context) {
	roomId := c.Query("room_id")

	if roomId == "" {
		data := dto.BaseResponse{
			Status: http.StatusBadRequest,
			Error:  "room id is required",
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}

	res, err := h.roomUsecase.GetRoomInfo(roomId)
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
func (h *RoomHandler) StartRoom(c *gin.Context) {
	req := dto.StartRoomRequest{}
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
	res, err := h.roomUsecase.StartRoom(user, req)
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
func (h *RoomHandler) EndRoom(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	res, err := h.roomUsecase.EndRoom(user)
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
