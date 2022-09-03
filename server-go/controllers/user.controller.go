package controllers

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"example/web-service-gin/models"
	"example/web-service-gin/services"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	cart, _ := uc.userService.GetCart(currentUser.ID)
	history, _ := uc.userService.GetHistory(currentUser.ID)

	fmt.Println(history)

	ctx.JSON(http.StatusOK, models.FilteredResponse(currentUser, cart.Items))
}

func (pc *UserController) AddToCart(ctx *gin.Context) {
	var queryParams = ctx.Request.URL.Query()
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	var productId string

	if len(queryParams["productId"]) != 0 {
		productId = queryParams["productId"][0]
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	cart, err := pc.userService.AddToCart(currentUser.ID, productId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "cart": cart})
	return
}

func (pc *UserController) RemoveFromCart(ctx *gin.Context) {
	var queryParams = ctx.Request.URL.Query()
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	var productId string

	if len(queryParams["productId"]) != 0 {
		productId = queryParams["productId"][0]
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	cart, err := pc.userService.RemoveFromCart(currentUser.ID, productId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "cart": cart})
	return
}

func (pc *UserController) SuccessBuy(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)
	var paymentData *models.SuccessPaymentInput

	fmt.Println(paymentData)

	if err := ctx.ShouldBindJSON(&paymentData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	pc.userService.CreatePaymentHistory(currentUser.ID, paymentData)
	pc.userService.ClearCart(currentUser.ID)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
	return
}

func (pc *UserController) GetHistory(ctx *gin.Context) {
	fmt.Println("historyyyy")
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	paymentHistory, _ := pc.userService.GetHistory(currentUser.ID)

	ctx.JSON(http.StatusOK, gin.H{"success": true, "history": paymentHistory})
	return
}
