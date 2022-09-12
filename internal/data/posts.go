package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/AthfanFasee/blog-post-backend/internal/validator"
	"github.com/lib/pq"
)

// Later change this struct and validator for Post
type Post struct {
ID int64 `json:"id"`
CreatedAt time.Time `json:"createdAt"`
Title string `json:"title"`
PostText string `json:"postText"`
Img string `json:"img"`
ReadTime ReadTime `json:"readTime"` // If we use our custom Runtime type here (which has the underlying type int32) go will use Runtime type's method MarshalJSON to encode this to JSON and it will be encoded to Runtime type (a string in the format "<runtime> mins") instead of int
LikedBy []int `json:"likedBy"`
CreatedBy int64 `json:"createdBy"`
Version int32 `json:"-"`
}

type PostModel struct {
	DB *sql.DB
}

func (p PostModel) Insert(post *Post) error {
	query := `
		INSERT INTO posts (title, post_text, img, read_time, created_by) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at`
		
	args := []interface{}{post.Title, post.PostText, post.Img, post.ReadTime, post.CreatedBy}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&post.ID, &post.CreatedAt)
}



func (p PostModel) Get(id int64) (*Post, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, title, post_text, img, read_time, liked_by, created_by, created_at, version
	FROM posts
	WHERE id = $1`

	var post Post

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.PostText,
		&post.Img,
		&post.ReadTime,
		pq.Array(&post.LikedBy),
		&post.CreatedBy,
		&post.CreatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &post, nil	
}



func (p PostModel) Update(post *Post) (sql.Result, error) {
	query := `
	UPDATE posts 
	SET title = $1, post_text = $2, img = $3, read_time = $4, version = version + 1
	WHERE id = $5 AND version = $6`

	args := []interface{}{
		post.Title,
		post.PostText,
		post.Img,
		post.ReadTime,
		post.ID,
		post.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := p.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {		
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrEditConflict
	}

	return result, nil
}

func (p PostModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM posts
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrEditConflict
	}

	return nil
}

func ValidatePost(v *validator.Validator, post *Post) {
	v.Check(post.Title != "", "title", "must be provided")
	v.Check(len(post.Title) <= 100, "title", "can only contain 100 characters or less")
	
	v.Check(post.PostText != "", "postText", "must be provided")
	v.Check(post.Img != "", "img", "must be provided")
	
	v.Check(post.ReadTime != 0, "readTime", "must be provided")
	v.Check(post.ReadTime > 0, "readTime", "must be provided")
}