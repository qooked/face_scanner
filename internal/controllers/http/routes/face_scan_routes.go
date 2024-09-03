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
// @description     Authorization - Ключ, который можно получить при регистрации, basic base64(login:password), требуется для всех запросов группы task
// @description Статусы заданий
// @description 	Новое задание   1
// @description Задание в обработке 2
// @description Успешное задание    3
// @description Частично успешное задание 4
// @description Неуспешное задание  5
// @host            localhost:8080
// @BasePath        /task
// @Security        ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in              header
// @name            Authorization
func AttachTaskRoutes(router fiber.Router, middleware Middleware, handlers FaceScanHandlers) {
	router.Use(middleware.AuthorizationMiddleware())
	router.Post("/extend/:taskUUID", handlers.ExtendFaceScannerTask)
	router.Get("/get/:taskUUID", handlers.GetFaceScannerTask)
	router.Get("/start/:taskUUID", handlers.StartFaceScannerTask)
	router.Delete("/delete/:taskUUID", handlers.DeleteFaceScannerTask)
	router.Post("/create", handlers.CreateFaceScannerTask)

}
