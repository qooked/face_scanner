package handlers

import (
	"context"
	"faceScanner/pkg/utills"
	"faceScanner/pkg/validator"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type AuthHandlers struct {
	authUsecase AuthUsecase
}

type AuthUsecase interface {
	SaveUserCredentials(ctx context.Context, email string, unhashedPassword string) (err error)
}

func NewAuthHandlers(
	AuthUsecase AuthUsecase,
) *AuthHandlers {
	return &AuthHandlers{
		authUsecase: AuthUsecase,
	}
}

type RegisterParams struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	AuthorizationToken string `json:"authorization_token"`
}

// Register godoc
// @Summary         Регистрация нового пользователя и получение токена авторизации.
// @Tags            api
// @Accept          json
// @Produce         json
// @Security        ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth
// @description     Регистрация нового пользователя и получение токена авторизации, используется для всех других запросов
// @Param           RegisterParams body RegisterParams true "Логин и пароль пользователя"
// @Success         200 {object} RegisterResponse "Успешный ответ с токеном авторизации"
// @Failure         400 {string} string "Неверный запрос, например, если логин не является email"
// @Failure         500 {string} string "Внутренняя ошибка сервера"
// @Router          /register [post]
func (a *AuthHandlers) Register(c *fiber.Ctx) error {
	var (
		params   RegisterParams
		response RegisterResponse
	)
	if err := utills.UnmarshalFromJSON(c, &params); err != nil {
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if valid := validator.ValidateEmail(params.Login); !valid {
		return c.Status(fiber.StatusBadRequest).
			SendString("Login must be an email")
	}

	err := a.authUsecase.SaveUserCredentials(c.Context(), params.Login, params.Password)
	if err != nil {
		err = fmt.Errorf("a.authUsecase.SaveUserCredentials(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	response.AuthorizationToken = utills.CreateBasicAuthToken(params.Login, params.Password)

	return c.Status(fiber.StatusOK).JSON(response)
}
