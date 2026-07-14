package routes

import (
	controller "github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/controllers"
	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleWare())

	router.POST("/addset", controller.AddSet())
	router.PUT("/set/:set_num", controller.UpdateSet())
	router.DELETE("/set/:set_num", controller.DeleteSet())

}
