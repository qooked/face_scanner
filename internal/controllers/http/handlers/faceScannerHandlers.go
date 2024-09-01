package handlers

import (
	"context"
	"errors"
	scannerErrors "faceScanner/internal/errors"
	"faceScanner/internal/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log/slog"
)

type FaceScannerHandlers struct {
	faceScannerUsecase FaceScannerUsecase
}

type FaceScannerUsecase interface {
	ExtendFaceScannerTask(ctx context.Context, task models.ExtendFaceScannerTaskUsecase) (err error)
	GetFaceScannerTask(ctx context.Context, taskUUID string) (task models.GetFaceScannerTaskResponseUsecase, err error)
	StartFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	CreateFaceScannerTask(ctx context.Context, task models.CreateFaceScannerTaskParamsUsecase) (err error)
}

func NewFaceScannerHandlers(
	faceScannerUsecase FaceScannerUsecase,
) *FaceScannerHandlers {
	return &FaceScannerHandlers{
		faceScannerUsecase: faceScannerUsecase,
	}
}

type CreateFaceScannerTaskResponse struct {
	TaskUUID string `json:"task_uuid"`
}

func (h *FaceScannerHandlers) CreateFaceScannerTask(c *fiber.Ctx) error {
	var (
		taskUUID  = uuid.New().String()
		imageData = c.Body()
	)

	if c.Get("Content-Type") != "image/jpeg" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid content type")
	}

	err := h.faceScannerUsecase.CreateFaceScannerTask(
		c.UserContext(),
		models.CreateFaceScannerTaskParamsUsecase{
			Image:    imageData,
			TaskUUID: taskUUID,
		},
	)
	if err != nil {
		if errors.Is(err, scannerErrors.ErrDuplicateTask) {
			return c.Status(fiber.StatusConflict).SendString("Task already created")
		}
		err = fmt.Errorf("h.faceScannerUsecase.ExtendFaceScannerTask(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(CreateFaceScannerTaskResponse{TaskUUID: taskUUID})
}

type ExtendFaceScannerTaskParams struct {
	Image    []byte `json:"image"`
	TaskUUID string
}

func (h *FaceScannerHandlers) ExtendFaceScannerTask(c *fiber.Ctx) error {
	var (
		imageData = c.Body()
	)
	taskUUID := c.Params("taskUUID")
	if taskUUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("taskUUID is empty")
	}
	if c.Get("Content-Type") != "image/jpeg" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid content type")
	}

	err := h.faceScannerUsecase.ExtendFaceScannerTask(
		c.UserContext(),
		models.ExtendFaceScannerTaskUsecase{
			Image:    imageData,
			TaskUUID: taskUUID,
		},
	)
	if err != nil {
		if errors.Is(err, scannerErrors.ErrTaskNotFound) {
			return c.Status(fiber.StatusNoContent).SendString("Task not found")
		}
		if errors.Is(err, scannerErrors.ErrTaskAlreadyStarted) {
			return c.Status(fiber.StatusBadRequest).SendString("Task already started")
		}
		err = fmt.Errorf("h.faceScannerUsecase.ExtendFaceScannerTask(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

type GetFaceScannerTaskResponse struct {
	TaskUUID   string              `json:"taskUUID"`
	Status     int                 `json:"status"`
	ImagesData []SingleTaskPicture `json:"imagesData"`
}

type SingleTaskPicture struct {
	ImageData   []byte `json:"imageData"`
	ApiResponse string `json:"apiResponse"`
}

func (h *FaceScannerHandlers) GetFaceScannerTask(c *fiber.Ctx) error {
	var response GetFaceScannerTaskResponse

	taskUUID := c.Params("taskUUID")
	if taskUUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("taskUUID is empty")
	}
	task, err := h.faceScannerUsecase.GetFaceScannerTask(c.UserContext(), taskUUID)
	if err != nil {
		if errors.Is(err, scannerErrors.ErrTaskNotFound) {
			return c.Status(fiber.StatusNoContent).SendString("Task not found")
		}
		err = fmt.Errorf("h.faceScannerUsecase.GetFaceScannerTask(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	response.TaskUUID = task.TaskUUID
	response.Status = task.Status
	for i := 0; i < len(task.ImagesData); i++ {
		response.ImagesData = append(response.ImagesData, SingleTaskPicture{
			ImageData:   task.ImagesData[i].ImageData,
			ApiResponse: task.ImagesData[i].ApiResponse,
		})
	}

	return c.JSON(response)
}

func (h *FaceScannerHandlers) StartFaceScannerTask(c *fiber.Ctx) error {
	taskUUID := c.Params("taskUUID")
	if taskUUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("taskUUID is empty")
	}

	err := h.faceScannerUsecase.StartFaceScannerTask(c.UserContext(), taskUUID)
	if err != nil {
		if errors.Is(err, scannerErrors.ErrTaskNotFound) {
			return c.Status(fiber.StatusNoContent).SendString("Task not found")
		}

		if errors.Is(err, scannerErrors.ErrTaskAlreadyStarted) {
			return c.Status(fiber.StatusBadRequest).SendString("Task already started")
		}

		err = fmt.Errorf("h.faceScannerUsecase.StartFaceScannerTask(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
func (h *FaceScannerHandlers) DeleteFaceScannerTask(c *fiber.Ctx) error {
	taskUUID := c.Params("taskUUID")
	if taskUUID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("taskUUID is empty")
	}

	err := h.faceScannerUsecase.DeleteFaceScannerTask(c.UserContext(), taskUUID)
	if err != nil {
		if errors.Is(err, scannerErrors.ErrTaskNotFound) {
			return c.Status(fiber.StatusNoContent).SendString("Task not found")
		}

		if errors.Is(err, scannerErrors.ErrTaskAlreadyStarted) {
			return c.Status(fiber.StatusBadRequest).SendString("Task already started")
		}

		err = fmt.Errorf("h.faceScannerUsecase.DeleteFaceScannerTask(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
