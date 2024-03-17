package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zaenalarifin12/product-service/docs"
	"github.com/zaenalarifin12/product-service/internal/api"
	"github.com/zaenalarifin12/product-service/internal/config"
	"github.com/zaenalarifin12/product-service/internal/middleware"
	"github.com/zaenalarifin12/product-service/internal/repository"
	"github.com/zaenalarifin12/product-service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// @title Product Service API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5001
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @in header
// @name Authorization

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		panic("can't load config")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	dbName := client.Database("db-product-service")
	collectionName := "products"

	// Creating a new  repository
	productRepo := repository.NewProductRepository(dbName, collectionName)

	// creating services
	productService := service.NewProductService(productRepo)

	// Initialize Gin router
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// list router
	api.NewProduct(r, productService, middleware.JWTMiddleware())

	// Run server on port 8080
	if err := r.Run(fmt.Sprintf(":%v", conf.ServerPort)); err != nil {
		panic(err)
	}
}
