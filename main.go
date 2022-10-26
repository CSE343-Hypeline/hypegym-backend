package main

import (
	"hypegym-backend/controllers"
	"hypegym-backend/database"
	"hypegym-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect("root:@tcp(localhost:3306)/jwt_demo?parseTime=true")
	database.Migrate()

	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/", controllers.Home)
		api.POST("/login", controllers.GenerateToken)
		api.POST("/user/register", controllers.RegisterUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", controllers.Ping)
		}
	}
	return router
}
