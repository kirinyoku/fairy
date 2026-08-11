package fairy

// Attribute represents the elemental attribute of an agent.
type Attribute string

const (
	// AttributePhysical represents the agent's Physical attribute.
	AttributePhysical Attribute = "Physical"
	// AttributeHonedEdge represents the agent's Honed Edge attribute.
	AttributeHonedEdge Attribute = "HonedEdge"
	// AttributeFire represents the agent's Fire attribute.
	AttributeFire Attribute = "Fire"
	// AttributeIce represents the agent's Ice attribute.
	AttributeIce Attribute = "Ice"
	// AttributeFrost represents the agent's Frost attribute.
	AttributeFrost Attribute = "Frost"
	// AttributeElectric represents the agent's Electric attribute.
	AttributeElectric Attribute = "Electric"
	// AttributeEther represents the agent's Ether attribute.
	AttributeEther Attribute = "Ether"
	// AttributeAuricInk represents the agent's Auric Ink attribute.
	AttributeAuricInk Attribute = "AuricInk"
	// AttributeWind represents the agent's Wind attribute.
	AttributeWind Attribute = "Wind"
	// AttributeLumiflux represents the agent's Lumiflux attribute.
	AttributeLumiflux Attribute = "Lumiflux"
)

// Specialty represents the combat role or class of an agent.
type Specialty string

const (
	// SpecialtyAttack represents the Attack combat role.
	SpecialtyAttack Specialty = "Attack"
	// SpecialtyStun represents the Stun combat role.
	SpecialtyStun Specialty = "Stun"
	// SpecialtyAnomaly represents the Anomaly combat role.
	SpecialtyAnomaly Specialty = "Anomaly"
	// SpecialtySupport represents the Support combat role.
	SpecialtySupport Specialty = "Support"
	// SpecialtyDefense represents the Defense combat role.
	SpecialtyDefense Specialty = "Defense"
	// SpecialtyRupture represents the Rupture combat role.
	SpecialtyRupture Specialty = "Rupture"
)

// Rarity represents the rarity tier of agents and equipment.
type Rarity string

const (
	// RarityS represents the S-rank tier.
	RarityS Rarity = "S"
	// RarityA represents the A-rank tier.
	RarityA Rarity = "A"
	// RarityB represents the B-rank tier.
	RarityB Rarity = "B"
)

// Skin represents the equipped skin (outfit) of an agent.
type Skin struct {
	ID           int    `json:"id"`             // The internal ID of the skin.
	Name         string `json:"name"`           // The localized name of the skin.
	Description  string `json:"description"`    // The localized description of the skin.
	SplashArtURL string `json:"splash_art_url"` // The URL to the skin's splash art.
}

// Skill represents an agent's combat skill or passive ability.
type Skill struct {
	Level       int    `json:"level"`       // The level of the skill.
	Name        string `json:"name"`        // The localized name of the skill.
	Description string `json:"description"` // The localized description of the skill.
}

// Agent represents an enriched agent (character) showcased on a player's profile.
// A profile can showcase a maximum of 6 agents. It contains the agent's combat
// metadata, equipped gear, and final stats.
type Agent struct {
	ID                   int                 `json:"id"`                     // The internal ID of the agent.
	Name                 string              `json:"name"`                   // The localized name of the agent (e.g., "Ellen").
	Level                int                 `json:"level"`                  // The current level of the agent (1-60).
	Promotion            int                 `json:"promotion"`              // The promotion/ascension phase of the agent (0-5).
	MindscapeCinema      int                 `json:"mindscape_cinema"`       // The unlocked Mindscape Cinema level (0-6).
	CoreSkillEnhancement int                 `json:"core_skill_enhancement"` // The Core Skill enhancement level (0-6).
	Attribute            Attribute           `json:"attribute"`              // The elemental damage type (e.g., Ice).
	AttributeName        string              `json:"attribute_name"`         // The localized name of the attribute.
	Specialty            Specialty           `json:"specialty"`              // The combat role (e.g., Attack).
	SpecialtyName        string              `json:"specialty_name"`         // The localized name of the specialty.
	Rarity               Rarity              `json:"rarity"`                 // The rarity tier (S or A).
	Skin                 *Skin               `json:"skin"`                   // The currently equipped skin (can be nil if not found).
	SplashArtURL         string              `json:"splash_art_url"`         // The URL to the agent's splash art.
	Skills               []Skill             `json:"skills"`                 // The agent's skills and passives.
	WEngine              *WEngine            `json:"w_engine"`               // The currently equipped W-Engine (can be nil).
	DriveDiscs           []DriveDisc         `json:"drive_discs"`            // The equipped Drive Discs (up to 6).
	ActiveSetBonuses     []DriveDiscSetBonus `json:"active_set_bonuses"`     // The active 2-piece or 4-piece set bonuses.
	BaseStats            Stats               `json:"base_stats"`             // The agent's base combat stats before gear/buffs.
	Stats                Stats               `json:"stats"`                  // The agent's final combat stats including all gear/buffs.
}

// SubStatTotals calculates the sum of all sub-stats across all equipped Drive Discs.
// It groups them by PropertyID and sums the Rolls and Values.
// The returned slice is guaranteed to preserve the initial appearance order of sub-stats.
func (a *Agent) SubStatTotals() []StatValue {
	totals := make(map[PropertyID]StatValue)
	var order []PropertyID // Keep track of the order to ensure deterministic output

	for _, disc := range a.DriveDiscs {
		for _, sub := range disc.SubStats {
			if curr, exists := totals[sub.PropertyID]; exists {
				curr.Value += sub.Value
				curr.Rolls += sub.Rolls
				totals[sub.PropertyID] = curr
			} else {
				totals[sub.PropertyID] = sub
				order = append(order, sub.PropertyID)
			}
		}
	}

	result := make([]StatValue, 0, len(order))
	for _, id := range order {
		result = append(result, totals[id])
	}
	return result
}

// CountEffectiveRolls returns the total number of sub-stat rolls across all Drive Discs
// that match any of the provided target property IDs (also known as "effective" or "useful" rolls).
func (a *Agent) CountEffectiveRolls(targetProps ...PropertyID) int {
	total := 0
	targetMap := make(map[PropertyID]bool)
	for _, p := range targetProps {
		targetMap[p] = true
	}

	for _, disc := range a.DriveDiscs {
		for _, sub := range disc.SubStats {
			if targetMap[sub.PropertyID] {
				total += sub.Rolls
			}
		}
	}
	return total
}

// FormattedUIStats generates a complete breakdown of base vs added stats for UI display.
// This structure precisely matches the visual representation and layout seen in the in-game
// stat panel or on platforms like Enka.Network.
func (a *Agent) FormattedUIStats() UIStats {
	return UIStats{
		HP:                 formatBreakdown(a.BaseStats.HP, a.Stats.HP, false, 0, 1),
		ATK:                formatBreakdown(a.BaseStats.ATK, a.Stats.ATK, false, 0, 1),
		DEF:                formatBreakdown(a.BaseStats.DEF, a.Stats.DEF, false, 0, 1),
		Impact:             formatBreakdown(a.BaseStats.Impact, a.Stats.Impact, false, 0, 1),
		CritRate:           formatBreakdown(a.BaseStats.CritRate, a.Stats.CritRate, true, 0, 1),
		CritDMG:            formatBreakdown(a.BaseStats.CritDMG, a.Stats.CritDMG, true, 0, 1),
		AnomalyMastery:     formatBreakdown(a.BaseStats.AnomalyMastery, a.Stats.AnomalyMastery, false, 0, 1),
		AnomalyProficiency: formatBreakdown(a.BaseStats.AnomalyProficiency, a.Stats.AnomalyProficiency, false, 0, 1),
		PenRatio:           formatBreakdown(a.BaseStats.PenRatio, a.Stats.PenRatio, true, 0, 1),
		PenFlat:            formatBreakdown(a.BaseStats.PenFlat, a.Stats.PenFlat, false, 0, 1),
		EnergyRegen:        formatBreakdown(a.BaseStats.EnergyRegen, a.Stats.EnergyRegen, false, 2, 2),
		SheerForce:         formatBreakdown(a.BaseStats.SheerForce, a.Stats.SheerForce, false, 0, 1),
	}
}
