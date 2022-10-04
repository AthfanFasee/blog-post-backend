package data

import (
	"time"

	"github.com/AthfanFasee/blog-post-backend/internal/dto"
)

var mockComment = &Comment{
	ID:        1,
	CreatedAt: time.Now(),
	Text:      "Mocked Comment",
	CreatedBy: 1,
	PostID:    1,
}

var mockCommentResponseBody = &dto.CommentResponseBody{
	ID:        mockComment.ID,
	Text:      mockComment.Text,
	CreatedBy: mockComment.CreatedBy,
	PostID:    1,
}

type MockCommentModel struct{}

func (c MockCommentModel) GetAllForPost(postID int64) ([]*dto.CommentResponseBody, error) {
	switch postID {
	case 1:
		return []*dto.CommentResponseBody{mockCommentResponseBody}, nil
	default:
		return nil, ErrRecordNotFound
	}
}

func (c MockCommentModel) Insert(comment *Comment) error {
	return nil
}

func (c MockCommentModel) Delete(id int64) error {
	switch id {
	case 1:
		return nil
	default:
		return ErrRecordNotFound
	}
}
