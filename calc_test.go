package fairy

import (
	"testing"

	"github.com/kirinyoku/fairy/store"
)

type mockStore struct {
	store.MetadataStore
}

func (m mockStore) WeaponLevelTemplate(rarity, level int) (store.WeaponLevelTemplate, bool) {
	if rarity == 4 && level == 60 {
		return store.WeaponLevelTemplate{
			MainStat:           2000,
			SubStatDenominator: 4000,
		}, true
	}
	return store.WeaponLevelTemplate{}, false
}

func (m mockStore) WeaponStarTemplate(rarity, phase int) (store.WeaponStarTemplate, bool) {
	if rarity == 4 && phase == 0 {
		return store.WeaponStarTemplate{
			MainStat: 500,
			SubStat:  500,
		}, true
	}
	return store.WeaponStarTemplate{}, false
}

func TestCalcAgentBaseStat(t *testing.T) {
	meta := store.AvatarMeta{
		BaseProps: map[int]int{
			int(PropBaseHP): 1000,
		},
		GrowthProps: map[int]int{
			int(PropBaseHP): 50000, // Growth is divided by 10000 -> 5.0 per level
		},
		PromotionProps: []map[int]int{
			{int(PropBaseHP): 100}, // Promotion 1
			{int(PropBaseHP): 250}, // Promotion 2
		},
		CoreEnhancementProps: []map[int]int{
			{int(PropBaseHP): 0},  // 0
			{int(PropBaseHP): 50}, // 1
		},
	}

	tests := []struct {
		name      string
		level     int
		promotion int
		core      int
		expected  float64
	}{
		{"Level 1, no promo, no core", 1, 0, 0, 1000.0},
		{"Level 2, no promo, no core", 2, 0, 0, 1005.0},
		{"Level 2, promo 1, no core", 2, 1, 0, 1105.0},
		{"Level 10, promo 2, core 1", 10, 2, 1, 1000.0 + 45.0 + 250.0 + 50.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := calcAgentBaseStat(meta, int(PropBaseHP), tt.level, tt.promotion, tt.core)
			if val != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}

func TestCalcWEngineMainStat(t *testing.T) {
	meta := store.WeaponMeta{
		Rarity: 4,
		MainStat: store.PropertyStat{
			PropertyID:    int(PropBaseATK),
			PropertyValue: 50,
		},
	}

	ms := mockStore{}

	tests := []struct {
		name     string
		level    int
		phase    int
		expected int
	}{
		{"Level 60, Phase 0", 60, 0, 62}, // 50 * (1 + 2000/10000 + 500/10000) = 50 * 1.25 = 62.5 -> 62
		{"Level 1, Phase 0", 1, 0, 52},   // Level mult = 0, star mult = 500. 50 * 1.05 = 52.5 -> 52
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := calcWEngineMainStat(ms, meta, tt.level, tt.phase)
			if val != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}

func TestCalcWEngineSecondaryStat(t *testing.T) {
	meta := store.WeaponMeta{
		Rarity: 4,
		SecondaryStat: store.PropertyStat{
			PropertyID:    int(PropATKPercent),
			PropertyValue: 100, // Say 100 base
		},
	}

	ms := mockStore{}

	tests := []struct {
		name     string
		level    int
		phase    int
		expected int
	}{
		{"Level 60, Phase 0", 60, 0, 262}, // 100 * (10000/4000) * (1 + 500/10000) = 100 * 2.5 * 1.05 = 262.5 -> 262
		{"Level 1, Phase 0", 1, 0, 105},   // Level mult = 1, star mult = 500. 100 * 1.05 = 105
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := calcWEngineSecondaryStat(ms, meta, tt.level, tt.phase)
			if val != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}
