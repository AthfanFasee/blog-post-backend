package dto

type CommentRequestBody struct {
	Text   string `json:"text"`
	PostID int64  `json:"post"`
}

type CommentResponseBody struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	CreatedBy int64  `json:"createdBy"`
	PostID    int64  `json:"post"`
	UserName  string `json:"userName"`
}
