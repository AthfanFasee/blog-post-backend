package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/AthfanFasee/blog-post-backend/internal/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/api/v1/healthcheck")

	assert.Equal(t, code, http.StatusOK)

	// Unmarshal the JSON and check the actual values. This way,
	// tests won't break just because the formatting of the JSON changes (writeJSON helper adds whitespaces and newlines to the output by default).
	var data envelope
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}

	if status, ok := data["status"].(string); ok {
		assert.Equal(t, status, "available")
	} else {
		t.Error("status is not a string")
	}
	// data["systemInfo"] is originally unmarshalled from JSON into a map[string]interface{}.
	// JSON objects unmarshal into map[string]interface{} by default.
	systemInfo, ok := data["systemInfo"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'systemInfo' to be in the format map[string]interface{}")
	}

	_, ok = systemInfo["environment"]
	assert.Equal(t, ok, true)

	_, ok = systemInfo["version"]
	assert.Equal(t, ok, true)

	_, ok = systemInfo["build_time"]
	assert.Equal(t, ok, true)

}
