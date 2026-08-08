package fairy

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/fairy/store"
)

func TestMapperSnapshot(t *testing.T) {
	data, err := os.ReadFile("internal/testdata/response.json")
	if err != nil {
		t.Fatalf("failed to read response.json: %v", err)
	}

	var raw zzz.Profile
	if err := json.Unmarshal(data, &raw); err != nil {
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

	snapshotPath := "internal/testdata/snapshot.json"

	// Create snapshot if it doesn't exist or if UPDATE_SNAPSHOTS=true
	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) || os.Getenv("UPDATE_SNAPSHOTS") == "true" {
		if err := os.WriteFile(snapshotPath, outData, 0644); err != nil {
			t.Fatalf("failed to write snapshot: %v", err)
		}
	}

	snapshot, err := os.ReadFile(snapshotPath)
	if err != nil {
		t.Fatalf("failed to read snapshot.json: %v", err)
	}

	if string(outData) != string(snapshot) {
		t.Errorf("snapshot mismatch!\nExpected:\n%s\nGot:\n%s", string(snapshot), string(outData))
	}
}
