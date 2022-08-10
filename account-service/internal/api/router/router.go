package router

import (
	"api/internal/pkg/handler"
	"api/pkg/infrastucture/db"
	"api/pkg/share/middleware"
	"api/pkg/share/validators"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
	DB     db.Database
}

func (r *Router) Routes() {

	r.DB.MigrateDBWithGorm()
	validators.CustomValidate()
	r.DB.MigrateDBWithGorm()

	hUserCus := handler.NewCustomerHandler(r.DB)
	hUserBaned := handler.NewuserBannedHandler(r.DB)
	hRoom := handler.NewRoomHandler(r.DB)
	api := r.Engine.Group("/api")
	{
		accountAPI := api.Group("/account")
		{
			accountAPI.POST("/login", hUserCus.Login)
			accountAPI.POST("/register", hUserCus.CreateUser)
			accountAPI.Use(middleware.AuthMiddleware(r.DB))
			{
				accountAPI.GET("/check_login", hUserCus.CheckLogin)
				accountAPI.GET("/profile", hUserCus.GetProfile)
				accountAPI.PATCH("/update_profile", hUserCus.UpdateProfile)
				accountAPI.PATCH("/change_password", hUserCus.ChangePassWord)
				accountAPI.POST("/ban_user", hUserBaned.BanUser)
				accountAPI.DELETE("/unban_user", hUserBaned.UnBanUser)
				accountAPI.GET("/get_userbanned_list", hUserBaned.GetUserBannedList)
				accountAPI.GET("/get_users_in_room", hUserCus.GetUsersInRoom)
			}

		}
		roomAPI := api.Group("/room")
		{
			roomAPI.GET("/get_room_list", hRoom.GetRoomList)
			roomAPI.GET("/get_room_info", hRoom.GetRoomList)
			roomAPI.Use(middleware.AuthRoomMiddleware(r.DB))
			{
				roomAPI.PUT("/start_room", hRoom.StartRoom)
				roomAPI.PATCH("/end_room", hRoom.EndRoom)
			}
		}
	}

}
func NewRouter() Router {
	var r Router
	r.Engine = gin.Default()
	database, err := db.NewDB()
	if err != nil {
		return Router{}
	}
	r.DB = database
	r.Routes()
	return r
}
