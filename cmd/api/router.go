package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"tek-bank/cmd/api/handler/v1/account"
	"tek-bank/cmd/api/handler/v1/auth"
	"tek-bank/cmd/api/handler/v1/profile"
	"tek-bank/cmd/api/middleware/authware"
	"tek-bank/cmd/api/middleware/transaction"
	"tek-bank/internal/db/repository"
	"tek-bank/internal/service"
	"tek-bank/pkg/converter"
	"tek-bank/pkg/crypto"
)

// HealthCheck godoc
// @Summary Health Check API
// @Description Health Check for the API
// @Tags Health Check
// @Accept application/json
// @Produce application/json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func health(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}

func InitializeRouters(app *fiber.App, connection *gorm.DB, redis *redis.Client) {

	// Middleware
	authorizationConfig := authware.Config{
		DBConnection:            connection,
		AuthorizationHeaderKey:  "Authorization",
		AuthorizationTypeBearer: "Bearer",
	}

	authentication := authware.New(authorizationConfig)

	// Packages
	pkgConverter := converter.NewConverter()
	pkgCrypto := crypto.NewCrypto()

	// Repositories
	userRepository := repository.NewUserRepository(connection)
	accountRepository := repository.NewAccountRepository(connection, redis)
	transferHistoryRepository := repository.NewTransferHistoryRepository(connection)

	// Services
	authService := service.NewAuthService(userRepository, pkgCrypto)
	accountService := service.NewAccountService(accountRepository, userRepository, transferHistoryRepository, pkgCrypto, pkgConverter)
	profileService := service.NewProfileService(accountRepository, transferHistoryRepository, userRepository)

	// Handlers
	authHandler := auth.NewAuthHandler(authService)
	accountHandler := account.NewAccountHandler(accountService)
	profileHandler := profile.NewProfileHandler(profileService)

	// Initialize the routes for the application here
	v1 := app.Group("/v1")

	// Swagger documentation
	v1.Get("/docs/*", swagger.HandlerDefault)

	// Health check
	v1.Get("/health", health)

	// Initialize the routes for the application here
	// Auth routes
	authRouter := v1.Group("/auth")
	authRouter.Post("/login", authHandler.Login)
	authRouter.Get("/user-info", authentication, authHandler.GetUserInfo)

	// Account routes
	accountRouter := v1.Group("/account")
	accountRouter.Post("/register", transaction.Tx(connection), accountHandler.RegisterAccount)
	accountRouter.Post("/create", transaction.Tx(connection), accountHandler.CreateNewAccount)
	accountRouter.Put("/add-money/:accountNumber", authentication, transaction.Tx(connection), accountHandler.AddMoney)
	accountRouter.Post("/transfer", authentication, transaction.Tx(connection), accountHandler.TransferMoney)
	accountRouter.Get("/transfer-approval", transaction.Tx(connection), accountHandler.TransferApproval)

	// Profile routes
	profileRouter := v1.Group("/profile")
	profileRouter.Get("/", authentication, profileHandler.MyProfile)
	profileRouter.Get("/transfer-history", authentication, profileHandler.MyTransferHistory)

}
