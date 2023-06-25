package dto

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidReadtimeFormat = errors.New("invalid readtime format")

type ReadTime int32

// This is a custom MarshalJSON func. Go will call this method to encode any value which got Readtime type into JSON.
func (r ReadTime) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d mins", r)

	// A JSON string must be wrapped in double quotes.
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

// This is a custom UnmarshalJSON func. Go will call this method to decode any JSON value which got Readtime type in the destination.
func (r *ReadTime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidReadtimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidReadtimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidReadtimeFormat
	}

	// Convert the int32 type to Readtime type, deference the receiver, and assign it to the underlying value of r.
	*r = ReadTime(i)

	return nil
}
