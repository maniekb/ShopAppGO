package controllers

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"github.com/gin-gonic/gin"
	"example/web-service-gin/config"
	"example/web-service-gin/models"
	"example/web-service-gin/services"
	"example/web-service-gin/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthController(authService services.AuthService, userService services.UserService) AuthController {
	return AuthController{authService, userService}
}

func (ac *AuthController) GetAuthCodeNotProvidedErrorMessage () (string) {
	return "Authorization code not provided!"
}


// SignUp User
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if user.Password != user.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	newUser, err := ac.authService.SignUpUser(user)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "data": gin.H{"user": models.FilteredResponse(newUser, []models.CartItemDBResponse {})}})
}

// SignIn User
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserByEmail(credentials.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email or Password"})
		return
	}

	config, _ := config.LoadConfig(".")

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"success": true, "access_token": access_token})
}

// Refresh Access Token
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := config.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserById(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (ac *AuthController) InitGoogleLogin(ctx *gin.Context) {
    var OAuth2Config = utils.GetGoogleOAuthConfig()
    url := OAuth2Config.AuthCodeURL("/")
    ctx.Redirect(http.StatusFound, url)
}

func (ac *AuthController) GoogleOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": ac.GetAuthCodeNotProvidedErrorMessage})
		return
	}

	// Use the code to get the id and access tokens
	tokenRes, err := utils.GetGoogleOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}

	user, err := utils.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}

	createdAt := time.Now()
	resBody := &models.UpdateDBUser{
		Email:     user.Email,
		Name:      user.Name,
		Photo:     user.Picture,
		Provider:  "google",
		Verified:  true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	updatedUser, err := ac.userService.UpsertUser(user.Email, resBody)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}

	config, _ := config.LoadConfig(".")

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, updatedUser.ID.Hex(), config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, updatedUser.ID.Hex(), config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.ClientOrigin, pathUrl))
}

func (ac *AuthController) InitGitHubLogin(ctx *gin.Context) {
    var OAuth2Config = utils.GetGitHubOAuthConfig()
    url := OAuth2Config.AuthCodeURL("/")
    ctx.Redirect(http.StatusFound, url)
}

func (ac *AuthController) GitHubOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": ac.GetAuthCodeNotProvidedErrorMessage})
		return
	}

	// Use the code to get the id and access tokens
	tokenRes, err := utils.GetGitHubOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}

	user, err := utils.GetGitHubUser(tokenRes.Access_token)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}

	createdAt := time.Now()
	resBody := &models.UpdateDBUser{
		Email:     user.Email,
		Name:      user.Name,
		Photo:     user.Photo,
		Provider:  "github",
		Verified:  true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	updatedUser, err := ac.userService.UpsertUser(user.Email, resBody)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}

	config, _ := config.LoadConfig(".")

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, updatedUser.ID.Hex(), config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, updatedUser.ID.Hex(), config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.ClientOrigin, pathUrl))
}

func (ac *AuthController) InitFacebookLogin(ctx *gin.Context) {
    var OAuth2Config = utils.GetFacebookOAuthConfig()
    url := OAuth2Config.AuthCodeURL("/")
    ctx.Redirect(http.StatusFound, url)
}

func (ac *AuthController) FacebookOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": ac.GetAuthCodeNotProvidedErrorMessage})
		return
	}

	var OAuth2Config = utils.GetFacebookOAuthConfig()

    token, err := OAuth2Config.Exchange(oauth2.NoContext, code)

    if err != nil || token == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

    user, err := utils.GetFacebookUser(token.AccessToken)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}

	createdAt := time.Now()
	resBody := &models.UpdateDBUser{
		Email:     user.Email,
		Name:      user.Name,
		Provider:  "facebook",
		Verified:  true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	updatedUser, err := ac.userService.UpsertUser(user.Email, resBody)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}

	config, _ := config.LoadConfig(".")

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, updatedUser.ID.Hex(), config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, updatedUser.ID.Hex(), config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

    ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.ClientOrigin, pathUrl))
}


func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}