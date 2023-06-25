package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/AthfanFasee/blog-post-backend/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateAuthenticationTokenHandler(t *testing.T) {
	// Create a function to compare maps.
	mapEquals := func(x, y map[string]interface{}) bool {
		if len(x) != len(y) {
			return false
		}
		for k, v1 := range x {
			if v2, ok := y[k]; !ok || fmt.Sprintf("%v", v1) != fmt.Sprintf("%v", v2) {
				return false
			}
		}
		return true
	}
	// Create hashed password using bcrypt.
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 12)
	// Assign the hashed password to mockUser.
	data.SetMockUserPassword(hashedPassword)

	app := &application{
		models: data.Models{
			Users: data.MockUserModel{
				UserActivated: true,
				UserAnonymous: false,
			},
			Tokens: data.MockTokenModel{
				MockNew: func(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
					if userID != 1 || ttl != 30*24*time.Hour || scope != data.ScopeAuthentication {
						t.Error("Unexpected parameters to MockNew")
					}

					return &data.Token{
						UserID: userID,
						Expiry: time.Now().Add(ttl),
						Scope:  scope,
						Hash:   []byte("mocked-hash"),
					}, nil
				},
			},
		},
	}

	tests := []struct {
		name           string
		requestBody    string
		wantStatusCode int
		wantBody       map[string]interface{}
	}{
		{
			name:           "Valid request",
			requestBody:    `{"email":"mocked@email.com","password":"password"}`,
			wantStatusCode: http.StatusCreated,
			wantBody: map[string]interface{}{
				"authentication_token": map[string]interface{}{
					"token":  "",
					"userID": float64(1), // After unmarshalling, numeric values become float64, not int
				},
				"userName": "Mocked Name",
			},
		},
		{
			name:           "Invalid email",
			requestBody:    `{"email":"invalid","password":"password"}`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "Non-existant user",
			requestBody:    `{"email":"non-existent@example.com","password":"password"}`,
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			responseRecorder := httptest.NewRecorder()

			app.createAuthenticationTokenHandler(responseRecorder, request)

			res := responseRecorder.Result()
			body, _ := ioutil.ReadAll(res.Body)

			if responseRecorder.Code != tt.wantStatusCode {
				t.Errorf("Want status %d; got %d. Body: %s", tt.wantStatusCode, res.StatusCode, body)
			}

			if tt.wantBody != nil {
				// Unmarshal the response into a map, then delete the "expiry" field before comparison.
				var responseMap map[string]interface{}
				if err := json.Unmarshal(body, &responseMap); err != nil {
					t.Errorf("Failed to unmarshal response body: %v", err)
				}
				tokenMap := responseMap["authentication_token"].(map[string]interface{})
				// Remove the expiry key as it changes rapidly and is hard to test.
				delete(tokenMap, "expiry")

				if !mapEquals(responseMap, tt.wantBody) {
					t.Errorf("Unexpected body. Want \n%+v\n; got \n%+v", tt.wantBody, responseMap)
				}
			}
			res.Body.Close()
		})
	}
}
