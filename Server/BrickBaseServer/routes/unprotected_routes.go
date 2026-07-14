package routes

import (
	controller "github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine) {
	router.GET("/sets", controller.GetSets())
	router.GET("/set/:set_num", controller.GetSet())
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
}
