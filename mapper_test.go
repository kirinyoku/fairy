package fairy

import (
	_ "embed"
	"encoding/json"
	"os"
	"testing"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/fairy/store"
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
