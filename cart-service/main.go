package main

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/cart-service/docs"
	"github.com/cart-service/internal/api"
	"github.com/cart-service/internal/config"
	"github.com/cart-service/internal/middleware"
	"github.com/cart-service/internal/repository"
	"github.com/cart-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
)

//	@title			Cart Service API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:5002
//	@BasePath	/

// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @in header
// @name Authorization

// @tag.name Carts
// @tag.name TransactionCarts

func main() {

	conf, err := config.LoadConfig(".")
	if err != nil {
		panic("can't load config")
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.RedisURL,      // Assuming you have Redis URL in your config
		Password: conf.RedisPassword, // If Redis is password protected
		DB:       conf.RedisDB,       // Select Redis DB
	})

	// Ping Redis to test the connection
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic("can't connect to Redis")
	}
	fmt.Println("Connected to Redis:", pong)
	runMigrationDB(conf.MigrationUrl, conf.DBSource)
	pool, err := connectToDatabase(conf)

	//--------
	// repository
	cartRepository := repository.NewCart(pool, redisClient)

	// services
	cartService := service.NewCartService(cartRepository)

	// setup gin
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// list router
	api.NewCart(router, cartService, middleware.JWTMiddleware())

	port := fmt.Sprintf(":%v", conf.ServerPort)

	err = router.Run(port)
	if err != nil {
		panic("can't start server")
	}
}

func connectToDatabase(conf config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), conf.DBSource)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func runMigrationDB(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance: ", err)
	}

	UpMigration(migration)
	//rollbackMigration(migration)

}

func rollbackMigration(migration *migrate.Migrate) {
	// Perform rollback
	if err := migration.Down(); err != nil {
		log.Fatal("failed to rollback migration: ", err)
	}
	log.Println("rollback successful")
	os.Exit(1)
}

func UpMigration(migration *migrate.Migrate) {
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to run migrate up! ", err)
	}

	log.Println("db migrate successfully")

}
