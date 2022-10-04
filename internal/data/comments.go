package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/AthfanFasee/blog-post-backend/internal/dto"
	"github.com/AthfanFasee/blog-post-backend/internal/validator"
)

type Comment struct {
	ID        int64
	CreatedAt time.Time
	Text      string
	CreatedBy int64
	PostID    int64
}

type CommentModel struct {
	DB *sql.DB
}

func (c CommentModel) GetAllForPost(postID int64) ([]*dto.CommentResponseBody, error) {
	query := `SELECT c.id, c.text, c.created_by, c.post_id, u.name
	FROM comments c
	INNER JOIN users u ON c.created_by = u.id
	WHERE post_id = $1
	ORDER BY id DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*dto.CommentResponseBody{}

	for rows.Next() {
		var comment Comment
		var userName string

		err := rows.Scan(
			&comment.ID,
			&comment.Text,
			&comment.CreatedBy,
			&comment.PostID,
			&userName,
		)

		if err != nil {
			return nil, err
		}

		CommentResponseBody := dto.CommentResponseBody{
			ID:        comment.ID,
			Text:      comment.Text,
			CreatedBy: comment.CreatedBy,
			PostID:    comment.PostID,
			UserName:  userName,
		}

		comments = append(comments, &CommentResponseBody)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c CommentModel) Insert(comment *Comment) error {
	query := `
	INSERT INTO comments (text, post_id, created_by)
	VALUES ($1, $2, $3)
	RETURNING id`

	args := []interface{}{comment.Text, comment.PostID, comment.CreatedBy}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.ID)
}

func (c CommentModel) Delete(id int64) error {
	query := `
	DELETE FROM comments
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	v.Check(comment.Text != "", "text", "Comment cannot be empty")
	v.Check(len(comment.Text) <= 200, "text", "Comment can only contain 200 characters or less")

	v.Check(comment.PostID != 0, "post", "Post id must be provided")
	v.Check(comment.PostID > 0, "post", "Post id must be valid")
}
