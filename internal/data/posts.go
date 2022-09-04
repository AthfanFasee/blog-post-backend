package data

import (
	"time"
)

type Movie struct {
ID int64 
CreatedAt time.Time 
Title string 
Year int32 
Runtime Runtime // If we use our custom Runtime type here (which has the underlying type int32) go will use Runtime type's method MarshalJSON to encode this to JSON and it will be encoded to Runtime type (a string in the format "<runtime> mins") instead of int
Genres []string 
Version int32 
}