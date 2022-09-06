package routes

import (
	"github.com/gin-gonic/gin"
	"example/web-service-gin/controllers"
)

type SessionRouteController struct {
	authController controllers.AuthController
}

func NewSessionRouteController(authController controllers.AuthController) SessionRouteController {
	return SessionRouteController{authController}
}

func (rc *SessionRouteController) SessionRoute(rg *gin.RouterGroup) {
	router := rg.Group("/sessions/oauth")

	router.GET("/google/init", rc.authController.InitGoogleLogin)
	router.GET("/google", rc.authController.GoogleOAuth)
	router.GET("/github/init", rc.authController.InitGitHubLogin)
	router.GET("/github", rc.authController.GitHubOAuth)
	router.GET("/facebook/init", rc.authController.InitFacebookLogin)
	router.GET("/facebook", rc.authController.FacebookOAuth)
}