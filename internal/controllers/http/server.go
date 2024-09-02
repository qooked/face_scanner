package http

import (
	"context"
	"faceScanner/internal/controllers/http/handlers"
	"faceScanner/internal/controllers/http/middleware"
	"faceScanner/internal/controllers/http/routes"
	"faceScanner/internal/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type ServerOptions struct {
	FiberAppParams   FiberAppParams
	AuthorizationKey string
}

type FiberAppParams struct {
	Host string
	Port string
}

type Server struct {
	app     *fiber.App
	options ServerOptions
}

type Usecase interface {
	ExtendFaceScannerTask(ctx context.Context, task models.ExtendFaceScannerTaskUsecase) (err error)
	GetFaceScannerTask(ctx context.Context, taskUUID string) (task models.GetFaceScannerTaskResponseUsecase, err error)
	StartFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	CreateFaceScannerTask(ctx context.Context, task models.CreateFaceScannerTaskParamsUsecase) (err error)
}

func NewServer(
	host string,
	port string,
	authorizationKey string,
) *Server {
	return &Server{
		app: fiber.New(),
		options: ServerOptions{
			FiberAppParams: FiberAppParams{
				Host: host,
				Port: port,
			},
			AuthorizationKey: authorizationKey,
		},
	}
}

func (s *Server) Run() {
	err := s.app.Listen(
		fmt.Sprintf("%s:%s", s.options.FiberAppParams.Host, s.options.FiberAppParams.Port),
	)
	if err != nil {
		err = fmt.Errorf("srv.listener.Listen(...): %w", err)
		panic(err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.app.ShutdownWithContext(ctx)
	if err != nil {
		err = fmt.Errorf("srv.listener.ShutdownWithContext(...): %w", err)

		return err
	}

	return nil
}

func (s *Server) AttachHandlers(ctx context.Context, taskUsecase Usecase) {
	middleware := middleware.NewHttpMiddleware(
		s.options.AuthorizationKey,
	)

	s.app.Use(middleware.AuthorizationMiddleware())

	taskGroup := s.app.Group("/task")
	taskHandlers := handlers.NewFaceScannerHandlers(
		taskUsecase)

	routes.AttachTaskRoutes(taskGroup, taskHandlers)
}
