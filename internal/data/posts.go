package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/AthfanFasee/blog-post-backend/internal/dto"
	"github.com/AthfanFasee/blog-post-backend/internal/validator"
	"github.com/lib/pq"
)

type Post struct {
	ID        int64
	CreatedAt time.Time
	Title     string
	PostText  string
	Img       string
	ReadTime  dto.ReadTime
	LikedBy   []int64
	CreatedBy int64
	Version   int32
}

type PostModel struct {
	DB *sql.DB
}

func (p PostModel) GetAll(title string, filters Filters) ([]*dto.PostResponseBody, Metadata, error) {
	// Get post data along with name of the user who created it
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), p.id, p.title, p.post_text, p.img, p.read_time, p.liked_by, p.created_by, p.created_at, u.name
	FROM posts p
	INNER JOIN users u ON p.created_by = u.id
	WHERE (to_tsvector('english', title) @@ plainto_tsquery('english', $1) OR $1 = '')
	AND (created_by = $2 OR $2 = 0)
	ORDER BY %s %s, id %s
	LIMIT $3 OFFSET $4`, filters.sortParam(), filters.sortDirection(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, filters.ID, filters.limit(), filters.offset()}

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	posts := []*dto.PostResponseBody{}

	for rows.Next() {
		var post Post
		var userName string
		err := rows.Scan(
			&totalRecords,
			&post.ID,
			&post.Title,
			&post.PostText,
			&post.Img,
			&post.ReadTime,
			pq.Array(&post.LikedBy),
			&post.CreatedBy,
			&post.CreatedAt,
			&userName,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		PostResponseBody := dto.PostResponseBody{
			ID:        post.ID,
			Title:     post.Title,
			PostText:  post.PostText,
			Img:       post.Img,
			ReadTime:  post.ReadTime,
			LikedBy:   post.LikedBy,
			CreatedAt: post.CreatedAt,
			CreatedBy: post.CreatedBy,
			UserName:  userName,
		}

		posts = append(posts, &PostResponseBody)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.Limit)

	return posts, metadata, nil
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

func (p PostModel) GetWithUserName(id int64) (*Post, *string, error) {
	if id < 1 {
		return nil, nil, ErrRecordNotFound
	}

	query := `
	SELECT p.id, p.title, p.post_text, p.img, p.read_time, p.liked_by, p.created_by, p.created_at, u.name
	FROM posts p
	INNER JOIN users u ON p.created_by = u.id
	WHERE p.id = $1`

	var post Post
	var userName string

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
		&userName,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil, ErrRecordNotFound
		default:
			return nil, nil, err
		}
	}

	return &post, &userName, nil
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

func (p PostModel) Update(post *Post) error {
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
		return err
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

func (p PostModel) AddLike(post *Post, userID int64) error {
	// This SQL statement will prevent a user from liking a post twice
	query := `
	UPDATE posts SET 
	liked_by = (select array_agg(distinct x) from unnest(array_append(liked_by, $1)) t(x))
	WHERE id = $2
	RETURNING liked_by`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, userID, post.ID).Scan(pq.Array(&post.LikedBy))
}

func (p PostModel) RemoveLike(post *Post, userID int64) error {
	query := `
	UPDATE posts SET liked_by = array_remove(liked_by, $1)
	WHERE id = $2
	RETURNING liked_by`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, userID, post.ID).Scan(pq.Array(&post.LikedBy))
}

func ValidatePost(v *validator.Validator, post *Post) {
	v.Check(post.Title != "", "title", "Title must be provided")
	v.Check(len(post.Title) <= 100, "title", "Title can only contain 100 characters or less")

	v.Check(post.PostText != "", "postText", "Text must be provided")
	v.Check(post.Img != "", "img", "Image must be provided")

	v.Check(post.ReadTime != 0, "readTime", "Read time must be provided")
	v.Check(post.ReadTime > 0, "readTime", "Read time must be provided")
}
