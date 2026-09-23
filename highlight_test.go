package fairy_test

import (
	"encoding/json"
	"testing"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/fairy"
)

func TestAgent_HighlightProps_JSON(t *testing.T) {
	// 1. Agent with highlight props
	rawAvatar := &zzz.AvatarData{
		ID: 1011, // Anby (Stun -> Impact 12201)
	}
	agent, err := fairy.EnrichAgent(rawAvatar)
	if err != nil {
		t.Fatalf("EnrichAgent failed: %v", err)
	}

	if len(agent.HighlightProps) == 0 {
		t.Fatalf("expected Anby to have HighlightProps, got empty")
	}
	if agent.HighlightProps[0] != fairy.PropBaseImpact {
		t.Errorf("expected Anby HighlightProps[0] to be %v (12201), got %v", fairy.PropBaseImpact, agent.HighlightProps[0])
	}

	data, err := json.Marshal(agent)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	if string(data) == "" {
		t.Fatalf("empty JSON output")
	}

	// 2. Empty HighlightProps and RecommendedSubStats should serialize as [] and not null
	emptyAgent := &fairy.Agent{
		HighlightProps:      make([]fairy.PropertyID, 0),
		RecommendedSubStats: make([]fairy.PropertyID, 0),
	}
	emptyData, err := json.Marshal(emptyAgent)
	if err != nil {
		t.Fatalf("json.Marshal(emptyAgent) failed: %v", err)
	}
	type jsonCheck struct {
		HighlightProps      []fairy.PropertyID `json:"highlight_props"`
		RecommendedSubStats []fairy.PropertyID `json:"recommended_sub_stats"`
	}
	var check jsonCheck
	if err := json.Unmarshal(emptyData, &check); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if check.HighlightProps == nil {
		t.Errorf("expected non-nil slice in JSON, got nil (rendered as null)")
	}
	if check.RecommendedSubStats == nil {
		t.Errorf("expected non-nil RecommendedSubStats slice in JSON, got nil (rendered as null)")
	}
}

func TestAgent_IsHighlightProp(t *testing.T) {
	// Nil receiver safety
	var nilAgent *fairy.Agent
	if nilAgent.IsHighlightProp(fairy.PropBaseATK) {
		t.Errorf("expected nil Agent.IsHighlightProp to return false")
	}

	// Empty highlight props safety
	emptyAgent := &fairy.Agent{HighlightProps: []fairy.PropertyID{}}
	if emptyAgent.IsHighlightProp(fairy.PropBaseATK) {
		t.Errorf("expected empty Agent.IsHighlightProp to return false")
	}

	// Attack Agent: Soldier 11 (ID 1041) has Base ATK (12101), Crit Rate (20101), Crit DMG (21101)
	rawAvatar := &zzz.AvatarData{ID: 1041}
	agent, err := fairy.EnrichAgent(rawAvatar)
	if err != nil {
		t.Fatalf("EnrichAgent failed: %v", err)
	}

	// Substats matching ATK
	if !agent.IsHighlightProp(fairy.PropATKPercent) {
		t.Errorf("expected ATK%% (12102) to be highlight prop for Soldier 11")
	}
	if !agent.IsHighlightProp(fairy.PropATKFlat) {
		t.Errorf("expected Flat ATK (12103) to be highlight prop for Soldier 11")
	}
	if !agent.IsHighlightProp(fairy.PropCritRate) {
		t.Errorf("expected CRIT Rate (20103) to be highlight prop for Soldier 11")
	}
	if !agent.IsHighlightProp(fairy.PropCritDMG) {
		t.Errorf("expected CRIT DMG (21103) to be highlight prop for Soldier 11")
	}

	// Non-matching stats
	if agent.IsHighlightProp(fairy.PropHPPercent) {
		t.Errorf("did not expect HP%% to be highlight prop for Soldier 11")
	}
	if agent.IsHighlightProp(fairy.PropDEFPercent) {
		t.Errorf("did not expect DEF%% to be highlight prop for Soldier 11")
	}

	// Rina (ID 1211) has PropBasePENRatio (23101)
	rinaRaw := &zzz.AvatarData{ID: 1211}
	rina, err := fairy.EnrichAgent(rinaRaw)
	if err != nil {
		t.Fatalf("EnrichAgent Rina failed: %v", err)
	}

	// Verify PEN Ratio matches for Rina, but Flat PEN does NOT match
	if !rina.IsHighlightProp(fairy.PropBasePENRatio) {
		t.Errorf("expected PropBasePENRatio to match for Rina")
	}
	if !rina.IsHighlightProp(fairy.PropPENRatio) {
		t.Errorf("expected PropPENRatio to match for Rina")
	}
	if rina.IsHighlightProp(fairy.PropBasePENFlat) {
		t.Errorf("expected PropBasePENFlat NOT to match for Rina")
	}
	if rina.IsHighlightProp(fairy.PropPENFlat) {
		t.Errorf("expected PropPENFlat NOT to match for Rina")
	}
}

func TestAgent_CountEffectiveRolls(t *testing.T) {
	// Nil safety
	var nilAgent *fairy.Agent
	if nilAgent.CountEffectiveRolls() != 0 {
		t.Errorf("expected nil Agent.CountEffectiveRolls() to return 0")
	}

	// Construct agent with equipped discs
	agent := &fairy.Agent{
		HighlightProps: []fairy.PropertyID{
			fairy.PropBaseATK,      // 12101
			fairy.PropBaseCritRate, // 20101
			fairy.PropBaseCritDMG,  // 21101
		},
		DriveDiscs: fairy.DriveDiscs{
			Slots: []fairy.DriveDisc{
				{
					Slot: 1,
					SubStats: []fairy.StatValue{
						{PropertyID: fairy.PropATKPercent, Rolls: 2}, // highlight (+2)
						{PropertyID: fairy.PropCritRate, Rolls: 3},   // highlight (+3)
						{PropertyID: fairy.PropHPPercent, Rolls: 1},  // not highlight
					},
				},
				{
					Slot: 2,
					SubStats: []fairy.StatValue{
						{PropertyID: fairy.PropCritDMG, Rolls: 4},    // highlight (+4)
						{PropertyID: fairy.PropDEFPercent, Rolls: 1}, // not highlight
					},
				},
			},
		},
	}

	// 1. Default evaluation against agent.HighlightProps: 2 + 3 + 4 = 9 rolls
	effective := agent.CountEffectiveRolls()
	if effective != 9 {
		t.Errorf("CountEffectiveRolls() = %d; want 9", effective)
	}

	// 2. Custom evaluation against only Crit Rate: 3 rolls
	critOnly := agent.CountEffectiveRolls(fairy.PropCritRate)
	if critOnly != 3 {
		t.Errorf("CountEffectiveRolls(PropCritRate) = %d; want 3", critOnly)
	}

	// 3. Custom evaluation against Crit Rate & Crit DMG: 3 + 4 = 7 rolls
	critsOnly := agent.CountEffectiveRolls(fairy.PropCritRate, fairy.PropCritDMG)
	if critsOnly != 7 {
		t.Errorf("CountEffectiveRolls(PropCritRate, PropCritDMG) = %d; want 7", critsOnly)
	}
}

func TestGlobal_AgentHighlightProps(t *testing.T) {
	// Anby (1011)
	props := fairy.AgentHighlightProps(1011)
	if len(props) != 1 || props[0] != fairy.PropBaseImpact {
		t.Errorf("AgentHighlightProps(1011) = %v; want [%v]", props, fairy.PropBaseImpact)
	}

	// Defensive copy verification
	props[0] = fairy.PropertyID(99999)
	propsAgain := fairy.AgentHighlightProps(1011)
	if propsAgain[0] != fairy.PropBaseImpact {
		t.Errorf("AgentHighlightProps returned a mutable internal slice!")
	}

	// Unknown agent ID returns nil
	unknown := fairy.AgentHighlightProps(99999999)
	if unknown != nil {
		t.Errorf("AgentHighlightProps(unknown) = %v; want nil", unknown)
	}
}

func TestGlobal_AllAgentHighlightProps(t *testing.T) {
	allProps := fairy.AllAgentHighlightProps()
	if len(allProps) < 40 {
		t.Fatalf("expected at least 40 agents in AllAgentHighlightProps, got %d", len(allProps))
	}

	// Verify key agents are present
	anbyProps, ok := allProps[1011]
	if !ok || len(anbyProps) == 0 {
		t.Errorf("missing Anby (1011) in AllAgentHighlightProps")
	}

	rinaProps, ok := allProps[1211]
	if !ok || len(rinaProps) == 0 || rinaProps[0] != fairy.PropBasePENRatio {
		t.Errorf("missing or incorrect Rina (1211) in AllAgentHighlightProps: %v", rinaProps)
	}
}
