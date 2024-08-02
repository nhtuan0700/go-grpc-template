package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/nhtuan0700/go-grpc-template/internal/config/constants"
	"github.com/nhtuan0700/go-grpc-template/internal/utils"
)

func RequestMiddlewareWith(ctx context.Context) utils.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.New()
			ctx = context.WithValue(ctx, constants.RequestIDKey, requestID.String())
			next.ServeHTTP(w, r)
		})
	}
}
