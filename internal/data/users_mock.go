package data

import (
	"time"
)

var mockUser = &User{
	ID:        1,
	CreatedAt: time.Now(),
	Name:      "Mocked Name",
	Email:     "Mocked Email",
	Password:  Password{},
	Activated: true,
	Version:   1,
}

var mockUserModel = MockUserModel{
	UserActivated: true,
	UserAnonymous: false,
}

type MockUserModel struct {
	UserActivated bool
	UserAnonymous bool
}

func (u MockUserModel) Insert(user *User) error {
	return nil
}

func (MockUserModel) GetByEmail(email string) (*User, error) {
	switch email {
	case "mocked@email.com":
		return mockUser, nil
	default:
		if mockUserModel.UserAnonymous {
			return AnonymousUser, nil
		}
		return nil, ErrRecordNotFound
	}
}

func (MockUserModel) Update(user *User) error {
	return nil
}

func (MockUserModel) GetForToken(tokenScope, tokenPlainText string) (*User, error) {
	user := *mockUser
	user.Activated = mockUserModel.UserActivated
	switch {
	case tokenScope == ScopeActivation:
		return &user, nil
	case tokenPlainText != "":
		return &user, nil
	default:
		if mockUserModel.UserAnonymous {
			return AnonymousUser, nil
		}
		return nil, ErrRecordNotFound
	}
}

// SetMockUserPassword sets the password for the mockUser.
func SetMockUserPassword(passwordHash []byte) {
	mockUser.Password.hash = passwordHash
}
