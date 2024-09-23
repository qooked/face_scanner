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

type CreateFaceScannerTaskParams struct {
	Image []byte `json:"image" validate:"required" swaggerType:"image"`
}

type CreateFaceScannerTaskResponse struct {
	TaskUUID  string `json:"task_uuid"`
	ImageUUID string `json:"image_uuid"`
}

// CreateFaceScannerTask godoc
// @Summary         Создание задания для распознавания лиц.
// @Tags            api
// @Accept          jpeg
// @Produce         json
// @Security        ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth
// @Param			Authorization	header	string	true	"Ключ, который можно получить при регистрации, basic base64(login:password)."
// @Param           request body CreateFaceScannerTaskParams true "Загружаемое jpeg изображение"
// @Success         200 {object} CreateFaceScannerTaskResponse
// @Failure         400 {string} string "Invalid request"
// @Failure         500 {string} string "Internal Server Error"
// @Router          /create [post]
func (h *FaceScannerHandlers) CreateFaceScannerTask(c *fiber.Ctx) error {
	var (
		taskUUID  = uuid.New().String()
		imageUUID = uuid.New().String()
		imageData = c.Body()
		response  CreateFaceScannerTaskResponse
	)

	if c.Get(fiber.HeaderContentType) != "image/jpeg" {
		return c.Status(fiber.StatusBadRequest).
			SendString("Invalid content type")
	}

	err := h.faceScannerUsecase.CreateFaceScannerTask(
		c.UserContext(),
		models.CreateFaceScannerTaskParamsUsecase{
			Image:     imageData,
			TaskUUID:  taskUUID,
			ImageUUID: imageUUID,
		},
	)
	if err != nil {
		if errors.Is(err, scannerErrors.ErrDuplicateTask) {
			return c.Status(fiber.StatusConflict).
				SendString("Task already created")
		}
		err = fmt.Errorf("h.faceScannerUsecase.CreateFaceScannerTask(...): %w", err)
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	response.ImageUUID = imageUUID
	response.TaskUUID = taskUUID

	return c.Status(fiber.StatusOK).
		JSON(response)
}

type ExtendFaceScannerTaskParams struct {
	Image    []byte `json:"image" validate:"required" swaggerType:"image"`
	TaskUUID string `json:"-"`
}

// ExtendFaceScannerTask godoc
// @Summary         Добавление файлов в задание.
// @Tags            api
// @Accept          jpeg
// @Produce         jpeg
// @Security        ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth
// @Param			Authorization	header	string	true	"Ключ, который можно получить при регистрации, basic base64(login:password)."
// @Param           taskUUID path string true "UUID задания, которое нужно расширить"
// @Param           request body ExtendFaceScannerTaskParams true "Загружаемое jpeg изображение"
// @Success         200
// @Failure         400 {string} string "Bad request"
// @Failure         500 {string} string "Internal Server Error"
// @Router          /extend/{taskUUID} [post]
func (h *FaceScannerHandlers) ExtendFaceScannerTask(c *fiber.Ctx) error {
	var (
		imageData = c.Body()
		imageUUID = uuid.New().String()
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
			Image:     imageData,
			TaskUUID:  taskUUID,
			ImageUUID: imageUUID,
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
	Stats      Stats               `json:"stats"`
}

type SingleTaskPicture struct {
	ImageUUID string `json:"imageUUID"`
	FileName  string `json:"fileName"`
	Faces     []Face `json:"face"`
}

type Stats struct {
	FacesCount       int `json:"facesCount"`
	MaleFemaleCount  int `json:"maleFemaleCount"`
	AverageMaleAge   int `json:"averageMaleAge"`
	AverageFemaleAge int `json:"averageFemaleAge"`
}
type Face struct {
	BoundingBox `json:"boundingBox"`
	Sex         string  `json:"sex"`
	Age         float64 `json:"age"`
}

type BoundingBox struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// GetFaceScannerTask godoc
// @Summary         Получение данных по заданию.
// @Tags            api
// @Accept          json
// @Produce         json
// @Security        ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth
// @Param			Authorization	header	string	true	"Ключ, который можно получить при регистрации, basic base64(login:password)."
// @Param           taskUUID path string true "UUID задания"
// @Success         200 {object} GetFaceScannerTaskResponse
// @Failure         400 {string} string "Bad request"
// @Failure         500 {string} string "Internal Server Error"
// @Router          /get/{taskUUID} [get]
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
	response.Stats.FacesCount = task.FacesCount
	response.Stats.MaleFemaleCount = task.MaleFemaleCount
	response.Stats.AverageFemaleAge = int(task.AverageFemaleAge)
	response.Stats.AverageMaleAge = int(task.AverageMaleAge)
	for _, image := range task.ImagesData {

		var faces []Face
		for _, face := range image.Faces {
			singleBoundingBox := BoundingBox{
				X: face.BoundingBox.X,
				Y: face.BoundingBox.Y,
				W: face.BoundingBox.W,
				H: face.BoundingBox.H,
			}

			faces = append(faces, Face{
				BoundingBox: singleBoundingBox,
				Sex:         face.Sex,
				Age:         face.Age,
			})
		}

		response.ImagesData = append(response.ImagesData, SingleTaskPicture{
			ImageUUID: image.ImageUUID,
			Faces:     faces,
			FileName:  image.FileName,
		})
	}

	return c.JSON(response)
}

// StartFaceScannerTask godoc
// @Summary         Запуск задания на распознавание лиц.
// @Tags            api
// @Accept          json
// @Produce         json
// @Security        ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth
// @Param           taskUUID path string true "UUID задания"
// @Param			Authorization	header	string	true	"Ключ, который можно получить при регистрации, basic base64(login:password)."
// @Success         200
// @Failure         400 {string} string "Bad request"
// @Failure         500 {string} string "Internal Server Error"
// @Router          /start/{taskUUID} [post]
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

// DeleteFaceScannerTask godoc
// @Summary         Удаление задания на распознавание лиц.
// @Tags            api
// @Accept          json
// @Produce         json
// @Security        ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth
// @Param			Authorization	header	string	true	"Ключ, который можно получить при регистрации, basic base64(login:password)."
// @Param           taskUUID path string true "UUID задания"
// @Success         200
// @Failure         400 {string} string "Bad request"
// @Failure         500 {string} string "Internal Server Error"
// @Router          /delete/{taskUUID} [delete]
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
