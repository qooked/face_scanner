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
type Middleware interface {
	AuthorizationMiddleware() fiber.Handler
}

// AttachTaskRoutes
//
// @title           Face Scanner
// @version         1.0
// @description     Документация к сервису по распознаванию лиц
// @host            localhost:8080
// @BasePath        /task
// @Security        ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in              header
// @name            Authorization
// @description     Ключ, который можно получить при регистрации, basic base64(login:password)
func AttachTaskRoutes(router fiber.Router, middleware Middleware, handlers FaceScanHandlers) {
	router.Use(middleware.AuthorizationMiddleware())
	router.Post("/extend/:taskUUID", handlers.ExtendFaceScannerTask)
	router.Get("/get/:taskUUID", handlers.GetFaceScannerTask)
	router.Get("/start/:taskUUID", handlers.StartFaceScannerTask)
	router.Delete("/delete/:taskUUID", handlers.DeleteFaceScannerTask)
	router.Post("/create", handlers.CreateFaceScannerTask)

}
