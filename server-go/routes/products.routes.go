package routes

import (
	"github.com/gin-gonic/gin"
	"example/web-service-gin/controllers"
	//"example/web-service-gin/middleware"
	"example/web-service-gin/services"
)

type ProductsRouteController struct {
	productsController controllers.ProductsController
}

func NewProductsRouteController(productsController controllers.ProductsController) ProductsRouteController {
	return ProductsRouteController{productsController}
}

func (rc *ProductsRouteController) ProductsRoute(rg *gin.RouterGroup, productsService services.ProductsService) {
	router := rg.Group("/product")

	router.GET("/getProducts", rc.productsController.GetProducts)
	router.POST("/uploadProduct", rc.productsController.UploadProduct)
	router.GET("", rc.productsController.GetProduct)
	//router.POST("/uploadImage", rc.productsController.UploadImage)
}