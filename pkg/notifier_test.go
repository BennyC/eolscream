package pkg

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSlackNotifier_Notify(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier := NewSlackNotifier(server.URL)
	product := Product{Name: "TestProduct", Version: "1.0.0"}
	releaseInfo := ReleaseInfo{ReleaseDate: "2024-01-01", EndOfLifeDate: "2025-01-01"}

	notifier.Notify(product, releaseInfo)

	// Extend this test to inspect the request payload if necessary
}
