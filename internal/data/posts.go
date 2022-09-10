package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/AthfanFasee/blog-post-backend/internal/validator"
	"github.com/lib/pq"
)

// Later change this struct and validator for Post
type Post struct {
ID int64 
CreatedAt time.Time 
Title string 
PostText string
Img string 
ReadTime ReadTime // If we use our custom Runtime type here (which has the underlying type int32) go will use Runtime type's method MarshalJSON to encode this to JSON and it will be encoded to Runtime type (a string in the format "<runtime> mins") instead of int
LikedBy []int 
CreatedBy int64 
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

	return p.DB.QueryRow(query, args...).Scan(&post.ID, &post.CreatedAt)
}



func (p PostModel) Get(id int64) (*Post, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, title, post_text, img, read_time, liked_by, created_by, created_at 
	FROM posts
	WHERE id = $1`

	var post Post

	err := p.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.PostText,
		&post.Img,
		&post.ReadTime,
		pq.Array(&post.LikedBy),
		&post.CreatedBy,
		&post.CreatedAt,
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



func (p PostModel) Update(post *Post) error {
	query := `
	UPDATE posts 
	SET title = $1, post_text = $2, img = $3, read_time = $4 
	WHERE id = $5`

	args := []interface{}{
		post.Title,
		post.PostText,
		post.Img,
		post.ReadTime,
		post.ID,
	}

	result, err := p.DB.Exec(query, args...)
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

func (p PostModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM posts
	WHERE id = $1`

	result, err := p.DB.Exec(query, id)
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

func ValidatePost(v *validator.Validator, post *Post) {
	v.Check(post.Title != "", "title", "must be provided")
	v.Check(len(post.Title) <= 100, "title", "can only contain 100 characters or less")
	
	v.Check(post.PostText != "", "postText", "must be provided")
	v.Check(post.Img != "", "img", "must be provided")
	
	v.Check(post.ReadTime != 0, "readTime", "must be provided")
	v.Check(post.ReadTime > 0, "readTime", "must be provided")
}