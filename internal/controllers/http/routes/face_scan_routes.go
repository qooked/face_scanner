package routes

import (
	"github.com/gofiber/fiber/v2"
)

type FaceScanHandlers interface {
	ExtendFaceScannerTask(c *fiber.Ctx) error
	GetFaceScannerTask(c *fiber.Ctx) error
	StartFaceScannerTask(c *fiber.Ctx) error
	DeleteFaceScannerTask(c *fiber.Ctx) error
}

func AttachCascadeRoutes(router fiber.Router, handlers FaceScanHandlers) {

	router.Post("extend", handlers.ExtendFaceScannerTask)
	router.Get("get", handlers.GetFaceScannerTask)
	router.Post("start", handlers.StartFaceScannerTask)
	router.Delete("delete", handlers.DeleteFaceScannerTask)

}
