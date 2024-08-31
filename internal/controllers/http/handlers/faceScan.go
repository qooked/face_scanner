package handlers

import (
	"context"
	"faceScanner/internal/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io/ioutil"
	"log/slog"
)

type FaceScannerHandlers struct {
	faceScannerUsecase FaceScannerUsecase
}

type FaceScannerUsecase interface {
	ExtendFaceScannerTask(ctx context.Context, task models.TaskParams) (err error)
	GetFaceScannerTask(ctx context.Context) (task models.TaskResponse, err error)
	StartFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	CreateFaceScannerTask(ctx context.Context, task models.TaskParams) (err error)
}

func NewFaceScannerHandlers(
	faceScannerUsecase FaceScannerUsecase,
) *FaceScannerHandlers {
	return &FaceScannerHandlers{
		faceScannerUsecase: faceScannerUsecase,
	}
}

func (h *FaceScannerHandlers) CreateFaceScannerTask(c *fiber.Ctx) error {
	return nil
}
func (h *FaceScannerHandlers) ExtendFaceScannerTask(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		err = fmt.Errorf("c.FormFile(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusBadRequest)
	}

	fileContent, err := file.Open()
	if err != nil {
		err = fmt.Errorf("file.Open(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusBadRequest)
	}

	fileBytes, err := ioutil.ReadAll(fileContent)
	if err != nil {
		err = fmt.Errorf("ioutil.ReadAll(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = h.faceScannerUsecase.ExtendFaceScannerTask(
		c.UserContext(),
		models.TaskParams{
			Image:    fileBytes,
			TaskUUID: uuid.New().String(),
		},
	)
	if err != nil {
		err = fmt.Errorf("h.faceScannerUsecase.ExtendFaceScannerTask(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *FaceScannerHandlers) GetFaceScannerTask(c *fiber.Ctx) error {
	return nil
}
func (h *FaceScannerHandlers) StartFaceScannerTask(c *fiber.Ctx) error {
	return nil
}
func (h *FaceScannerHandlers) DeleteFaceScannerTask(c *fiber.Ctx) error {
	return nil
}
