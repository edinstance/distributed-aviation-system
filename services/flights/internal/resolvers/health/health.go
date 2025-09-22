package health

import (
	"net/http"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
)

// HealthHandler responds to health check requests with HTTP 200 and a JSON body.
// It sets the Content-Type header to "application/json" and writes `{"status":"UP"}` to the response.
// Any error returned when writing the body is ignored.
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	logger.Debug("Healthcheck called", "path", "/health")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"UP"}`))
}
