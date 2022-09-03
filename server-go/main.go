package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"example/web-service-gin/config"
	"example/web-service-gin/controllers"
	"example/web-service-gin/routes"
	"example/web-service-gin/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection         *mongo.Collection
	productsCollection	   *mongo.Collection
	cartsCollection	   	   *mongo.Collection
	paymentsCollection	   *mongo.Collection

	authService            services.AuthService
	AuthController         controllers.AuthController
	AuthRouteController    routes.AuthRouteController
	SessionRouteController routes.SessionRouteController

	productsService services.ProductsService
	ProductsController controllers.ProductsController
	ProductsRouteController routes.ProductsRouteController
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// Collections
	authCollection = mongoclient.Database("golang_mongodb").Collection("users")
	productsCollection = mongoclient.Database("golang_mongodb").Collection("products")
	cartsCollection = mongoclient.Database("golang_mongodb").Collection("carts")
	paymentsCollection = mongoclient.Database("golang_mongodb").Collection("payments")

	userService = services.NewUserServiceImpl(authCollection, cartsCollection, productsCollection, paymentsCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	productsService = services.NewProductsServiceImpl(productsCollection, ctx)

	AuthController = controllers.NewAuthController(authService, userService)
	AuthRouteController = routes.NewAuthRouteController(AuthController)
	SessionRouteController = routes.NewSessionRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	ProductsController = controllers.NewProductsController(productsService)
	ProductsRouteController = routes.NewProductsRouteController(ProductsController)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "ok"})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	ProductsRouteController.ProductsRoute(router, productsService)
	SessionRouteController.SessionRoute(router)
	log.Fatal(server.Run(":" + config.Port))
}