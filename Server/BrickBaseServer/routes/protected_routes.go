package routes

import (
	controller "github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/controllers"
	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	apiProtected := router.Group("/api")
	apiProtected.Use(middleware.AuthMiddleWare())
	{
		apiProtected.POST("/addset", controller.AddSet(client))
		apiProtected.PUT("/set/:set_num", controller.UpdateSet(client))
		apiProtected.DELETE("/set/:set_num", controller.DeleteSet(client))
		apiProtected.PATCH("/sets/:set_num/review", controller.AddUserReview(client))
	}

}
