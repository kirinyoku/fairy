package fairy

import (
	"testing"

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
