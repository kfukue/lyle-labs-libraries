package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Middleware(handler)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}

	for key, value := range expectedHeaders {
		if rr.Header().Get(key) != value {
			t.Errorf("Expected header %s: %s, but got %s", key, value, rr.Header().Get(key))
		}
	}
}
