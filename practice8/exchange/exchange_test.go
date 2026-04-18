package exchange

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRate(t *testing.T) {
	t.Run("Successful scenario", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := RateResponse{Base: "USD", Target: "EUR", Rate: 0.92}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		rate, err := service.GetRate("USD", "EUR")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if rate != 0.92 {
			t.Errorf("expected 0.92, got %f", rate)
		}
	})

	t.Run("API Business Error - 404 with error message", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid currency pair"})
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "XXX")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Malformed JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"base": "USD", "target":`))
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")

		if err == nil {
			t.Error("expected decode error, got nil")
		}
	})

	t.Run("Slow Response/Timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(6 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		service.Client.Timeout = 2 * time.Second

		_, err := service.GetRate("USD", "EUR")
		if err == nil {
			t.Error("expected timeout error, got nil")
		}
	})

	t.Run("Server Panic / 500 Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("something went wrong")
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Empty Body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")

		if err == nil {
			t.Error("expected decode error, got nil")
		}
	})
}