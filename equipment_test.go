package fairy

import (
	"testing"

	"github.com/kirinyoku/fairy/internal/store"
)

func TestDriveDisc_CountEffectiveRolls(t *testing.T) {
	tests := []struct {
		name        string
		subStats    []StatValue
		targetProps []PropertyID
		expected    int
	}{
		{
			name: "single matching target",
			subStats: []StatValue{
				{PropertyID: PropCritRate, Rolls: 2, Value: 4.8},
				{PropertyID: PropCritDMG, Rolls: 1, Value: 4.8},
				{PropertyID: PropATKFlat, Rolls: 1, Value: 19},
				{PropertyID: PropHPFlat, Rolls: 1, Value: 112},
			},
			targetProps: []PropertyID{PropCritRate},
			expected:    2,
		},
		{
			name: "multiple matching targets",
			subStats: []StatValue{
				{PropertyID: PropCritRate, Rolls: 3, Value: 7.2},
				{PropertyID: PropCritDMG, Rolls: 2, Value: 9.6},
				{PropertyID: PropATKPercent, Rolls: 1, Value: 3.0},
				{PropertyID: PropDEFPercent, Rolls: 1, Value: 4.8},
			},
			targetProps: []PropertyID{PropCritRate, PropCritDMG, PropATKPercent},
			expected:    6,
		},
		{
			name: "no matching targets",
			subStats: []StatValue{
				{PropertyID: PropHPFlat, Rolls: 2, Value: 224},
				{PropertyID: PropDEFFlat, Rolls: 3, Value: 45},
			},
			targetProps: []PropertyID{PropCritRate, PropCritDMG},
			expected:    0,
		},
		{
			name: "empty targets",
			subStats: []StatValue{
				{PropertyID: PropCritRate, Rolls: 3, Value: 7.2},
			},
			targetProps: []PropertyID{},
			expected:    0,
		},
		{
			name:        "empty sub-stats",
			subStats:    []StatValue{},
			targetProps: []PropertyID{PropCritRate},
			expected:    0,
		},
		{
			name: "duplicate target props do not double-count",
			subStats: []StatValue{
				{PropertyID: PropCritRate, Rolls: 3, Value: 7.2},
				{PropertyID: PropCritDMG, Rolls: 2, Value: 9.6},
			},
			targetProps: []PropertyID{PropCritRate, PropCritRate, PropCritDMG},
			expected:    5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disc := &DriveDisc{
				SubStats: tt.subStats,
			}
			got := disc.CountEffectiveRolls(tt.targetProps...)
			if got != tt.expected {
				t.Errorf("CountEffectiveRolls() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestDriveDiscs_BySlot(t *testing.T) {
	discs := DriveDiscs{
		{Slot: 1, Level: 15, UID: "disc-1"},
		{Slot: 3, Level: 12, UID: "disc-3"},
		{Slot: 6, Level: 15, UID: "disc-6"},
	}

	if d := discs.BySlot(1); d == nil || d.UID != "disc-1" {
		t.Errorf("BySlot(1) = %v, want UID disc-1", d)
	}
	if d := discs.BySlot(2); d != nil {
		t.Errorf("BySlot(2) = %v, want nil", d)
	}
	if d := discs.BySlot(3); d == nil || d.UID != "disc-3" {
		t.Errorf("BySlot(3) = %v, want UID disc-3", d)
	}
	if d := discs.BySlot(6); d == nil || d.UID != "disc-6" {
		t.Errorf("BySlot(6) = %v, want UID disc-6", d)
	}
}

func TestDriveDiscs_SetHelpers(t *testing.T) {
	setA := Set{ID: SetWoodpeckerElectro, Name: "Woodpecker Electro"}
	setB := Set{ID: SetPolarMetal, Name: "Polar Metal"}

	discs := DriveDiscs{
		{Slot: 1, Set: setA},
		{Slot: 2, Set: setA},
		{Slot: 3, Set: setA},
		{Slot: 4, Set: setA},
		{Slot: 5, Set: setB},
		{Slot: 6, Set: setB},
	}

	// 1. SetCounts
	counts := discs.SetCounts()
	if counts[setA] != 4 {
		t.Errorf("SetCounts()[setA] = %d, want 4", counts[setA])
	}
	if counts[setB] != 2 {
		t.Errorf("SetCounts()[setB] = %d, want 2", counts[setB])
	}

	// 2. Has2Piece
	if !discs.Has2Piece(SetWoodpeckerElectro) {
		t.Errorf("Has2Piece(SetWoodpeckerElectro) = false, want true")
	}
	if !discs.Has2Piece(SetPolarMetal) {
		t.Errorf("Has2Piece(SetPolarMetal) = false, want true")
	}
	if discs.Has2Piece(SetID(99999)) {
		t.Errorf("Has2Piece(99999) = true, want false")
	}

	// 3. Has4Piece
	if !discs.Has4Piece(SetWoodpeckerElectro) {
		t.Errorf("Has4Piece(SetWoodpeckerElectro) = false, want true")
	}
	if discs.Has4Piece(SetPolarMetal) {
		t.Errorf("Has4Piece(SetPolarMetal) = true, want false")
	}
	if discs.Has4Piece(SetID(99999)) {
		t.Errorf("Has4Piece(99999) = true, want false")
	}

	// Empty collection
	var empty DriveDiscs
	if len(empty.SetCounts()) != 0 {
		t.Errorf("empty.SetCounts() = %v, want empty map", empty.SetCounts())
	}
	if empty.Has2Piece(SetWoodpeckerElectro) || empty.Has4Piece(SetWoodpeckerElectro) {
		t.Errorf("empty.Has2Piece / Has4Piece returned true on empty discs")
	}
}

func TestAllSetIDs(t *testing.T) {
	sets := AllSetIDs()
	if len(sets) != 30 {
		t.Errorf("AllSetIDs() returned %d sets, want 30", len(sets))
	}

	// Verify defensive copy (modifying returned slice does not mutate internal array)
	sets[0] = SetID(99999)
	if AllSetIDs()[0] == SetID(99999) {
		t.Errorf("AllSetIDs() is not returning a defensive copy")
	}
}

func TestSetID_IsValid(t *testing.T) {
	for _, id := range AllSetIDs() {
		if !id.IsValid() {
			t.Errorf("SetID(%d).IsValid() = false, want true", id)
		}
	}

	if SetID(0).IsValid() {
		t.Errorf("SetID(0).IsValid() = true, want false")
	}
	if SetID(99999).IsValid() {
		t.Errorf("SetID(99999).IsValid() = true, want false")
	}
}

func TestAllSetIDs_StoreCoverage(t *testing.T) {
	st, err := store.Default()
	if err != nil {
		t.Fatalf("failed to load default store: %v", err)
	}

	allSetsMap := make(map[SetID]bool)
	for _, id := range AllSetIDs() {
		allSetsMap[id] = true
	}

	// 1. Check that every SetID constant has corresponding metadata in store
	for _, id := range AllSetIDs() {
		meta, ok := st.EquipmentSuitMeta(int(id))
		if !ok {
			t.Errorf("SetID %d has constant in Go code but no metadata in store", id)
			continue
		}
		name := st.Localize(meta.Name, "en")
		if name == "" || name == meta.Name {
			t.Errorf("SetID %d has empty or unlocalized name in store", id)
		}
	}

	// 2. Scan all possible suit IDs (30000..35000) and ensure none are missing from AllSetIDs()
	for id := 30000; id <= 35000; id += 100 {
		if _, ok := st.EquipmentSuitMeta(id); ok {
			if !allSetsMap[SetID(id)] {
				t.Errorf("Suit ID %d exists in store data but is MISSING from fairy SetID constants!", id)
			}
		}
	}
}
