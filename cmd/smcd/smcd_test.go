package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/cebarks/smcd"
	"github.com/stretchr/testify/assert"
)

func init() {
	var err error
	smcd.WorkingDir, err = filepath.Abs("test/")
	if err != nil {
		log.Fatalf("error setting up test dir: %v", err)
	}
}

func TestPingRoute(t *testing.T) {
	router := setupRouter(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
