package cresponse

import "github.com/gofiber/fiber/v2"

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(ctx *fiber.Ctx, status int, data interface{}, msg ...string) error {
	if len(msg) == 0 {
		msg = append(msg, "Success")
	}

	return ctx.Status(status).JSON(BaseResponse{
		Success: true,
		Message: msg[0],
		Data:    data,
	})
}

func ErrorResponse(ctx *fiber.Ctx, status int, msg string, data ...interface{}) error {
	if len(data) == 0 {
		data = nil
	}

	return ctx.Status(status).JSON(BaseResponse{
		Success: false,
		Message: msg,
		Data:    data,
	})
}

func RedirectResponse(ctx *fiber.Ctx, url string) error {
	return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
}
