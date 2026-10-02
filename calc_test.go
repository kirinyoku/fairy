package fairy

import (
	"math"
	"testing"

	"github.com/kirinyoku/fairy/internal/store"
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
	if rarity == 4 && phase == 5 {
		return store.WeaponStarTemplate{
			MainStat: 44610,
			SubStat:  15000,
		}, true
	}
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
		name       string
		breakLevel int
		expected   int
	}{
		{"BreakLevel 5 (Max)", 5, 250}, // 100 * (1 + 15000/10000) = 250
		{"BreakLevel 0 (Base)", 0, 105}, // 100 * (1 + 500/10000) = 105
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := calcWEngineSecondaryStat(ms, meta, tt.breakLevel)
			if val != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}

func TestCalcSheerForce(t *testing.T) {
	meta := store.AvatarMeta{
		BaseProps: map[int]int{
			int(PropBaseSheerForce): 200,
		},
		GrowthProps: map[int]int{
			int(PropBaseSheerForce): 10000, // 1 per level
		},
	}

	val := calcAgentBaseStat(meta, int(PropBaseSheerForce), 10, 0, 0)
	if val != 209.0 {
		t.Errorf("expected 209.0 base SheerForce, got %v", val)
	}
}

func TestCalcAdrenalineAccumulation_Rupture(t *testing.T) {
	meta := store.AvatarMeta{
		BaseProps: map[int]int{
			int(PropBaseRpRecover):   200, // 200 / 100 = 2.00
			int(PropBaseEnergyRegen): 0,
		},
		GrowthProps: map[int]int{},
	}

	val := calcAgentBaseStat(meta, int(PropBaseRpRecover), 1, 0, 0) / 100.0
	if val != 2.00 {
		t.Errorf("expected 2.00 base Adrenaline Accumulation, got %v", val)
	}
}

func TestCalcWEngineMainStat_BaseDEF(t *testing.T) {
	meta := store.WeaponMeta{
		Rarity: 4,
		MainStat: store.PropertyStat{
			PropertyID:    int(PropBaseDEF),
			PropertyValue: 29,
		},
	}

	ms := mockStore{}
	// At Lv60, phase 0: 29 * (1 + 2000/10000 + 500/10000) = 29 * 1.25 = 36.25 -> 36
	val := calcWEngineMainStat(ms, meta, 60, 0)
	if val != 36 {
		t.Errorf("expected 36 base DEF, got %d", val)
	}
}

func TestCalculateAgentStats_Armorer_Claret(t *testing.T) {
	st, err := store.Default()
	if err != nil {
		t.Fatalf("failed to load default store: %v", err)
	}

	agent := &Agent{
		ID:                   1611, // Claret
		Level:                60,
		Promotion:            6,
		CoreSkillEnhancement: 6,
		Specialty:            SpecialtyArmorer,
		Attribute:            AttributeElectric,
		WEngine: &WEngine{
			ID:           14161, // Scarlet Thirst
			Level:        60,
			Phase:        5,
			Modification: 1,
		},
	}

	calculateAgentStats(agent, st)

	// Base DEF: 441 (agent) + 431 (w-engine) = 872
	if agent.BaseStats.DEF != 872 {
		t.Errorf("expected Base DEF 872, got %v", agent.BaseStats.DEF)
	}

	// Base CRIT Rate: 5% + 28.8% (Core F) = 33.8%
	if math.Abs(agent.BaseStats.CritRate-0.338) > 1e-4 {
		t.Errorf("expected Base CritRate 0.338, got %v", agent.BaseStats.CritRate)
	}

	// Base SharpCritDMG: 150% (1.50)
	if math.Abs(agent.BaseStats.SharpCritDMG-1.50) > 1e-4 {
		t.Errorf("expected Base SharpCritDMG 1.50, got %v", agent.BaseStats.SharpCritDMG)
	}

	// Claret passive conversion: 50% base CD * 0.35 = 17.5% added to CRIT Rate
	// Final CritRate without discs: 33.8% + 17.5% = 51.3%
	expectedCritRate := 0.338 + (0.50 * 0.35)
	if math.Abs(agent.Stats.CritRate-expectedCritRate) > 1e-4 {
		t.Errorf("expected CritRate %v, got %v", expectedCritRate, agent.Stats.CritRate)
	}

	// UIStats formatting check
	agent.UIStats = formatAgentUIStats(agent, st, LangRU)
	uiList := agent.UIStats.List()

	// In UIStats.List(), SharpCritDMG should be placed right after CritDMG (index 6)
	if len(uiList) < 7 || uiList[6].PropertyID != PropBaseSharpCritDMG {
		t.Errorf("expected SharpCritDMG at index 6 in UIStats.List(), got %+v", uiList)
	}

	// Verify EnergyRegen ("0.00") and SheerForce ("0") are omitted for Armorer
	for _, s := range uiList {
		if s.PropertyID == PropBaseEpRecover || s.PropertyID == PropBaseEnergyRegen || s.PropertyID == PropBaseRpRecover {
			t.Errorf("unexpected energy regen stat in Armorer UIStats.List(): %+v", s)
		}
		if s.PropertyID == PropBaseSheerForce {
			t.Errorf("unexpected SheerForce in Armorer UIStats.List(): %+v", s)
		}
	}

	// Verify Stats.List() numeric alignment
	numericList := agent.Stats.List()
	if len(numericList) < 7 || numericList[6].PropertyID != PropBaseSharpCritDMG {
		t.Errorf("expected SharpCritDMG at index 6 in Stats.List(), got %+v", numericList)
	}
	for _, s := range numericList {
		if s.PropertyID == PropBaseEpRecover || s.PropertyID == PropBaseEnergyRegen || s.PropertyID == PropBaseRpRecover {
			t.Errorf("unexpected energy regen stat in Armorer Stats.List(): %+v", s)
		}
		if s.PropertyID == PropBaseSheerForce {
			t.Errorf("unexpected SheerForce in Armorer Stats.List(): %+v", s)
		}
	}
}

func TestCalculateAgentStats_Sunna_CoreBonus(t *testing.T) {
	st, err := store.Default()
	if err != nil {
		t.Fatalf("failed to load default store: %v", err)
	}

	agent := &Agent{
		ID:                   1491, // Sunna
		Level:                60,
		Promotion:            6,
		CoreSkillEnhancement: 6, // Rank F (gives 21% ATK and 75 Base ATK)
		Specialty:            SpecialtySupport,
		Attribute:            AttributePhysical,
		WEngine: &WEngine{
			ID:           13115, // Kaboom the Cannon (Бум-пушка)
			Level:        60,
			Phase:        5,
			Modification: 5,
		},
		DriveDiscs: DriveDiscs{
			Slots: []DriveDisc{
				{
					MainStat: StatValue{PropertyID: PropHPFlat, Value: 2200},
					SubStats: []StatValue{
						{PropertyID: PropDEFPercent, Value: 0.144},
						{PropertyID: PropATKPercent, Value: 0.090},
						{PropertyID: PropATKFlat, Value: 19},
						{PropertyID: PropCritDMG, Value: 0.048},
					},
				},
				{
					MainStat: StatValue{PropertyID: PropATKFlat, Value: 316},
					SubStats: []StatValue{
						{PropertyID: PropHPFlat, Value: 112},
						{PropertyID: PropATKPercent, Value: 0.120},
						{PropertyID: PropCritDMG, Value: 0.048},
						{PropertyID: PropAnomalyProficiency, Value: 18},
					},
				},
				{
					MainStat: StatValue{PropertyID: PropDEFFlat, Value: 184},
					SubStats: []StatValue{
						{PropertyID: PropPENFlat, Value: 18},
						{PropertyID: PropCritRate, Value: 0.024},
						{PropertyID: PropATKPercent, Value: 0.090},
						{PropertyID: PropCritDMG, Value: 0.096},
					},
				},
				{
					MainStat: StatValue{PropertyID: PropATKPercent, Value: 0.300},
					SubStats: []StatValue{
						{PropertyID: PropDEFFlat, Value: 45},
						{PropertyID: PropCritRate, Value: 0.048},
						{PropertyID: PropAnomalyProficiency, Value: 27},
						{PropertyID: PropCritDMG, Value: 0.048},
					},
				},
				{
					MainStat: StatValue{PropertyID: PropATKPercent, Value: 0.300},
					SubStats: []StatValue{
						{PropertyID: PropDEFFlat, Value: 30},
						{PropertyID: PropAnomalyProficiency, Value: 18},
						{PropertyID: PropHPFlat, Value: 224},
						{PropertyID: PropATKFlat, Value: 57},
					},
				},
				{
					MainStat: StatValue{PropertyID: PropEnergyRegenPercent, Value: 0.600},
					SubStats: []StatValue{
						{PropertyID: PropHPPercent, Value: 0.060},
						{PropertyID: PropPENFlat, Value: 18},
						{PropertyID: PropCritDMG, Value: 0.048},
						{PropertyID: PropATKFlat, Value: 57},
					},
				},
			},
			SetBonuses: []DriveDiscSetBonus{
				{Set: Set{ID: 33400}, Count: 4}, // Moonlit Lullaby 2pc (+20% Energy Regen)
				{Set: Set{ID: 31600}, Count: 2}, // Swing Jazz 2pc (+20% Energy Regen)
			},
		},
	}

	calculateAgentStats(agent, st)

	// Base ATK: 750 (Sunna + 75 from core) + 624 (weapon) = 1374
	if agent.BaseStats.ATK != 1374 {
		t.Errorf("expected Base ATK 1374, got %v", agent.BaseStats.ATK)
	}

	// Total ATK: floor(1374 * (1 + 0.90 + 0.21) + 449) = 3348 (was 3059 without 21% core bonus)
	if agent.Stats.ATK != 3348 {
		t.Errorf("expected total ATK 3348, got %v", agent.Stats.ATK)
	}

	// Energy Regen: base 1.0 * (1 + 0.50 [weapon] + 0.60 [disc 6] + 0.20 [moonlit] + 0.20 [swing]) = 2.50
	if math.Abs(agent.Stats.EnergyRegen-2.50) > 1e-4 {
		t.Errorf("expected Energy Regen 2.50, got %v", agent.Stats.EnergyRegen)
	}

	// HP: 11522
	if agent.Stats.HP != 11522 {
		t.Errorf("expected HP 11522, got %v", agent.Stats.HP)
	}

	// DEF: 945
	if agent.Stats.DEF != 945 {
		t.Errorf("expected DEF 945, got %v", agent.Stats.DEF)
	}
}

func TestCalculateAgentStats_Rina_PENRatio(t *testing.T) {
	st, err := store.Default()
	if err != nil {
		t.Fatalf("failed to load default store: %v", err)
	}

	agent := &Agent{
		ID:                   1211, // Rina
		Level:                60,
		Promotion:            5,
		CoreSkillEnhancement: 6, // Rank F (gives 14.4% Base PEN Ratio via 23101: 1440)
		Specialty:            SpecialtySupport,
		Attribute:            AttributeElectric,
	}

	calculateAgentStats(agent, st)

	// Base PEN Ratio should be 14.4% (0.144)
	if math.Abs(agent.BaseStats.PenRatio-0.144) > 1e-4 {
		t.Errorf("expected Base PEN Ratio 0.144 for Rina, got %v", agent.BaseStats.PenRatio)
	}
}
