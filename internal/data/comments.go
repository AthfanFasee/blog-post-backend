package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/AthfanFasee/blog-post-backend/internal/validator"
)

type Comment struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Text      string    `json:"text"`
	CreatedBy int64     `json:"createdBy"`
	PostID   int64     `json:"post"`
}

type CommentModel struct {
	DB *sql.DB 
}

func (c CommentModel) GetAllForPost(postID int64) ([]*Comment,  error) {
	query := `SELECT id, text, created_by, post_id, created_at
	FROM comments 
	WHERE post_id = $1
	ORDER BY id DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*Comment{}

	for rows.Next() {
		var comment Comment

		err := rows.Scan(
			&comment.ID,
			&comment.Text,
			&comment.CreatedBy,
			&comment.PostID,
			&comment.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentModel) Insert(comment *Comment) error {
	query := `
	INSERT INTO comments (text, post_id, created_by)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	args := []interface{}{comment.Text, comment.PostID, comment.CreatedBy}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.ID, &comment.CreatedAt)
}

func (c *CommentModel) Delete(id int64) error {
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
		return  ErrRecordNotFound
	}

	return nil
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	v.Check(comment.Text != "", "text", "must be provided")
	v.Check(len(comment.Text) <= 200, "text", "can only contain 200 characters or less")

	v.Check(comment.PostID != 0, "post", "must be provided")
	v.Check(comment.PostID > 0, "post", "must be valid")
}