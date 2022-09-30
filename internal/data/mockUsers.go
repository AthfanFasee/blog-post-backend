package data

import (
	"time"
)

var mockUser = &User{
	ID:       1,
	CreatedAt: time.Now(),
	Name:      "Mocked Name",
	Email:      "Mocked Email",
	Password: Password{},
	Activated:   true,
	Version: 1,
}

type MockUserModel struct{}

func (u MockUserModel) Insert(user *User) error {
	return nil
}

func (u MockUserModel) GetByEmail(email string) (*User, error) {
	switch email {
	case "Mocked Email":
		return mockUser, nil
	default:
		return nil, ErrRecordNotFound
	}
}

func (u MockUserModel) Update(user *User) error {
	return nil
}

func (u MockUserModel) GetForToken(tokenScope, tokenPlainText string) (*User, error) {
	switch {
	case tokenScope == ScopeActivation:
		return mockUser, nil
	case tokenPlainText == "Mocked token":
		return mockUser, nil
	default:
		return nil, ErrRecordNotFound
	}
}