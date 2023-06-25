package main

import (
	"context"
	"expvar"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AthfanFasee/blog-post-backend/internal/assert"
	"github.com/AthfanFasee/blog-post-backend/internal/data"
)

func TestApplication_RecoverPanic(t *testing.T) {
	app := newTestApplication(t)

	// Create a handler function which panics.
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Panic!")
	})

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()

	// Use the recoverPanic middleware with the panicHandler.
	app.recoverPanic(panicHandler).ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	assert.StringContains(t, responseRecorder.Body.String(), "the server encountered a problem and could not process your request")
}

func TestSecureHeaders(t *testing.T) {
	app := newTestApplication(t)

	// Create a new httptest.ResponseRecorder and dummy http.Request.
	responseRecorder := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Pass the response recorder and request to our secureHeaders middleware.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	app.secureHeaders(next).ServeHTTP(responseRecorder, r)

	headers := map[string]string{
		"Content-Security-Policy": "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		"Referrer-Policy":         "origin-when-cross-origin",
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "deny",
		"X-XSS-Protection":        "0",
	}

	for header, expectedValue := range headers {
		actualValue := responseRecorder.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("Header %q = %q; want %q", header, actualValue, expectedValue)
		}
	}
}

func TestApplication_Authenticate(t *testing.T) {
	validToken := data.GenerateTestToken()

	tests := []struct {
		name               string
		authorizationToken string
		wantStatus         int
	}{
		{"Authorization header missing", "", http.StatusOK},
		{"Authorization header in wrong format", "InvalidToken", http.StatusUnauthorized},
		{"Invalid token", "Bearer InvalidToken", http.StatusUnauthorized},
		{"Valid token and user found", "Bearer " + validToken, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApplication(t)
			// Use the MockUserModel instead of the real UserModel.
			app.models.Users = &data.MockUserModel{}

			request, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.authorizationToken != "" {
				request.Header.Set("Authorization", tt.authorizationToken)
			}

			responseRecorder := httptest.NewRecorder()

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})

			app.authenticate(next).ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tt.wantStatus {
				t.Errorf("got status %d; want %d", responseRecorder.Code, tt.wantStatus)
			}
		})
	}
}

// The tests for requireAuthenticatedUser should only focus on the logic within itself.
func TestApplication_RequireAuthenticatedUser(t *testing.T) {
	tests := []struct {
		name        string
		isAnonymous bool
		wantStatus  int
	}{
		{"Anonymous user", true, http.StatusUnauthorized},
		{"Authenticated user", false, http.StatusOK},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApplication(t)

			request, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			responseRecorder := httptest.NewRecorder()

			// Create a user and set it in the context of the request manually.
			var user *data.User
			if tt.isAnonymous {
				user = data.AnonymousUser
			} else {
				user = &data.User{Activated: true}
			}

			ctx := context.WithValue(request.Context(), userContextKey, user)
			request = request.WithContext(ctx)

			app.requireAuthenticatedUser(next).ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tt.wantStatus {
				t.Errorf("got status %d; want %d", responseRecorder.Code, tt.wantStatus)
			}
		})
	}
}

// The tests for requireActivatedUser should only focus on the logic within itself.
func TestApplication_RequireActivatedUser(t *testing.T) {
	tests := []struct {
		name       string
		activated  bool
		wantStatus int
	}{
		{"Not activated user", false, http.StatusForbidden},
		{"Activated user", true, http.StatusOK},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApplication(t)

			request, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			responseRecorder := httptest.NewRecorder()

			// Create a user with the activated status from the test case and set it in the context of the request manually.
			var user *data.User
			if tt.activated {
				user = &data.User{Activated: true}
			} else {
				user = &data.User{Activated: false}
			}

			ctx := context.WithValue(request.Context(), userContextKey, user)
			request = request.WithContext(ctx)

			// Call requireAuthenticatedUser function first to set the authenticated user.
			app.requireAuthenticatedUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(responseRecorder, request)

			// Create a new Recorder for the second middleware.
			responseRecorder = httptest.NewRecorder()

			app.requireActivatedUser(next).ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tt.wantStatus {
				t.Errorf("got status %d; want %d", responseRecorder.Code, tt.wantStatus)
			}
		})
	}
}

func TestEnableCORS(t *testing.T) {
	// Prepare the application with CORS configuration.
	app := newTestApplication(t)
	app.config.cors.trustedOrigins = []string{"http://localhost:8080"}

	// Create a simple next handler that will be wrapped by enableCORS middleware.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the next handler with the middleware.
	corsEnabledNext := app.enableCORS(next)

	t.Run("trusted origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Origin", "http://localhost:8080")

		rec := httptest.NewRecorder()
		corsEnabledNext.ServeHTTP(rec, req)

		if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:8080" {
			t.Errorf("Expected %s but got %s", "http://localhost:8080", rec.Header().Get("Access-Control-Allow-Origin"))
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("untrusted origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Origin", "http://malicious.com")

		rec := httptest.NewRecorder()
		corsEnabledNext.ServeHTTP(rec, req)

		if rec.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Errorf("Expected %s but got %s", "", rec.Header().Get("Access-Control-Allow-Origin"))
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, rec.Code)
		}
	})
}

func TestMetrics(t *testing.T) {
	app := newTestApplication(t)
	app.config.metrics = true

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Give a small time period to the process of measuring the processing time.
		time.Sleep(1 * time.Microsecond)
		w.WriteHeader(http.StatusOK)
	})

	metricsEnabledNext := app.metrics(next)

	t.Run("metrics are collected", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		metricsEnabledNext.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, rec.Code)
		}

		totalRequestsReceived := expvar.Get("total_requests_received").(*expvar.Int)
		if totalRequestsReceived.Value() != 1 {
			t.Errorf("Expected %d but got %d", 1, totalRequestsReceived.Value())
		}

		totalResponsesSent := expvar.Get("total_responses_sent").(*expvar.Int)
		if totalResponsesSent.Value() != 1 {
			t.Errorf("Expected %d but got %d", 1, totalResponsesSent.Value())
		}

		totalProcessingTimeMicroseconds := expvar.Get("total_processing_time_μs").(*expvar.Int)
		if totalProcessingTimeMicroseconds.Value() <= 0 {
			t.Errorf("Expected processing time to be greater than %d but got %d", 0, totalProcessingTimeMicroseconds.Value())
		}

		totalResponsesSentByStatus := expvar.Get("total_responses_sent_by_status").(*expvar.Map)
		var responseCount expvar.Int
		totalResponsesSentByStatus.Do(func(kv expvar.KeyValue) {
			if kv.Key == "200" {
				responseCount = *(kv.Value.(*expvar.Int))
			}
		})
		if responseCount.Value() != 1 {
			t.Errorf("Expected %d but got %d", 1, responseCount.Value())
		}
	})
}
