package utills

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func UnmarshalFromJSON(c *fiber.Ctx, data interface{}) (err error) {
	err = c.BodyParser(&data)
	if err != nil {
		err = fmt.Errorf("c.BodyParser(...): %w", err)
		return err
	}
	return nil
}
