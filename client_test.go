package fairy

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
)

func TestNewClient_Defaults(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() unexpected error: %v", err)
	}

	if client.lang != LangEN {
		t.Errorf("client.lang = %v, want %v", client.lang, LangEN)
	}

	if client.store == nil {
		t.Errorf("client.store is nil, expected default store")
	}

	if client.apiClient == nil {
		t.Errorf("client.apiClient is nil")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	mockSt := mockStore{}
	enkaOpts := zzz.Options{
		UserAgent: "TestUserAgent/1.0",
	}

	client, err := NewClient(
		WithDefaultLang(LangRU),
		WithStore(mockSt),
		WithEnkaOptions(enkaOpts),
	)
	if err != nil {
		t.Fatalf("NewClient() unexpected error: %v", err)
	}

	if client.lang != LangRU {
		t.Errorf("client.lang = %v, want %v", client.lang, LangRU)
	}

	if client.store != mockSt {
		t.Errorf("client.store = %v, want %v", client.store, mockSt)
	}

	if client.enkaOpts.UserAgent != "TestUserAgent/1.0" {
		t.Errorf("client.enkaOpts.UserAgent = %q, want %q", client.enkaOpts.UserAgent, "TestUserAgent/1.0")
	}
}

func TestClient_Localize(t *testing.T) {
	client, err := NewClient(WithDefaultLang(LangEN))
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	t.Run("nil raw profile returns error", func(t *testing.T) {
		p, err := client.Localize(nil, LangEN)
		if err == nil {
			t.Errorf("expected error for nil raw profile, got nil")
		}
		if p != nil {
			t.Errorf("expected nil profile, got %v", p)
		}
	})

	t.Run("valid raw profile localizes successfully", func(t *testing.T) {
		raw := &zzz.Profile{
			Region: "Europe",
			PlayerInfo: zzz.PlayerInfo{
				SocialDetail: &zzz.SocialDetail{
					ProfileDetail: &zzz.ProfileDetail{
						UID:      100000001,
						Nickname: "Player1",
						Level:    50,
					},
				},
				ShowcaseDetail: &zzz.ShowcaseDetail{
					AvatarList: []zzz.AvatarData{
						{
							ID:    1011,
							Level: 60,
						},
					},
				},
			},
		}

		profile, err := client.Localize(raw, LangRU)
		if err != nil {
			t.Fatalf("client.Localize() failed: %v", err)
		}

		if profile == nil {
			t.Fatal("expected non-nil profile")
		}
		if profile.UID != "100000001" {
			t.Errorf("profile.UID = %q, want %q", profile.UID, "100000001")
		}
		if profile.Nickname != "Player1" {
			t.Errorf("profile.Nickname = %q, want %q", profile.Nickname, "Player1")
		}
		if len(profile.Agents) != 1 {
			t.Fatalf("expected 1 agent, got %d", len(profile.Agents))
		}
		if profile.Agents[0].ID != 1011 {
			t.Errorf("profile.Agents[0].ID = %d, want 1011", profile.Agents[0].ID)
		}
		// Verify agent.UIStats is correctly localized in Russian
		if profile.Agents[0].UIStats.CritRate.Name != "Шанс крит. попадания" {
			t.Errorf("profile.Agents[0].UIStats.CritRate.Name = %q, want %q", profile.Agents[0].UIStats.CritRate.Name, "Шанс крит. попадания")
		}
	})
}

func TestGlobal_Localize(t *testing.T) {
	t.Run("nil raw profile returns error", func(t *testing.T) {
		p, err := Localize(nil, LangEN)
		if err == nil {
			t.Errorf("expected error for nil raw profile, got nil")
		}
		if p != nil {
			t.Errorf("expected nil profile, got %v", p)
		}
	})

	t.Run("valid raw profile localizes via global function", func(t *testing.T) {
		raw := &zzz.Profile{
			Region: "Europe",
			PlayerInfo: zzz.PlayerInfo{
				SocialDetail: &zzz.SocialDetail{
					ProfileDetail: &zzz.ProfileDetail{
						UID:      100000001,
						Nickname: "Player1",
						Level:    50,
					},
				},
				ShowcaseDetail: &zzz.ShowcaseDetail{
					AvatarList: []zzz.AvatarData{
						{
							ID:    1011,
							Level: 60,
						},
					},
				},
			},
		}

		profile, err := Localize(raw, LangRU)
		if err != nil {
			t.Fatalf("Localize() failed: %v", err)
		}
		if profile == nil {
			t.Fatal("expected non-nil profile")
		}
		// Check pre-computed agent.UIStats
		if profile.Agents[0].UIStats.CritRate.Name != "Шанс крит. попадания" {
			t.Errorf("profile.Agents[0].UIStats.CritRate.Name = %q, want %q", profile.Agents[0].UIStats.CritRate.Name, "Шанс крит. попадания")
		}
		// FormattedUIStats(LangEN) dynamically formats into English on demand
		enStats := profile.Agents[0].FormattedUIStats(LangEN)
		if enStats.CritRate.Name != "CRIT Rate" {
			t.Errorf("enStats.CritRate.Name = %q, want %q", enStats.CritRate.Name, "CRIT Rate")
		}
	})
}

func TestClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedError error
	}{
		{
			name:          "404 player not found",
			statusCode:    http.StatusNotFound,
			responseBody:  `{"message":"player not found"}`,
			expectedError: ErrProfileNotFound,
		},
		{
			name:          "429 rate limited",
			statusCode:    http.StatusTooManyRequests,
			responseBody:  `{"message":"rate limit exceeded"}`,
			expectedError: ErrRateLimit,
		},
		{
			name:          "500 internal server error",
			statusCode:    http.StatusInternalServerError,
			responseBody:  `{"message":"internal server error"}`,
			expectedError: ErrMaintenance,
		},
		{
			name:          "503 service unavailable",
			statusCode:    http.StatusServiceUnavailable,
			responseBody:  `{"message":"service under maintenance"}`,
			expectedError: ErrMaintenance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client, err := NewClient(WithEnkaOptions(zzz.Options{
				BaseURL: server.URL,
				Retry:   &zzz.RetryOptions{MaxAttempts: 1},
			}))
			if err != nil {
				t.Fatalf("NewClient() failed: %v", err)
			}

			// Test GetProfile
			_, err = client.GetProfile(context.Background(), "100000001")
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("GetProfile() error = %v, want %v", err, tt.expectedError)
			}

			// Test GetProfileWithLang
			_, err = client.GetProfileWithLang(context.Background(), "100000001", LangRU)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("GetProfileWithLang() error = %v, want %v", err, tt.expectedError)
			}

			// Test GetRawProfile
			_, err = client.GetRawProfile(context.Background(), "100000001")
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("GetRawProfile() error = %v, want %v", err, tt.expectedError)
			}
		})
	}

	t.Run("network error on invalid host", func(t *testing.T) {
		client, err := NewClient(WithEnkaOptions(zzz.Options{
			BaseURL: "http://127.0.0.1:0",
			Retry:   &zzz.RetryOptions{MaxAttempts: 1},
		}))
		if err != nil {
			t.Fatalf("NewClient() failed: %v", err)
		}

		_, err = client.GetProfile(context.Background(), "100000001")
		if err == nil {
			t.Fatalf("expected network error, got nil")
		}
		if !errors.Is(err, ErrNetwork) {
			t.Errorf("expected ErrNetwork, got %v", err)
		}
	})

	t.Run("canceled context returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client, err := NewClient(WithEnkaOptions(zzz.Options{
			BaseURL: server.URL,
			Retry:   &zzz.RetryOptions{MaxAttempts: 1},
		}))
		if err != nil {
			t.Fatalf("NewClient() failed: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		_, err = client.GetProfile(ctx, "100000001")
		if err == nil {
			t.Fatalf("expected error on canceled context, got nil")
		}
	})

	t.Run("deadline exceeded context returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client, err := NewClient(WithEnkaOptions(zzz.Options{
			BaseURL: server.URL,
			Retry:   &zzz.RetryOptions{MaxAttempts: 1},
		}))
		if err != nil {
			t.Fatalf("NewClient() failed: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		_, err = client.GetProfile(ctx, "100000001")
		if err == nil {
			t.Fatalf("expected timeout error, got nil")
		}
	})

	t.Run("successful profile fetch through mock server", func(t *testing.T) {
		mockProfileJSON := `{
			"region": "Europe",
			"PlayerInfo": {
				"SocialDetail": {
					"ProfileDetail": {
						"Uid": 1504687050,
						"Nickname": "Belle",
						"Level": 50
					}
				},
				"ShowcaseDetail": {
					"AvatarList": [
						{
							"Id": 1011,
							"Level": 60
						}
					]
				}
			}
		}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(mockProfileJSON))
		}))
		defer server.Close()

		client, err := NewClient(WithEnkaOptions(zzz.Options{
			BaseURL: server.URL,
		}))
		if err != nil {
			t.Fatalf("NewClient() failed: %v", err)
		}

		profile, err := client.GetProfile(context.Background(), "1504687050")
		if err != nil {
			t.Fatalf("GetProfile() failed: %v", err)
		}
		if profile == nil {
			t.Fatal("expected non-nil profile")
		}
		if profile.UID != "1504687050" {
			t.Errorf("profile.UID = %q, want %q", profile.UID, "1504687050")
		}
		if profile.Nickname != "Belle" {
			t.Errorf("profile.Nickname = %q, want %q", profile.Nickname, "Belle")
		}
		if len(profile.Agents) != 1 {
			t.Fatalf("expected 1 agent, got %d", len(profile.Agents))
		}
	})
}
