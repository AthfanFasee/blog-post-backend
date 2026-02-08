package dto

import (
	"errors"
	"testing"

	"github.com/AthfanFasee/blog-post-backend/internal/assert"
)

func TestReadTime_MarshalJSON(t *testing.T) {
	r := ReadTime(120)
	jsonValue, err := r.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedJSONValueST1003default := "\"120 mins\""
	assert.Equal(t, string(jsonValue), expectedJSONValueST1003default)
}

func TestReadTime_UnmarshalJSON(t *testing.T) {
	jsonValue := "\"120 mins\""
	var r ReadTime
	err := r.UnmarshalJSON([]byte(jsonValue))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.Equal(t, r, ReadTime(120))

	invalidJSONValueST1003default := "\"invalid format\""
	err = r.UnmarshalJSON([]byte(invalidJSONValueST1003default))
	if !errors.Is(err, ErrInvalidReadtimeFormat) {
		t.Fatalf("Expected error for invalid format")
	}
}
