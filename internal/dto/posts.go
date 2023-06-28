package dto

import (
	"time"
)

type CreatePostRequestBody struct {
	Title    string   `json:"title"`
	PostText string   `json:"postText"`
	ReadTime ReadTime `json:"readTime"`
	Img      string   `json:"img"`
}

// Define the input struct in a way, all the field got zero value 'nil'.
type UpdatePostRequestBody struct {
	Title    *string   `json:"title"`
	PostText *string   `json:"postText"`
	ReadTime *ReadTime `json:"readTime"`
	Img      *string   `json:"img"`
}

type PostResponseBody struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Title     string    `json:"title"`
	PostText  string    `json:"postText"`
	Img       string    `json:"img"`
	ReadTime  ReadTime  `json:"readTime"` // If we use our custom ReadTime type here (which has the underlying type int32) go will use ReadTime type's method MarshalJSON to encode this to JSON and it will be encoded to ReadTime type (a string in the format "<readtime> mins") instead of int.
	LikedBy   []int64   `json:"likedBy,omitempty"`
	CreatedBy int64     `json:"createdBy"`
	UserName  string    `json:"userName"`
}
