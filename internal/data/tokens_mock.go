package data

import (
	"fmt"
	"time"
)

var mockToken = &Token{
	Plaintext: "Mocked text",
	Hash:      []byte("Mocked Hash"),
	Expiry:    time.Now().Add(5 * time.Hour),
	UserID:    1,
	Scope:     ScopeActivation,
}

type MockTokenModel struct {
	MockNew func(userID int64, ttl time.Duration, scope string) (*Token, error)
}

func (c MockTokenModel) Insert(token *Token) error {
	return nil
}

func (t MockTokenModel) New(userID int64, timeToLive time.Duration, scope string) (*Token, error) {
	if t.MockNew != nil {
		return t.MockNew(userID, timeToLive, scope)
	}

	switch {
	case userID == 1 && scope == ScopeActivation:
		return mockToken, nil
	default:
		return nil, fmt.Errorf("mocked Error")
	}
}

func (t MockTokenModel) DeleteAllForUser(scope string, userID int64) error {
	return nil
}
