package validator

import (
	"testing"

	"github.com/AthfanFasee/blog-post-backend/internal/assert"
)

func TestMatches(t *testing.T) {
	t.Run("valid email", func(t *testing.T) {
		email := "test@example.com"
		got := Matches(email, EmailRX)
		assert.Equal(t, got, true)
	})

	t.Run("invalid email", func(t *testing.T) {
		email := "invalid-email"
		got := Matches(email, EmailRX)
		assert.Equal(t, got, false)
	})
}

func TestUnique(t *testing.T) {
	t.Run("unique values", func(t *testing.T) {
		values := []string{"a", "b", "c", "d", "e"}
		got := Unique(values)
		assert.Equal(t, got, true)
	})

	t.Run("duplicate values", func(t *testing.T) {
		values := []string{"a", "b", "c", "a", "b"}
		got := Unique(values)
		assert.Equal(t, got, false)
	})
}

func TestIn(t *testing.T) {
	t.Run("value in list", func(t *testing.T) {
		value := "b"
		list := []string{"a", "b", "c", "d", "e"}
		got := In(value, list...)
		assert.Equal(t, got, true)
	})

	t.Run("value not in list", func(t *testing.T) {
		value := "z"
		list := []string{"a", "b", "c", "d", "e"}
		got := In(value, list...)
		assert.Equal(t, got, false)
	})
}

func TestValidator(t *testing.T) {
	t.Run("valid Validator", func(t *testing.T) {
		v := New()
		v.Check(true, "key", "message")
		got := v.Valid()
		assert.Equal(t, got, true)
	})

	t.Run("invalid Validator", func(t *testing.T) {
		v := New()
		v.Check(false, "key", "message")
		got := v.Valid()
		assert.Equal(t, got, false)
	})

	t.Run("add error", func(t *testing.T) {
		v := New()
		v.AddError("key", "message")
		got := v.Valid()
		assert.Equal(t, got, false)
	})

	t.Run("multiple errors", func(t *testing.T) {
		v := New()
		v.AddError("key1", "message1")
		v.AddError("key2", "message2")
		errCount := len(v.Errors)
		assert.Equal(t, errCount, 2)
	})

	t.Run("duplicate error keys", func(t *testing.T) {
		v := New()
		v.AddError("key", "message1")
		v.AddError("key", "message2")
		errCount := len(v.Errors)
		assert.Equal(t, errCount, 1)
	})
}
