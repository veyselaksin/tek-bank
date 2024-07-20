package config

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"tek-bank/pkg/cresponse"

	"github.com/gofiber/fiber/v2/log"
)

type ServerConfig struct {
	Host string
	Port string
}

var FiberConfig = fiber.Config{
	AppName:   "Tek-Bank",
	BodyLimit: 1024 * 1024 * 50, // 50 MB

	// Override default error handler
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		// Status code defaults to 500
		var code int = fiber.StatusInternalServerError

		// Retrieve the custom status code if it's a fiber.*Error
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		log.Error("Error occurred: ", err)

		return cresponse.ErrorResponse(ctx, code, "Unexpected error occurred")
	},
}

func GetLanguage(ctx *fiber.Ctx) string {
	return ctx.Get("Accept-Language", "en")
}
