package health

import (
	"net/http"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	logger.Debug("Healthcheck called", "path", "/health")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
