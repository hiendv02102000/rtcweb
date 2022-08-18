package middleware

import (
	"chat-service/dto"
	"chat-service/pkg/share/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) dto.UserProfileResponse {

	value, exist := c.Get("user")
	if !exist {
		return dto.UserProfileResponse{}
	}
	return value.(dto.UserProfileResponse)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")
		extractedToken := strings.Split(clientToken, "Bearer ")
		clientToken = strings.TrimSpace(extractedToken[1])
		dataJ, err := utils.SendRequest("GET", utils.HOST_ACCOUNT_SERVICE+"/api/account/profile", clientToken, nil)
		if err != nil {
			data := dto.BaseResponse{
				Status: http.StatusUnauthorized,
				Error:  err.Error(),
			}
			c.JSON(http.StatusUnauthorized, data)
			c.Abort()
			return
		}
		res := dto.BaseResponse{}
		err = json.Unmarshal(dataJ, &res)
		if err != nil {
			data := dto.BaseResponse{
				Status: http.StatusUnauthorized,
				Error:  err.Error(),
			}
			c.JSON(http.StatusUnauthorized, data)
			c.Abort()
			return
		}
		if http.StatusOK != res.Status {
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}
		pro := dto.UserProfileResponse{}
		utils.ConvertToObject(res.Result, &pro)
		c.Set("user", pro)
		c.Next()
	}
}
func AuthUserBanned() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")
		extractedToken := strings.Split(clientToken, "Bearer ")
		clientToken = strings.TrimSpace(extractedToken[1])
		roomId := c.Query("room_id")
		dataJ, err := utils.SendRequest("GET", utils.HOST_ACCOUNT_SERVICE+"/api/account/check_banned_user?room_id="+roomId, clientToken, nil)
		if err != nil {
			data := dto.BaseResponse{
				Status: http.StatusUnauthorized,
				Error:  err.Error(),
			}
			c.JSON(http.StatusUnauthorized, data)
			c.Abort()
			return
		}
		res := dto.BaseResponse{}
		err = json.Unmarshal(dataJ, &res)
		if err != nil {
			data := dto.BaseResponse{
				Status: http.StatusUnauthorized,
				Error:  err.Error(),
			}
			c.JSON(http.StatusUnauthorized, data)
			c.Abort()
			return
		}
		if http.StatusOK != res.Status {
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}
		isBanned, ok := res.Result.(bool)
		if isBanned || !ok {
			res.Status = http.StatusUnauthorized
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
		}
		c.Next()
	}
}
