package dto

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type ReadTime int32

// This is a custom MarshalJSON func. Go will call this method to encode any value which got Runtime type in to JSON
func (r ReadTime) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d mins", r)

	// A JSON string must be wrapped in double quotes
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

// This is a custom UnmarshalJSON func. Go will call this method to decode any JSON value which got Runtime type in the destination
func (r *ReadTime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Convert the int32 type to Runtime type, deference the receiver, and assign it to the underlying value of r
	*r = ReadTime(i)

	return nil
}
