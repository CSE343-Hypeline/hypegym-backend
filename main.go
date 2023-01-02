package main

import (
	"hypegym-backend/controllers"
	"hypegym-backend/initializers"
	"hypegym-backend/middlewares"
	"hypegym-backend/models/enums"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.MigrateDB()
	initializers.ConnectToRedis()
}

func main() {
	router := initRouter()
	router.Run()
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", controllers.Home)
	router.POST("/login", controllers.Login)

	api := router.Group("/api")
	api.Use(middlewares.Auth())
	{
		initSharedAPI(api)
		initAdminAPI(api)
		initMemberAPI(api)
	}
	return router
}

func initSharedAPI(api *gin.RouterGroup) {
	sharedAPI := api.Group("/")
	{
		sharedAPI.GET("/me", controllers.Me)
		sharedAPI.POST("/logout", controllers.Logout)
	}
}

func initAdminAPI(api *gin.RouterGroup) {
	adminAPI := api.Group("/", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "ADMIN"}))
	{
		adminAPI.GET("/users/:gymID", controllers.UserGetAll)
		adminAPI.GET("/users/pts/:gymID", controllers.UserGetAllByRole)
		adminAPI.GET("/users/members/:gymID", controllers.UserGetAllByRole)
		adminAPI.POST("/user", controllers.UserCreate)
		adminAPI.GET("/user/:id", controllers.UserGet)
		adminAPI.DELETE("/user/:id", controllers.UserDelete)

		adminAPI.POST("/pt/:id/assign-member", controllers.AssignMembers)
		adminAPI.DELETE("/pt/:id/dismiss-member", controllers.DismissMember)
		adminAPI.GET("/pt/:id/members", controllers.GetMembers)
		adminAPI.GET("/gym/:gymID/attendance/month", controllers.GetAttendance)
		adminAPI.GET("/gym/:gymID/attendance/day", controllers.GetAttendance)
		adminAPI.GET("/gym/:gymID/online", controllers.GetAllOnlines)

		adminAPI.POST("/gym", controllers.GymCreate)
	}
}

func initMemberAPI(api *gin.RouterGroup) {
	memberAPI := api.Group("/", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "MEMBER"}))
	{
		memberAPI.POST("/member/:userID/checkIn/:gymID", controllers.CheckIn)
		memberAPI.POST("/member/:userID/checkOut/:gymID", controllers.CheckOut)
	}
}
