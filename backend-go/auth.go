package main

import (
	"net/http"
	"os"
)

// authMiddleware fails closed: when JULES_API_KEY is unset, protected
// routes return 401 instead of silently allowing every caller.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("JULES_API_KEY")
		if apiKey == "" {
			http.Error(w, "Unauthorized: JULES_API_KEY not configured", http.StatusUnauthorized)
			return
		}
		if r.Header.Get("X-API-KEY") != apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// authMutatingOnly enforces auth only for non-GET/HEAD methods, so
// read-only GETs on mixed handlers (prune scan, backup list) stay open.
func authMutatingOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			next(w, r)
			return
		}
		authMiddleware(next)(w, r)
	}
}
