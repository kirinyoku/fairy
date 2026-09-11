package fairy

import (
	_ "embed"
	"encoding/json"
	"os"
	"testing"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/fairy/internal/store"
)

//go:embed internal/testdata/response.json
var testResponseJSON []byte

//go:embed internal/testdata/snapshot.json
var testSnapshotJSON []byte

func TestMapperSnapshot(t *testing.T) {
	var raw zzz.Profile
	if err := json.Unmarshal(testResponseJSON, &raw); err != nil {
		t.Fatalf("failed to unmarshal response.json: %v", err)
	}

	s, err := store.Default()
	if err != nil {
		t.Fatalf("failed to init default store: %v", err)
	}

	mapper := newMapper(s, LangEN)
	p, err := mapper.ToProfile(&raw)
	if err != nil {
		t.Fatalf("failed to map profile: %v", err)
	}

	outData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal profile: %v", err)
	}

	// Update snapshot file on disk if explicitly requested with UPDATE_SNAPSHOTS=true
	if os.Getenv("UPDATE_SNAPSHOTS") == "true" {
		snapshotPath := "internal/testdata/snapshot.json"
		if err := os.WriteFile(snapshotPath, outData, 0644); err != nil {
			t.Fatalf("failed to write snapshot: %v", err)
		}
	}

	if string(outData) != string(testSnapshotJSON) {
		t.Errorf("snapshot mismatch!\nExpected:\n%s\nGot:\n%s", string(testSnapshotJSON), string(outData))
	}
}

func TestMapAgentSkills_Fallback(t *testing.T) {
	s, err := store.Default()
	if err != nil {
		t.Fatalf("failed to init default store: %v", err)
	}

	mapper := newMapper(s, LangEN)

	t.Run("empty skill level list falls back gracefully without panic", func(t *testing.T) {
		rawAvatar := &zzz.AvatarData{
			ID:                   1011, // Anby
			Level:                50,
			SkillLevelList:       nil,
			CoreSkillEnhancement: 0,
		}

		skills := mapper.mapAgentSkills(rawAvatar)
		if len(skills) == 0 {
			t.Fatal("expected non-empty skills for known avatar 1011")
		}

		for _, sk := range skills {
			if sk.Level != 0 {
				t.Errorf("expected skill level 0 when missing from list, got %d", sk.Level)
			}
			if sk.Name == "" {
				t.Errorf("expected localized skill name, got empty string")
			}
		}
	})

	t.Run("unknown avatar ID returns nil skills", func(t *testing.T) {
		rawAvatar := &zzz.AvatarData{
			ID: 999999,
		}

		skills := mapper.mapAgentSkills(rawAvatar)
		if skills != nil {
			t.Errorf("expected nil skills for unknown avatar ID, got %v", skills)
		}
	})
}

func TestFilterUnknownEntities(t *testing.T) {
	s, err := store.Default()
	if err != nil {
		t.Fatalf("failed to init default store: %v", err)
	}

	mapper := newMapper(s, LangEN)

	t.Run("unknown avatar returns nil from ToAgent", func(t *testing.T) {
		agent := mapper.ToAgent(&zzz.AvatarData{ID: 999999})
		if agent != nil {
			t.Errorf("expected nil agent for unknown ID 999999, got %+v", agent)
		}
	})

	t.Run("unknown weapon returns nil from mapWEngine", func(t *testing.T) {
		weapon := mapper.mapWEngine(&zzz.Weapon{ID: 999999})
		if weapon != nil {
			t.Errorf("expected nil weapon for unknown ID 999999, got %+v", weapon)
		}
	})

	t.Run("showcase filters out unknown avatars cleanly", func(t *testing.T) {
		raw := &zzz.Profile{
			PlayerInfo: zzz.PlayerInfo{
				ShowcaseDetail: &zzz.ShowcaseDetail{
					AvatarList: []zzz.AvatarData{
						{ID: 1011, Level: 50},   // Known (Anby)
						{ID: 999999, Level: 60}, // Unknown
						{ID: 1021, Level: 60},   // Known (Nekomata)
					},
				},
			},
		}

		p, err := mapper.ToProfile(raw)
		if err != nil {
			t.Fatalf("ToProfile failed: %v", err)
		}
		if len(p.Agents) != 2 {
			t.Fatalf("expected exactly 2 known agents, got %d", len(p.Agents))
		}
		if p.Agents[0].ID != 1011 || p.Agents[1].ID != 1021 {
			t.Errorf("unexpected agents in profile: %+v", p.Agents)
		}
	})

	t.Run("known avatar with unknown weapon has nil WEngine", func(t *testing.T) {
		raw := &zzz.Profile{
			PlayerInfo: zzz.PlayerInfo{
				ShowcaseDetail: &zzz.ShowcaseDetail{
					AvatarList: []zzz.AvatarData{
						{
							ID:    1011,
							Level: 50,
							Weapon: &zzz.Weapon{
								ID: 999999,
							},
						},
					},
				},
			},
		}

		p, err := mapper.ToProfile(raw)
		if err != nil {
			t.Fatalf("ToProfile failed: %v", err)
		}
		if len(p.Agents) != 1 {
			t.Fatalf("expected 1 agent, got %d", len(p.Agents))
		}
		if p.Agents[0].WEngine != nil {
			t.Errorf("expected nil WEngine for unknown weapon, got %+v", p.Agents[0].WEngine)
		}
	})
}
