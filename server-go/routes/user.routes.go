package routes

import (
	"github.com/gin-gonic/gin"
	"example/web-service-gin/controllers"
	"example/web-service-gin/middleware"
	"example/web-service-gin/services"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup, userService services.UserService) {

	router := rg.Group("users")
	router.Use(middleware.DeserializeUser(userService))
	router.GET("/auth", uc.userController.GetMe)
	router.POST("/addToCart", uc.userController.AddToCart)
	router.DELETE("/removeFromCart", uc.userController.RemoveFromCart)
	router.POST("/successBuy", uc.userController.SuccessBuy)
	router.GET("/history", uc.userController.GetHistory)
}