package server

import (
	"net/http"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	
	// Health check
	mux.HandleFunc("/health", HealthHandler)

	return mux
}
