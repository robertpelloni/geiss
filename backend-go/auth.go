package main

import (
	"net/http"
	"os"
)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("JULES_API_KEY")
		// Fail-closed authentication
		if apiKey == "" {
			http.Error(w, "Server configuration error: JULES_API_KEY is not set", http.StatusInternalServerError)
			return
		}

		reqKey := r.Header.Get("X-API-KEY")
		if reqKey != apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
