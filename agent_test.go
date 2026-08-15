package fairy

import (
	"math"
	"testing"
)

func TestAgent_SubStatTotals(t *testing.T) {
	tests := []struct {
		name     string
		discs    []DriveDisc
		expected []StatValue
	}{
		{
			name:     "no drive discs",
			discs:    nil,
			expected: []StatValue{},
		},
		{
			name: "single disc preserves substats",
			discs: []DriveDisc{
				{
					SubStats: []StatValue{
						{PropertyID: PropCritRate, Name: "CRIT Rate", Value: 4.8, Rolls: 2, IsPercent: true},
						{PropertyID: PropATKFlat, Name: "ATK", Value: 19.0, Rolls: 1, IsPercent: false},
					},
				},
			},
			expected: []StatValue{
				{PropertyID: PropCritRate, Name: "CRIT Rate", Value: 4.8, Rolls: 2, IsPercent: true},
				{PropertyID: PropATKFlat, Name: "ATK", Value: 19.0, Rolls: 1, IsPercent: false},
			},
		},
		{
			name: "multiple discs aggregate values and rolls while preserving first appearance order",
			discs: []DriveDisc{
				{
					Slot: 1,
					SubStats: []StatValue{
						{PropertyID: PropCritRate, Name: "CRIT Rate", Value: 2.4, Rolls: 1, IsPercent: true},
						{PropertyID: PropATKPercent, Name: "ATK%", Value: 3.0, Rolls: 1, IsPercent: true},
					},
				},
				{
					Slot: 2,
					SubStats: []StatValue{
						{PropertyID: PropCritDMG, Name: "CRIT DMG", Value: 9.6, Rolls: 2, IsPercent: true},
						{PropertyID: PropCritRate, Name: "CRIT Rate", Value: 4.8, Rolls: 2, IsPercent: true},
					},
				},
				{
					Slot: 3,
					SubStats: []StatValue{
						{PropertyID: PropATKPercent, Name: "ATK%", Value: 6.0, Rolls: 2, IsPercent: true},
						{PropertyID: PropHPFlat, Name: "HP", Value: 112.0, Rolls: 1, IsPercent: false},
					},
				},
			},
			expected: []StatValue{
				// First seen: PropCritRate (from disc 1) -> 2.4 + 4.8 = 7.2, rolls = 1 + 2 = 3
				{PropertyID: PropCritRate, Name: "CRIT Rate", Value: 7.2, Rolls: 3, IsPercent: true},
				// Second seen: PropATKPercent (from disc 1) -> 3.0 + 6.0 = 9.0, rolls = 1 + 2 = 3
				{PropertyID: PropATKPercent, Name: "ATK%", Value: 9.0, Rolls: 3, IsPercent: true},
				// Third seen: PropCritDMG (from disc 2) -> 9.6, rolls = 2
				{PropertyID: PropCritDMG, Name: "CRIT DMG", Value: 9.6, Rolls: 2, IsPercent: true},
				// Fourth seen: PropHPFlat (from disc 3) -> 112.0, rolls = 1
				{PropertyID: PropHPFlat, Name: "HP", Value: 112.0, Rolls: 1, IsPercent: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &Agent{
				DriveDiscs: tt.discs,
			}
			got := agent.SubStatTotals()

			if len(got) != len(tt.expected) {
				t.Fatalf("SubStatTotals() returned %d items, want %d", len(got), len(tt.expected))
			}

			for i := range got {
				if got[i].PropertyID != tt.expected[i].PropertyID {
					t.Errorf("[%d] PropertyID = %v, want %v", i, got[i].PropertyID, tt.expected[i].PropertyID)
				}
				if got[i].Rolls != tt.expected[i].Rolls {
					t.Errorf("[%d] Rolls = %v, want %v", i, got[i].Rolls, tt.expected[i].Rolls)
				}
				// Float comparison with tolerance for precision
				if math.Abs(got[i].Value-tt.expected[i].Value) > 1e-6 {
					t.Errorf("[%d] Value = %v, want %v", i, got[i].Value, tt.expected[i].Value)
				}
				if got[i].Name != tt.expected[i].Name {
					t.Errorf("[%d] Name = %v, want %v", i, got[i].Name, tt.expected[i].Name)
				}
				if got[i].IsPercent != tt.expected[i].IsPercent {
					t.Errorf("[%d] IsPercent = %v, want %v", i, got[i].IsPercent, tt.expected[i].IsPercent)
				}
			}
		})
	}
}

func TestAgent_CountEffectiveRolls(t *testing.T) {
	agent := &Agent{
		DriveDiscs: []DriveDisc{
			{
				SubStats: []StatValue{
					{PropertyID: PropCritRate, Rolls: 2},
					{PropertyID: PropATKFlat, Rolls: 1},
				},
			},
			{
				SubStats: []StatValue{
					{PropertyID: PropCritDMG, Rolls: 3},
					{PropertyID: PropCritRate, Rolls: 1},
				},
			},
			{
				SubStats: []StatValue{
					{PropertyID: PropPENRatio, Rolls: 2},
					{PropertyID: PropHPPercent, Rolls: 1},
				},
			},
		},
	}

	tests := []struct {
		name        string
		agent       *Agent
		targetProps []PropertyID
		expected    int
	}{
		{
			name:        "crit stats only (rate + dmg)",
			agent:       agent,
			targetProps: []PropertyID{PropCritRate, PropCritDMG},
			expected:    6, // Disc 1 (2) + Disc 2 (3 + 1) = 6
		},
		{
			name:        "single prop",
			agent:       agent,
			targetProps: []PropertyID{PropPENRatio},
			expected:    2,
		},
		{
			name:        "non matching prop",
			agent:       agent,
			targetProps: []PropertyID{PropDEFPercent},
			expected:    0,
		},
		{
			name:        "duplicate target props do not overcount",
			agent:       agent,
			targetProps: []PropertyID{PropCritRate, PropCritRate, PropCritDMG},
			expected:    6,
		},
		{
			name:        "empty targets",
			agent:       agent,
			targetProps: []PropertyID{},
			expected:    0,
		},
		{
			name:        "agent with no discs",
			agent:       &Agent{},
			targetProps: []PropertyID{PropCritRate},
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.agent.CountEffectiveRolls(tt.targetProps...)
			if got != tt.expected {
				t.Errorf("CountEffectiveRolls() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestAgent_FormattedUIStats(t *testing.T) {
	agent := &Agent{
		Attribute: AttributeIce,
		BaseStats: Stats{
			HP:                 1000,
			ATK:                500,
			DEF:                300,
			Impact:             100,
			CritRate:           0.05,
			CritDMG:            0.50,
			AttributeDMGBonus:  0.10,
			AnomalyMastery:     90,
			AnomalyProficiency: 80,
			PenRatio:           0.0,
			PenFlat:            0,
			EnergyRegen:        1.20,
			SheerForce:         0,
		},
		Stats: Stats{
			HP:                 1500, // +500
			ATK:                1200, // +700
			DEF:                450,  // +150
			Impact:             120,  // +20
			CritRate:           0.65, // +0.60
			CritDMG:            1.40, // +0.90
			AttributeDMGBonus:  0.40, // +0.30
			AnomalyMastery:     115,  // +25
			AnomalyProficiency: 100,  // +20
			PenRatio:           0.24, // +0.24
			PenFlat:            30,   // +30
			EnergyRegen:        1.80, // +0.60
			SheerForce:         50,   // +50
		},
	}

	t.Run("stat breakdown and formatting values", func(t *testing.T) {
		ui := agent.FormattedUIStats(LangEN)

		// Flat stats
		if ui.HP.Base != "1000" || ui.HP.Added != "500" || ui.HP.Total != "1500" {
			t.Errorf("HP breakdown mismatch: Base=%s, Added=%s, Total=%s", ui.HP.Base, ui.HP.Added, ui.HP.Total)
		}
		if ui.ATK.Base != "500" || ui.ATK.Added != "700" || ui.ATK.Total != "1200" {
			t.Errorf("ATK breakdown mismatch: Base=%s, Added=%s, Total=%s", ui.ATK.Base, ui.ATK.Added, ui.ATK.Total)
		}

		// Percent stats
		if ui.CritRate.Base != "5.0%" || ui.CritRate.Added != "60.0%" || ui.CritRate.Total != "65.0%" {
			t.Errorf("CritRate breakdown mismatch: Base=%s, Added=%s, Total=%s", ui.CritRate.Base, ui.CritRate.Added, ui.CritRate.Total)
		}
		if ui.CritDMG.Base != "50.0%" || ui.CritDMG.Added != "90.0%" || ui.CritDMG.Total != "140.0%" {
			t.Errorf("CritDMG breakdown mismatch: Base=%s, Added=%s, Total=%s", ui.CritDMG.Base, ui.CritDMG.Added, ui.CritDMG.Total)
		}

		// Energy regen (precision 2)
		if ui.EnergyRegen.Base != "1.20" || ui.EnergyRegen.Added != "0.60" || ui.EnergyRegen.Total != "1.80" {
			t.Errorf("EnergyRegen breakdown mismatch: Base=%s, Added=%s, Total=%s", ui.EnergyRegen.Base, ui.EnergyRegen.Added, ui.EnergyRegen.Total)
		}

		// Check SVG Icon URLs are populated
		if ui.HP.IconURL == "" || ui.ATK.IconURL == "" || ui.CritRate.IconURL == "" {
			t.Errorf("expected IconURLs to be populated")
		}
	})

	t.Run("elemental damage property per attribute", func(t *testing.T) {
		tests := []struct {
			attr         Attribute
			expectedProp PropertyID
		}{
			{AttributePhysical, PropPhysicalDMGBonus},
			{AttributeFire, PropFireDMGBonus},
			{AttributeIce, PropIceDMGBonus},
			{AttributeElectric, PropElectricDMGBonus},
			{AttributeEther, PropEtherDMGBonus},
			{AttributeWind, PropWindDMGBonus},
		}

		for _, tt := range tests {
			t.Run(string(tt.attr), func(t *testing.T) {
				a := &Agent{Attribute: tt.attr}
				ui := a.FormattedUIStats(LangEN)
				if ui.AttributeDMGBonus.PropertyID != tt.expectedProp {
					t.Errorf("Attribute %s: PropertyID = %v, want %v", tt.attr, ui.AttributeDMGBonus.PropertyID, tt.expectedProp)
				}
			})
		}
	})

	t.Run("localization with language parameter", func(t *testing.T) {
		uiEN := agent.FormattedUIStats(LangEN)
		uiRU := agent.FormattedUIStats(LangRU)

		if uiEN.CritRate.Name == "" || uiRU.CritRate.Name == "" {
			t.Errorf("expected non-empty localized stat names")
		}
		if uiEN.CritRate.Name == uiRU.CritRate.Name {
			t.Errorf("expected localized names to differ between EN and RU, both got %q", uiEN.CritRate.Name)
		}
		if uiEN.CritRate.Name != "CRIT Rate" {
			t.Errorf("uiEN.CritRate.Name = %q, want %q", uiEN.CritRate.Name, "CRIT Rate")
		}
		if uiRU.CritRate.Name != "Шанс крит. попадания" {
			t.Errorf("uiRU.CritRate.Name = %q, want %q", uiRU.CritRate.Name, "Шанс крит. попадания")
		}
	})
}

func TestAllAttributes(t *testing.T) {
	attrs := AllAttributes()
	if len(attrs) != 10 {
		t.Fatalf("AllAttributes() returned %d attributes, want 10", len(attrs))
	}

	attrs[0] = Attribute("corrupted")
	freshAttrs := AllAttributes()
	if freshAttrs[0] == Attribute("corrupted") {
		t.Errorf("AllAttributes() is vulnerable to external slice mutation")
	}
}

func TestAttribute_IsValid(t *testing.T) {
	tests := []struct {
		attr  Attribute
		valid bool
	}{
		{AttributePhysical, true},
		{AttributeHonedEdge, true},
		{AttributeFire, true},
		{AttributeIce, true},
		{AttributeFrost, true},
		{AttributeElectric, true},
		{AttributeEther, true},
		{AttributeAuricInk, true},
		{AttributeWind, true},
		{AttributeLumiflux, true},
		{Attribute("Quantum"), false},
		{Attribute(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.attr), func(t *testing.T) {
			if got := tt.attr.IsValid(); got != tt.valid {
				t.Errorf("Attribute(%q).IsValid() = %v, want %v", tt.attr, got, tt.valid)
			}
		})
	}
}

func TestAllSpecialties(t *testing.T) {
	specs := AllSpecialties()
	if len(specs) != 6 {
		t.Fatalf("AllSpecialties() returned %d specialties, want 6", len(specs))
	}

	specs[0] = Specialty("corrupted")
	freshSpecs := AllSpecialties()
	if freshSpecs[0] == Specialty("corrupted") {
		t.Errorf("AllSpecialties() is vulnerable to external slice mutation")
	}
}

func TestSpecialty_IsValid(t *testing.T) {
	tests := []struct {
		spec  Specialty
		valid bool
	}{
		{SpecialtyAttack, true},
		{SpecialtyStun, true},
		{SpecialtyAnomaly, true},
		{SpecialtySupport, true},
		{SpecialtyDefense, true},
		{SpecialtyRupture, true},
		{Specialty("Healer"), false},
		{Specialty(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.spec), func(t *testing.T) {
			if got := tt.spec.IsValid(); got != tt.valid {
				t.Errorf("Specialty(%q).IsValid() = %v, want %v", tt.spec, got, tt.valid)
			}
		})
	}
}

func TestAllRarities(t *testing.T) {
	rarities := AllRarities()
	if len(rarities) != 3 {
		t.Fatalf("AllRarities() returned %d rarities, want 3", len(rarities))
	}

	rarities[0] = Rarity("corrupted")
	freshRarities := AllRarities()
	if freshRarities[0] == Rarity("corrupted") {
		t.Errorf("AllRarities() is vulnerable to external slice mutation")
	}
}

func TestRarity_IsValid(t *testing.T) {
	tests := []struct {
		rarity Rarity
		valid  bool
	}{
		{RarityS, true},
		{RarityA, true},
		{RarityB, true},
		{Rarity("SSR"), false},
		{Rarity("C"), false},
		{Rarity(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.rarity), func(t *testing.T) {
			if got := tt.rarity.IsValid(); got != tt.valid {
				t.Errorf("Rarity(%q).IsValid() = %v, want %v", tt.rarity, got, tt.valid)
			}
		})
	}
}

func TestAllSkillTypes(t *testing.T) {
	types := AllSkillTypes()
	if len(types) != 6 {
		t.Fatalf("AllSkillTypes() returned %d skill types, want 6", len(types))
	}

	types[0] = SkillType("corrupted")
	freshTypes := AllSkillTypes()
	if freshTypes[0] == SkillType("corrupted") {
		t.Errorf("AllSkillTypes() is vulnerable to external slice mutation")
	}
}

func TestSkillType_IsValid(t *testing.T) {
	tests := []struct {
		skillType SkillType
		valid     bool
	}{
		{SkillTypeBasic, true},
		{SkillTypeDodge, true},
		{SkillTypeAssist, true},
		{SkillTypeSpecial, true},
		{SkillTypeChain, true},
		{SkillTypePassive, true},
		{SkillType("ultimate"), false},
		{SkillType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.skillType), func(t *testing.T) {
			if got := tt.skillType.IsValid(); got != tt.valid {
				t.Errorf("SkillType(%q).IsValid() = %v, want %v", tt.skillType, got, tt.valid)
			}
		})
	}
}
