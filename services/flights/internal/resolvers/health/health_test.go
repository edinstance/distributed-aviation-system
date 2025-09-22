package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(testHelper *testing.T) {
	request, err := http.NewRequest("GET", "/health", nil)

	if err != nil {
		testHelper.Fatal(err)
	}

	responseWriter := httptest.NewRecorder()

	HealthHandler(responseWriter, request)

	if responseWriter.Code != http.StatusOK {
		testHelper.Fatalf("Health handler returned wrong status code: got %v want %v", responseWriter.Code, http.StatusOK)
	}

	expectedBody := `{"status":"UP"}`
	if responseWriter.Body.String() != expectedBody {
		testHelper.Errorf("handler returned unexpected body: got %v want %v",
			responseWriter.Body.String(), expectedBody)
	}

	expectedContentType := "application/json"
	if responseWriter.Header().Get("Content-Type") != expectedContentType {
		testHelper.Errorf("handler returned unexpected content type: got %v want %v",
			responseWriter.Header().Get("Content-Type"), expectedContentType)
	}
}
