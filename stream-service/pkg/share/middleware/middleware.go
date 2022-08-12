package middleware

import (
	"encoding/json"
	"net/http"
	"stream-service/dto"
	"stream-service/pkg/share/utils"
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
		utils.ConvertToPbject(res.Result, &pro)
		c.Set("user", pro)
		c.Next()
	}
}
func AuthRoomMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
