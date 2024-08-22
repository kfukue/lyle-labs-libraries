package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	expectedResponse := `{"error":"Internal Server Error"}`

	RespondWithError(w, code, message)

	if w.Code != code {
		t.Errorf("RespondWithError() returned wrong status code: got %d, want %d", w.Code, code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("RespondWithError() returned wrong content type: got %s, want application/json", w.Header().Get("Content-Type"))
	}

	if w.Body.String() != expectedResponse {
		t.Errorf("RespondWithError() returned unexpected body: got %s, want %s", w.Body.String(), expectedResponse)
	}
}

func TestRespondWithJSON(t *testing.T) {
	w := httptest.NewRecorder()
	code := http.StatusOK
	payload := map[string]interface{}{
		"message": "Success",
	}

	RespondWithJSON(w, code, payload)

	if w.Code != code {
		t.Errorf("RespondWithJSON() returned wrong status code: got %d, want %d", w.Code, code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("RespondWithJSON() returned wrong content type: got %s, want application/json", w.Header().Get("Content-Type"))
	}

	expectedResponse, _ := json.Marshal(payload)
	if w.Body.String() != string(expectedResponse) {
		t.Errorf("RespondWithJSON() returned unexpected body: got %s, want %s", w.Body.String(), string(expectedResponse))
	}
}
