package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

// This is a custom MarshalJSON func. Go will call this method to encode any value which uses this Runtime type in to JSON
func (r Runtime) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d mins", r)

	// By default this string under the hood, won't have a double quote. The double quote we see above is just a syntax thing in go
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}