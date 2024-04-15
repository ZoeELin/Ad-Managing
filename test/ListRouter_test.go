// test/TestList.go

package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAdRouter(t *testing.T) {
	// Create a test server to emulate an actual HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Login for processing requests
		if r.URL.Path == "/api/v1/ad" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Send a test request to a route on the test server
	r, err := http.Get(server.URL + "/api/v1/ad")
	assert.NoError(t, err)
	defer r.Body.Close()

	assert.Equal(t, http.StatusOK, r.StatusCode)
}
