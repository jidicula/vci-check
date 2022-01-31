package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckHandler(t *testing.T) {

	tests := map[string]struct {
		issURL       string
		wantCode     int
		wantResponse string
	}{
		"trueIssuer": {
			issURL:       "https://myvaccinerecord.cdph.ca.gov/creds",
			wantCode:     http.StatusOK,
			wantResponse: `{"message": true}`,
		},
		"falseIssuer": {
			issURL:       "https://mallory.me/creds",
			wantCode:     http.StatusOK,
			wantResponse: `{"message": false}`,
		},
		"emptyIssuer": {
			issURL:       "",
			wantCode:     http.StatusBadRequest,
			wantResponse: `{"message": "No issuer URL provided"}` + "\n",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			query := fmt.Sprintf("/?iss=%s", tt.issURL)

			req, err := http.NewRequest("GET", query, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(checkHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
			got := rr.Body.String()

			if got != tt.wantResponse {
				t.Errorf("%s returned wrong response: got %s, want %s", name, got, tt.wantResponse)
			}
		})
	}
	name := "incorrect GET"
	t.Run(name, func(t *testing.T) {
		wantCode := http.StatusMethodNotAllowed
		wantResponse := `{"message": "expect method GET at /?iss=<url>"}` + "\n"
		reader := strings.NewReader(`{"foo": "bar"}`)
		req, err := http.NewRequest("POST", "/", reader)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(checkHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != wantCode {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		got := rr.Body.String()
		if got != wantResponse {
			t.Errorf("%s returned wrong response: got %s, want %s", name, got, wantResponse)
		}
	})
}
