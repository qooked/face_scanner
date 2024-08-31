package http

import (
	"context"
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
