package utils

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func AddChainingMiddleware(h http.Handler, ms ...Middleware) http.Handler {
	if len(ms) == 0 {
		return h
	}

	wrapperMiddleware := h
	for i := len(ms) - 1; i >= 0; i-- {
		wrapperMiddleware = ms[i](wrapperMiddleware)
	}

	return wrapperMiddleware
}
