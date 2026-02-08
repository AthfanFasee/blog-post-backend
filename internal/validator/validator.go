// Package validator provides utilities for validating common data types
// such as email format and uniqueness within a slice. It includes a Validator
// struct to accumulate error messages and provide methods to check validation
// results. It also offers helper functions for pattern matching and checking
// membership within a list.
package validator

import (
	"regexp"
	"slices"
)

// Regular expression for sanity checking the format of email addresses.
var (
	EmailRX = regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	)
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map as long as no entry already exists for the given key.
func (v *Validator) AddError(key string, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key string, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Matches returns true if a string matches a specific regexp pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique returns true if all the strings in a slice are unique.
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	// When we are adding each value of slice to map, if a duplicate value comes, it will be replaced.
	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

// In returns true if a value is in a list of strings.
func In(value string, list ...string) bool {
	return slices.Contains(list, value)
}
