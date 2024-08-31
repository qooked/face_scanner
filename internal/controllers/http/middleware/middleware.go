package middleware

type HTTPMiddleware struct {
	authorizationKey string
}

func NewHttpMiddleware(
	authorizationKey string,
) *HTTPMiddleware {
	return &HTTPMiddleware{
		authorizationKey: authorizationKey,
	}
}
