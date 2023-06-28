package jsonlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/AthfanFasee/blog-post-backend/internal/assert"
)

func TestPrintInfo(t *testing.T) {
	// Using buffer as io.Writer for testing purpose.
	buffer := new(bytes.Buffer)
	logger := New(buffer, LevelInfo)

	message := "test message"
	properties := map[string]string{
		"key": "value",
	}

	logger.PrintInfo(message, properties)

	// Convert the content of buffer to a string.
	logged := buffer.String()
	assert.StringContains(t, logged, message)

	// Unmarshal logged JSON.
	var logOutput map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &logOutput)
	if err != nil {
		t.Errorf("Failed to unmarshal log output: %v", err)
	}

	assertedMessage, ok := logOutput["message"].(string)
	if !ok {
		t.Errorf("logOutput[\"message\"] is not of type string")
	}
	assert.Equal(t, assertedMessage, message)

	level, ok := logOutput["level"].(string)
	if !ok {
		t.Errorf("logOutput[\"level\"] is not of type string")
	}
	assert.Equal(t, level, "INFO")

	assertedKey, ok := logOutput["properties"].(map[string]interface{})["key"].(string)
	if !ok {
		t.Errorf("logOutput[\"properties\"][\"key\"] is not of type string")
	}
	assert.Equal(t, assertedKey, "value")
}

func TestPrintError(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := New(buffer, LevelError)

	errMessage := "error message"
	properties := map[string]string{
		"key": "value",
	}

	logger.PrintError(errors.New(errMessage), properties)

	logged := buffer.String()
	assert.StringContains(t, logged, errMessage)

	var logOutput map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &logOutput)
	if err != nil {
		t.Errorf("Failed to unmarshal log output: %v", err)
	}

	assertedMessage, ok := logOutput["message"].(string)
	if !ok {
		t.Errorf("logOutput[\"message\"] is not of type string")
	}
	assert.Equal(t, assertedMessage, errMessage)

	level, ok := logOutput["level"].(string)
	if !ok {
		t.Errorf("logOutput[\"level\"] is not of type string")
	}
	assert.Equal(t, level, "ERROR")

	assertedKey, ok := logOutput["properties"].(map[string]interface{})["key"].(string)
	if !ok {
		t.Errorf("logOutput[\"properties\"][\"key\"] is not of type string")
	}
	assert.Equal(t, assertedKey, "value")

	assert.StringContains(t, logOutput["trace"].(string), "goroutine") // Ensure stack trace exists.
}
