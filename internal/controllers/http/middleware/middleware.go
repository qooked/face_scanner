package middleware

import "context"

type HTTPMiddleware struct {
	authUsecase AuthUsecase
}

type AuthUsecase interface {
	GetUserCredentials(ctx context.Context, email string) (hashedPassword string, err error)
}

func NewHttpMiddleware(
	authUsecase AuthUsecase,
) *HTTPMiddleware {
	return &HTTPMiddleware{
		authUsecase: authUsecase,
	}
}
