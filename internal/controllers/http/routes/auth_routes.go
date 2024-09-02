package routes

import "github.com/gofiber/fiber/v2"

type AuthHandlers interface {
	Register(c *fiber.Ctx) error
}

// AttachAuthRoutes
//
// @title           Face Scanner
// @version         1.0
// @description     Документация к сервису по распознаванию лиц
// @host            localhost:8080
// @BasePath        /auth
// @Security        ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in              header
// @name            Authorization
// @description     Ключ, который можно получить при регистрации, basic base64(login:password)
func AttachAuthRoutes(router fiber.Router, handlers AuthHandlers) {
	router.Post("/register", handlers.Register)

}
