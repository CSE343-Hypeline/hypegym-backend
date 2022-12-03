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
}

func main() {
	router := initRouter()
	router.Run()
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", controllers.Home)
	router.POST("/login", controllers.UserLogin)

	api := router.Group("/api").Use(middlewares.Auth())
	{
		api.GET("/me", controllers.Me)

		api.GET("/users/:gymID", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "ADMIN"}), controllers.UserGetAll)
		api.GET("/users/pts/:gymID", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "ADMIN"}), controllers.UserGetAllByRole)
		api.GET("/users/members/:gymID", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "ADMIN"}), controllers.UserGetAllByRole)
		api.POST("/user", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "ADMIN"}), controllers.UserCreate)
		api.GET("/user/:id", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "ADMIN"}), controllers.UserGet)
		api.DELETE("/user/:id", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "ADMIN"}), controllers.UserDelete)

		api.POST("/gym", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "ADMIN"}), controllers.GymCreate)
	}
	return router
}
