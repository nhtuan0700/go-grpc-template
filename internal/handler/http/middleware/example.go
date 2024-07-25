package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://www.alexedwards.net/blog/making-and-using-middleware
func ExampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		data, err := io.ReadAll(r.Body)
		if err == nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(data))
		fmt.Println(string(data))
		next.ServeHTTP(w, r)
	})
}
