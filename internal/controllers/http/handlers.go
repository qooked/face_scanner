package http

import (
	"context"
	"faceScanner/internal/controllers/http/handlers"
	"faceScanner/internal/controllers/http/middleware"
	"faceScanner/internal/controllers/http/routes"
	"faceScanner/internal/models"
)

type Usecase interface {
	ExtendFaceScannerTask(ctx context.Context, task models.TaskParams) (err error)
	GetFaceScannerTask(ctx context.Context) (task models.TaskResponse, err error)
	StartFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	CreateFaceScannerTask(ctx context.Context, task models.TaskParams) (err error)
}

func (s *Server) AttachHandlers(ctx context.Context, taskUsecase Usecase) {
	middleware := middleware.NewHttpMiddleware(
		s.options.AuthorizationKey,
	)

	s.app.Use(middleware.AuthorizationMiddleware())

	taskGroup := s.app.Group("/task")
	taskHandlers := handlers.NewFaceScannerHandlers(
		taskUsecase)

	routes.AttachCascadeRoutes(taskGroup, taskHandlers)
}
