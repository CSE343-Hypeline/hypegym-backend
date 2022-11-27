package main

import (
	"hypegym-backend/controllers"
	"hypegym-backend/initializers"
	"hypegym-backend/middlewares"
	"hypegym-backend/models/enums"

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
	router.GET("/", controllers.Home)
	router.POST("/login", controllers.UserLogin)

	api := router.Group("/api").Use(middlewares.Auth())
	{
		api.GET("/ping", controllers.Ping)
		api.POST("/user", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "ADMIN"}), controllers.UserCreate)
		api.POST("/gym", middlewares.AccessControl([]enums.Role{"SUPERADMIN", "ADMIN"}), controllers.GymCreate)
	}
	return router
}
