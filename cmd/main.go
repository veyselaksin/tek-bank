package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"tek-bank/cmd/api"
	"tek-bank/cmd/config"
	"tek-bank/docs"
	"tek-bank/internal/db/connection"
	"tek-bank/internal/i18n"
	"time"
)

var once sync.Once
var conn *gorm.DB
var redisConn *redis.Client

var serverConf config.ServerConfig

func init() {
	once.Do(func() {
		conn = connection.PostgresSQLConnection(connection.DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
			AppName:  os.Getenv("APP_NAME"),
			SSLMode:  os.Getenv("DB_SSL_MODE"),
			Timezone: os.Getenv("DB_TIMEZONE"),
		})

		redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			log.Error("Redis DB conversion error", err)
		}

		var redisPassword string
		if os.Getenv("REDIS_PASSWORD") != "" {
			redisPassword = os.Getenv("REDIS_PASSWORD")
		}

		// Initialize the Redis connection
		redisConn = connection.RedisConnection(connection.RedisConfig{
			Host:               os.Getenv("REDIS_HOST"),
			Port:               os.Getenv("REDIS_PORT"),
			Password:           redisPassword,
			DB:                 redisDB,
			TLSEnable:          os.Getenv("REDIS_TLS_ENABLE") == "true",
			InsecureSkipVerify: os.Getenv("REDIS_INSECURE_SKIP_VERIFY") == "true",
		})
	})

	// Initialize the config configuration
	serverConf = config.ServerConfig{
		Host: os.Getenv("APP_HOST"),
		Port: os.Getenv("APP_PORT"),
	}

	//Swagger Info configuration
	docs.SwaggerInfo.Host = fmt.Sprint(os.Getenv("SWAGGER_HOST"))

	//Init i18n
	i18n.InitBundle("./internal/i18n/languages/")
}

// @title Teknasyon Case Study API
// @version 1.0
// @description This is a config for Teknasyon Case Study API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
func main() {
	app := fiber.New(config.FiberConfig)

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05.000Z",
		TimeZone:   "Europe/Istanbul",
	}))

	// Initialize routes
	api.InitializeRouters(app, conn, redisConn)

	// Start listening on port 8000
	go func() {
		if err := app.Listen(":" + serverConf.Port); err != nil {
			panic(err)
		}
	}()

	// Graceful shutdown
	err := GracefulShutdown(app, 5*time.Second)
	if err != nil {
		log.Error("Graceful shutdown error", err)
	}
}

func GracefulShutdown(app *fiber.App, timeout time.Duration) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	sig := <-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db, err := conn.DB()
	if err != nil {
		return err
	}

	if err := shutdownDatabase(ctx, db); err != nil {
		return err
	}

	if err := app.Shutdown(); err != nil {
		return err
	}

	log.Infof("Signal received: %v", sig)
	return nil
}

func shutdownDatabase(ctx context.Context, db *sql.DB) error {
	ch := make(chan error, 1)
	go func() {
		ch <- db.Close()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		if err != nil {
			log.Error("Database close error", err)
			return err
		}
		return nil
	}
}
