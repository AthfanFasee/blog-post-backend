package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Comments CommentModel
	Posts  PostModel
	Tokens TokenModel
	Users  UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Comments: CommentModel{DB: db},
		Posts:  PostModel{DB: db},
		Tokens: TokenModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
