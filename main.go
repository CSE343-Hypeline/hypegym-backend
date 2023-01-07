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
	router.POST("/contact", controllers.Contact)

	api := router.Group("/api")
	api.Use(middlewares.Auth())
	{
		initSharedAPI(api)
		initAdminAPI(api)
		initMemberAPI(api)
		initTrainerAPI(api)
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
		adminAPI.PUT("/user/:id", controllers.UserUpdate)

		adminAPI.POST("/pt/:id/assign-member", controllers.AssignMembers)
		adminAPI.DELETE("/pt/:id/dismiss-member", controllers.DismissMember)
		adminAPI.GET("/pt/:id/members", controllers.GetMembers)
		adminAPI.GET("/gym/:gymID/attendance/month", controllers.GetAttendance)
		adminAPI.GET("/gym/:gymID/attendance/day", controllers.GetAttendance)
		adminAPI.GET("/gym/:gymID/online", controllers.GetAllOnlines)

		adminAPI.POST("/gym", controllers.GymCreate)

		adminAPI.GET("/exercises", controllers.GetExercises)

	}
}

func initTrainerAPI(api *gin.RouterGroup) {
	trainerAPI := api.Group("/", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "PT"}))
	{
		trainerAPI.POST("/member/assign-program/:id", controllers.AssignProgram)
		trainerAPI.POST("/member/assign-programs/:id", controllers.AssignPrograms)
		trainerAPI.DELETE("/member/dismiss-program/:id", controllers.DismissProgram)
	}
}

func initMemberAPI(api *gin.RouterGroup) {
	memberAPI := api.Group("/", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "MEMBER"}))
	{
		memberAPI.POST("/member/:userID/checkIn/:gymID", controllers.CheckIn)
		memberAPI.POST("/member/:userID/checkOut/:gymID", controllers.CheckOut)

		memberAPI.GET("/member/programs/:id", controllers.GetPrograms)
	}
}
