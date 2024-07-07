package transaction

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"tek-bank/internal/i18n"
	"tek-bank/internal/i18n/messages"
)

var DbTx *gorm.DB

func Tx(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txHandle := db.Begin()

		defer func() {
			if r := recover(); r != nil {
				log.Print("rolling back transaction due to panic: ", r)
				txHandle.Rollback()
			}
		}()

		c.Locals(DbTx, txHandle)
		err := c.Next()
		if err != nil {
			return err
		}

		if c.Response().StatusCode() >= fiber.StatusOK && c.Response().StatusCode() < fiber.StatusMultipleChoices {
			if err := txHandle.Commit().Error; err != nil {
				log.Print("tx commit error: ", err)
			}
		} else {
			log.Print("rolling back transaction due to status code: ", c.Response().StatusCode())
			txHandle.Rollback()
		}

		return nil
	}
}

func GetDbTx(ctx *fiber.Ctx) (*gorm.DB, error) {
	tx, ok := ctx.Locals(DbTx).(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	return tx, nil
}
