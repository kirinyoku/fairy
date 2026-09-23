package fairy_test

import (
	"testing"

	"github.com/kirinyoku/fairy"
)

func TestAgent_RecommendedSubStats(t *testing.T) {
	remielle := &fairy.Agent{
		ID:                  1581,
		RecommendedSubStats: fairy.AgentRecommendedSubStats(1581),
	}

	nangong := &fairy.Agent{
		ID:                  1511,
		RecommendedSubStats: fairy.AgentRecommendedSubStats(1511),
	}

	rina := &fairy.Agent{
		ID:                  1211,
		RecommendedSubStats: fairy.AgentRecommendedSubStats(1211),
	}

	ellen := &fairy.Agent{
		ID:                  1191,
		RecommendedSubStats: fairy.AgentRecommendedSubStats(1191),
	}

	// 1. Remielle: Flat ATK IS recommended, ATK% is recommended, Anomaly Proficiency is recommended
	if !remielle.IsRecommendedSubStat(fairy.PropATKFlat) {
		t.Errorf("expected Remielle to have PropATKFlat recommended")
	}
	if !remielle.IsRecommendedSubStat(fairy.PropATKPercent) {
		t.Errorf("expected Remielle to have PropATKPercent recommended")
	}
	if !remielle.IsRecommendedSubStat(fairy.PropAnomalyProficiency) {
		t.Errorf("expected Remielle to have PropAnomalyProficiency recommended")
	}
	if remielle.IsRecommendedSubStat(fairy.PropDEFPercent) {
		t.Errorf("expected Remielle to NOT have PropDEFPercent recommended")
	}

	// 2. Nangong Yu: ATK% is recommended, Anomaly Proficiency is recommended, but Flat ATK is NOT!
	if nangong.IsRecommendedSubStat(fairy.PropATKFlat) {
		t.Errorf("expected Nangong Yu to NOT have PropATKFlat recommended")
	}
	if !nangong.IsRecommendedSubStat(fairy.PropATKPercent) {
		t.Errorf("expected Nangong Yu to have PropATKPercent recommended")
	}
	if !nangong.IsRecommendedSubStat(fairy.PropAnomalyProficiency) {
		t.Errorf("expected Nangong Yu to have PropAnomalyProficiency recommended")
	}

	// 3. Rina: Anomaly Proficiency and ATK% are recommended, but Flat PEN and Flat ATK are NOT!
	if rina.IsRecommendedSubStat(fairy.PropPENFlat) || rina.IsRecommendedSubStat(fairy.PropBasePENFlat) {
		t.Errorf("expected Rina to NOT have Flat PEN recommended")
	}
	if !rina.IsRecommendedSubStat(fairy.PropAnomalyProficiency) {
		t.Errorf("expected Rina to have PropAnomalyProficiency recommended")
	}
	if !rina.IsRecommendedSubStat(fairy.PropATKPercent) {
		t.Errorf("expected Rina to have PropATKPercent recommended")
	}
	if rina.IsRecommendedSubStat(fairy.PropATKFlat) {
		t.Errorf("expected Rina to NOT have PropATKFlat recommended")
	}

	// 4. Ellen: Crit Rate, Crit DMG, ATK% are recommended, but Flat ATK is NOT!
	if ellen.IsRecommendedSubStat(fairy.PropATKFlat) {
		t.Errorf("expected Ellen to NOT have PropATKFlat recommended")
	}
	if !ellen.IsRecommendedSubStat(fairy.PropATKPercent) {
		t.Errorf("expected Ellen to have PropATKPercent recommended")
	}
	if !ellen.IsRecommendedSubStat(fairy.PropCritRate) {
		t.Errorf("expected Ellen to have CRIT Rate recommended")
	}
	if !ellen.IsRecommendedSubStat(fairy.PropCritDMG) {
		t.Errorf("expected Ellen to have CRIT DMG recommended")
	}
	// Also verify base property IDs match via family in IsRecommendedSubStat:
	if !ellen.IsRecommendedSubStat(fairy.PropBaseCritRate) {
		t.Errorf("expected Ellen to match PropBaseCritRate via family")
	}
	if !ellen.IsRecommendedSubStat(fairy.PropBaseCritDMG) {
		t.Errorf("expected Ellen to match PropBaseCritDMG via family")
	}
}

func TestAgent_CountEffectiveRolls_Discs(t *testing.T) {
	// A disc with 3 rolls in ATK% (12102) and 2 rolls in Flat ATK (12103)
	disc := fairy.DriveDisc{
		Slot: 1,
		SubStats: []fairy.StatValue{
			{PropertyID: fairy.PropATKPercent, Value: 0.12, Rolls: 3},
			{PropertyID: fairy.PropATKFlat, Value: 38, Rolls: 2},
			{PropertyID: fairy.PropDEFPercent, Value: 0.05, Rolls: 1},
		},
	}

	// Remielle considers BOTH Flat ATK and ATK% useful: 3 + 2 = 5
	remielle := &fairy.Agent{
		ID:                  1581,
		RecommendedSubStats: fairy.AgentRecommendedSubStats(1581),
		DriveDiscs: fairy.DriveDiscs{
			Slots: []fairy.DriveDisc{disc},
		},
	}
	if rolls := remielle.CountEffectiveRolls(); rolls != 5 {
		t.Errorf("expected Remielle effective rolls to be 5, got %d", rolls)
	}

	// Nangong Yu considers ONLY ATK% useful: 3 rolls (Flat ATK is ignored)
	nangong := &fairy.Agent{
		ID:                  1511,
		RecommendedSubStats: fairy.AgentRecommendedSubStats(1511),
		DriveDiscs: fairy.DriveDiscs{
			Slots: []fairy.DriveDisc{disc},
		},
	}
	if rolls := nangong.CountEffectiveRolls(); rolls != 3 {
		t.Errorf("expected Nangong Yu effective rolls to be 3, got %d", rolls)
	}

	// Ellen with CRIT Rate (20103) and CRIT DMG (21103) rolls on DriveDiscs:
	critDisc := fairy.DriveDisc{
		Slot: 2,
		SubStats: []fairy.StatValue{
			{PropertyID: fairy.PropCritRate, Rolls: 3},
			{PropertyID: fairy.PropCritDMG, Rolls: 2},
			{PropertyID: fairy.PropHPPercent, Rolls: 1},
		},
	}
	ellen := &fairy.Agent{
		ID:                  1191,
		RecommendedSubStats: fairy.AgentRecommendedSubStats(1191),
		DriveDiscs: fairy.DriveDiscs{
			Slots: []fairy.DriveDisc{critDisc},
		},
	}
	// Direct agent evaluation:
	if rolls := ellen.CountEffectiveRolls(); rolls != 5 {
		t.Errorf("expected Ellen CountEffectiveRolls to be 5, got %d", rolls)
	}
	// DriveDiscs.CountEffectiveRolls with RecommendedSubStats directly:
	if rolls := ellen.DriveDiscs.CountEffectiveRolls(ellen.RecommendedSubStats...); rolls != 5 {
		t.Errorf("expected DriveDiscs.CountEffectiveRolls(RecommendedSubStats...) to be 5, got %d", rolls)
	}
}

func TestAllAgentRecommendedSubStats(t *testing.T) {
	all := fairy.AllAgentRecommendedSubStats()
	if len(all) == 0 {
		t.Fatal("expected non-empty map from AllAgentRecommendedSubStats")
	}
	if len(all) != 60 {
		t.Errorf("expected exactly 60 agents in AllAgentRecommendedSubStats, got %d", len(all))
	}
	for id, stats := range all {
		if len(stats) == 0 {
			t.Errorf("agent %d has empty recommended substats", id)
		}
	}
}

func TestNewAgents_RecommendedSubStats(t *testing.T) {
	tests := []struct {
		name     string
		agentID  int
		expected []fairy.PropertyID
	}{
		{
			name:     "Caesar",
			agentID:  1071,
			expected: []fairy.PropertyID{fairy.PropATKPercent, fairy.PropCritRate, fairy.PropCritDMG},
		},
		{
			name:     "Qingyi",
			agentID:  1251,
			expected: []fairy.PropertyID{fairy.PropATKPercent, fairy.PropCritRate, fairy.PropCritDMG},
		},
		{
			name:     "Norma",
			agentID:  1571,
			expected: []fairy.PropertyID{fairy.PropATKPercent, fairy.PropCritRate, fairy.PropCritDMG},
		},
		{
			name:     "Claret",
			agentID:  1611,
			expected: []fairy.PropertyID{fairy.PropDEFPercent, fairy.PropCritRate, fairy.PropCritDMG},
		},
		{
			name:     "Roxy",
			agentID:  1621,
			expected: []fairy.PropertyID{fairy.PropATKPercent, fairy.PropCritRate, fairy.PropCritDMG},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := fairy.AgentRecommendedSubStats(tc.agentID)
			if len(got) != len(tc.expected) {
				t.Fatalf("%s: expected %d recommended substats, got %d: %v", tc.name, len(tc.expected), len(got), got)
			}
			for i, prop := range tc.expected {
				if got[i] != prop {
					t.Errorf("%s: expected substat at %d to be %v, got %v", tc.name, i, prop, got[i])
				}
			}

			agent := &fairy.Agent{
				ID:                  tc.agentID,
				RecommendedSubStats: got,
			}
			for _, prop := range tc.expected {
				if !agent.IsRecommendedSubStat(prop) {
					t.Errorf("%s: expected IsRecommendedSubStat(%v) to be true", tc.name, prop)
				}
			}
		})
	}
}
