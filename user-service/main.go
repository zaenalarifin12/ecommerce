package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/exaring/otelpgx"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"golang.org/x/time/rate"
	"log"
	_ "user-service/docs"
	"user-service/internal/api"
	"user-service/internal/config"
	"user-service/internal/repository"
	"user-service/internal/service"
)

//	@title			User Service API
//	@version		1.0
//	@description	This is a sample user server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:5000
//	@BasePath	/api/v1

// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @in header
// @name Authorization

func main() {
	initTracer()

	conf, err := config.LoadConfig(".")
	if err != nil {
		panic("can't load config")
	}

	runMigrationDB(conf.MigrationUrl, conf.DBSource)
	pool, err := connectToDatabase(conf)

	// repository
	userRepository := repository.NewUser(pool)

	// services
	userService := service.NewUser(userRepository)

	// setup gin
	router := gin.Default()

	// open telemetry
	router.Use(otelgin.Middleware("user-service"))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// list router
	api.NewAuth(router, userService, rate.Limit(1e9), 100)
	//api.NewAuth(router, userService, rate.Limit(0.0033333), 100) // 2 second 5 hit
	api.NewUser(router, userService)

	// swagger

	port := fmt.Sprintf(":%v", conf.ServerPort)

	err = router.Run(port)
	if err != nil {
		panic("can't start server")
	}
}

func connectToDatabase(conf config.Config) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(conf.DBSource)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
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
}

func UpMigration(migration *migrate.Migrate) {
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to run migrate up! ", err)
	}

	log.Println("db migrate successfully")

}

func initTracer() {
	jaegerEndpoint := "http://localhost:14268/api/traces"

	jaegerExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)),
	)
	if err != nil {
		log.Fatalln("Couldn't initialize exporter", err)
	}

	// Create stdout exporter to be able to retrieve the collected spans.
	_, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalln("Couldn't initialize exporter", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(jaegerExporter),
		trace.WithResource(resource.NewSchemaless(attribute.KeyValue{
			Key:   semconv.ServiceNameKey,
			Value: attribute.StringValue("user-service"),
		})),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}
