package routes

import (
	"github.com/gofiber/fiber/v2"
)

type FaceScanHandlers interface {
	ExtendFaceScannerTask(c *fiber.Ctx) error
	GetFaceScannerTask(c *fiber.Ctx) error
	StartFaceScannerTask(c *fiber.Ctx) error
	DeleteFaceScannerTask(c *fiber.Ctx) error
	CreateFaceScannerTask(c *fiber.Ctx) error
}

// AttachTaskRoutes
//
//	@title						Face Scanner
//	@version					1.0
//	@description				Документация к сервису по распознаванию лиц
//
//
//	@host						localhost:8080
//	@BasePath					/task
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Ключ который выставляется в переменных окружения
func AttachTaskRoutes(router fiber.Router, handlers FaceScanHandlers) {

	router.Post("/extend/:taskUUID", handlers.ExtendFaceScannerTask)
	router.Get("/get/:taskUUID", handlers.GetFaceScannerTask)
	router.Get("/start/:taskUUID", handlers.StartFaceScannerTask)
	router.Delete("/delete/:taskUUID", handlers.DeleteFaceScannerTask)
	router.Post("/create", handlers.CreateFaceScannerTask)

}
