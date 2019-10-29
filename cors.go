package main

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const allowMethods = "POST, GET, OPTIONS"

const (
	internalErrorMessage = "Internal Server Error"
	internalErrorCode    = "500"
	badRequestErrorCode  = "400"
)

var allowedHeaders = []string{"Origin", "Accept", "Content-Type", "Authorization", "Connection", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Accept-Encoding"}

// See https://github.com/rs/cors for more details
// Allows all methods and all headers
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if r.Method == "OPTIONS" {
			if origin == "" {
				w.WriteHeader(http.StatusBadRequest)
				log.Error("Origin needed for preflight")
			}
			headers := w.Header()
			headers.Add("Access-Control-Allow-Origin", origin)
			headers.Add("Vary", "Origin")
			headers.Add("Vary", "Access-Control-Request-Method")
			headers.Add("Vary", "Access-Control-Request-Headers")
			headers.Set("Access-Control-Allow-Methods", allowMethods)
			headers.Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

			w.WriteHeader(http.StatusOK)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(allowedHeaders, ", "))
			next.ServeHTTP(w, r)
		}
	})
}
