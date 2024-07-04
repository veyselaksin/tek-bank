package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
	"tek-bank/cmd/api/handler/v1/auth"
	"tek-bank/internal/db/repository"
	"tek-bank/internal/service"
)

func InitializeRouters(app *fiber.App, connection *gorm.DB) {
	// Repositories
	userRepository := repository.NewUserRepository(connection)

	// Services
	authService := service.NewAuthService(userRepository)

	// Handlers
	authHandler := auth.NewAuthHandler(authService)

	// Initialize the routes for the application here
	v1 := app.Group("/v1")

	// HealthCheck godoc
	// @Summary Health Check API
	// @Description Health Check for the API
	// @Tags Health Check
	// @Accept application/json
	// @Produce application/json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	v1.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Swagger documentation
	v1.Get("/docs/*", swagger.HandlerDefault)

	// Initialize the routes for the application here
	// Auth routes
	authRouter := v1.Group("/auth")
	authRouter.Post("/register", authHandler.Register)
	authRouter.Post("/login", authHandler.Login)

}
