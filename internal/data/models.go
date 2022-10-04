package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/AthfanFasee/blog-post-backend/internal/dto"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Comments interface {
		GetAllForPost(postID int64) ([]*dto.CommentResponseBody, error)
		Insert(comment *Comment) error
		Delete(id int64) error
	}
	Posts interface {
		GetAll(title string, filters Filters) ([]*dto.PostResponseBody, Metadata, error)
		Get(id int64) (*Post, error)
		GetWithUserName(id int64) (*Post, *string, error)
		Insert(post *Post) error
		Update(post *Post) error
		Delete(id int64) error
		AddLike(post *Post, userID int64) error
		RemoveLike(post *Post, userID int64) error
	}
	Tokens interface {
		Insert(token *Token) error
		New(userID int64, timeToLive time.Duration, scope string) (*Token, error)
		DeleteAllForUser(scope string, userID int64) error
	}
	Users interface {
		Insert(user *User) error
		GetByEmail(email string) (*User, error)
		Update(user *User) error
		GetForToken(tokenScope, tokenPlainText string) (*User, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Comments: CommentModel{DB: db},
		Posts:    PostModel{DB: db},
		Tokens:   TokenModel{DB: db},
		Users:    UserModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Comments: MockCommentModel{},
		Posts:    MockPostModel{},
		Tokens:   MockTokenModel{},
		Users:    MockUserModel{},
	}
}
