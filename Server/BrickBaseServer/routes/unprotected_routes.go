package routes

import (
	controller "github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupUnProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	api := router.Group("/api")
	{
		api.GET("/sets", controller.GetSets(client))
		api.GET("/set/:set_num", controller.GetSet(client))
		api.POST("/register", controller.RegisterUser(client))
		api.POST("/login", controller.LoginUser(client))
		api.POST("/logout", controller.LogoutHandler(client))
		api.GET("/bestsets", controller.GetBestSets(client))
		api.POST("/refresh", controller.RefreshTokenHandler(client))
	}
}
