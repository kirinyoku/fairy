package fairy

import (
	"encoding/base64"

	"github.com/kirinyoku/fairy/internal/store"
)

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

var allAttributes = [...]Attribute{
	AttributePhysical,
	AttributeHonedEdge,
	AttributeFire,
	AttributeIce,
	AttributeFrost,
	AttributeElectric,
	AttributeEther,
	AttributeAuricInk,
	AttributeWind,
	AttributeLumiflux,
}

// AllAttributes returns a slice containing all agent elemental attributes.
func AllAttributes() []Attribute {
	attrs := make([]Attribute, len(allAttributes))
	copy(attrs, allAttributes[:])
	return attrs
}

// IsValid reports whether the attribute is a recognized agent elemental attribute.
func (a Attribute) IsValid() bool {
	for _, attr := range allAttributes {
		if a == attr {
			return true
		}
	}
	return false
}

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

var allSpecialties = [...]Specialty{
	SpecialtyAttack,
	SpecialtyStun,
	SpecialtyAnomaly,
	SpecialtySupport,
	SpecialtyDefense,
	SpecialtyRupture,
}

// AllSpecialties returns a slice containing all combat specialties/roles.
func AllSpecialties() []Specialty {
	specs := make([]Specialty, len(allSpecialties))
	copy(specs, allSpecialties[:])
	return specs
}

// IsValid reports whether the specialty is a recognized combat role.
func (s Specialty) IsValid() bool {
	for _, spec := range allSpecialties {
		if s == spec {
			return true
		}
	}
	return false
}

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

var allRarities = [...]Rarity{
	RarityS,
	RarityA,
	RarityB,
}

// AllRarities returns a slice containing all rarity tiers.
func AllRarities() []Rarity {
	rarities := make([]Rarity, len(allRarities))
	copy(rarities, allRarities[:])
	return rarities
}

// IsValid reports whether the rarity is a recognized rarity tier.
func (r Rarity) IsValid() bool {
	for _, rarity := range allRarities {
		if r == rarity {
			return true
		}
	}
	return false
}

// BaseAttribute returns the core elemental attribute that this attribute deals damage as.
// For example, AuricInk deals Ether DMG, HonedEdge deals Physical DMG, and Frost deals Ice DMG.
func (a Attribute) BaseAttribute() Attribute {
	switch a {
	case AttributeAuricInk:
		return AttributeEther
	case AttributeHonedEdge:
		return AttributePhysical
	case AttributeFrost:
		return AttributeIce
	default:
		return a
	}
}

// SVG returns the raw inline SVG markup string for the attribute.
func (a Attribute) SVG() string {
	return attributeSVGMap[a]
}

// IconURL returns the base64-encoded Data URI string containing the attribute's SVG icon.
func (a Attribute) IconURL() string {
	svg := a.SVG()
	if svg == "" {
		return ""
	}
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg))
}

// IconURL returns the official Enka CDN icon URL for the specialty.
func (s Specialty) IconURL() string {
	switch s {
	case SpecialtyAttack:
		return EnkaAssetBaseURL + "IconAttack.png"
	case SpecialtyStun:
		return EnkaAssetBaseURL + "IconStun.png"
	case SpecialtyAnomaly:
		return EnkaAssetBaseURL + "IconAnomaly.png"
	case SpecialtySupport:
		return EnkaAssetBaseURL + "IconSupport.png"
	case SpecialtyDefense:
		return EnkaAssetBaseURL + "IconDefense.png"
	case SpecialtyRupture:
		return EnkaAssetBaseURL + "IconRupture.png"
	default:
		return ""
	}
}

// IconURL returns the official Enka CDN icon URL for the rarity tier.
func (r Rarity) IconURL() string {
	switch r {
	case RarityS:
		return EnkaAssetBaseURL + "ItemRarityS.png"
	case RarityA:
		return EnkaAssetBaseURL + "ItemRarityA.png"
	case RarityB:
		return EnkaAssetBaseURL + "ItemRarityB.png"
	default:
		return ""
	}
}

// Skin represents the equipped skin (outfit) of an agent.
type Skin struct {
	ID           int    `json:"id"`             // The internal ID of the skin.
	Name         string `json:"name"`           // The localized name of the skin.
	Description  string `json:"description"`    // The localized description of the skin.
	SplashArtURL string `json:"splash_art_url"` // The URL to the skin's splash art.
}

// SkillType represents a categorized category of combat skill.
type SkillType string

const (
	SkillTypeBasic   SkillType = "basic"   // Basic Attack
	SkillTypeDodge   SkillType = "dodge"   // Dodge & Counter
	SkillTypeAssist  SkillType = "assist"  // Quick & Defensive Assists
	SkillTypeSpecial SkillType = "special" // Special & EX Special Attack
	SkillTypeChain   SkillType = "chain"   // Chain Attack & Ultimate
	SkillTypePassive SkillType = "passive" // Core Passive & Additional Ability
)

var allSkillTypes = [...]SkillType{
	SkillTypeBasic,
	SkillTypeDodge,
	SkillTypeAssist,
	SkillTypeSpecial,
	SkillTypeChain,
	SkillTypePassive,
}

// AllSkillTypes returns a slice containing all combat skill categories.
func AllSkillTypes() []SkillType {
	types := make([]SkillType, len(allSkillTypes))
	copy(types, allSkillTypes[:])
	return types
}

// IsValid reports whether the skill type is a recognized skill category.
func (st SkillType) IsValid() bool {
	for _, t := range allSkillTypes {
		if st == t {
			return true
		}
	}
	return false
}

// SkillParam represents a calculated numeric parameter or multiplier for a skill.
type SkillParam struct {
	Name  string `json:"name"`  // Localized parameter name.
	Value string `json:"value"` // Formatted value.
}

// Skill represents an agent's combat skill or passive ability.
type Skill struct {
	Level         int          `json:"level"`                    // The level of the skill.
	Name          string       `json:"name"`                     // The localized name of the skill.
	Description   string       `json:"description"`              // The localized description of the skill.
	FormattedHTML string       `json:"formatted_html,omitempty"` // Formatted HTML description with inline colors and evaluated formulas.
	Type          SkillType    `json:"type"`                     // Category type of the skill (basic, dodge, assist, special, chain, passive).
	TypeName      string       `json:"type_name"`                // Localized category type name.
	Params        []SkillParam `json:"params,omitempty"`         // Numeric parameters / multiplier table.
}

// SkillGroup represents a categorized group of skills matching the in-game UI / Enka buttons.
type SkillGroup struct {
	Type     SkillType `json:"type"`      // Group category key ("basic", "special", "dodge", "chain", "assist", "passive").
	TypeName string    `json:"type_name"` // Localized category group name.
	Level    int       `json:"level"`     // Group level (1-12 for active skills, 0-6 for core passives).
	Skills   []Skill   `json:"skills"`    // Individual skills belonging to this category group.
}

// EvaluatedDescription returns the skill description with all scaling formulas ({CAL:...})
// evaluated for the skill's current level.
func (s Skill) EvaluatedDescription() string {
	return evaluateFormulas(s.Description, s.Level)
}

// FormatHTML returns the skill description formatted as HTML with inline CSS colors,
// semantic icon spans, and scaling formulas evaluated for the skill's current level.
func (s Skill) FormatHTML() string {
	return formatHTML(s.Description, s.Level)
}

// FormatPlainText returns the skill description as clean plain text with all tags stripped
// and scaling formulas evaluated for the skill's current level.
func (s Skill) FormatPlainText() string {
	return formatPlainText(s.Description, s.Level)
}

// FormatMarkdown returns the skill description formatted in Markdown (bold tags for colored values)
// with scaling formulas evaluated for the skill's current level.
func (s Skill) FormatMarkdown() string {
	return formatMarkdown(s.Description, s.Level)
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
	SkillGroups          []SkillGroup        `json:"grouped_skills"`         // Skills categorized into 6 UI groups (Passives, Basic, Special, Dodge, Chain, Assist).
	Mindscapes           []MindscapeNode     `json:"mindscapes"`             // The agent's Mindscape Cinema levels (1-6).
	PotentialVision      *PotentialVision    `json:"potential_vision"`       // Potential Vision upgrade mechanics (can be nil if agent has none).
	WEngine              *WEngine            `json:"w_engine"`               // The currently equipped W-Engine (can be nil).
	DriveDiscs           []DriveDisc         `json:"drive_discs"`            // The equipped Drive Discs (up to 6).
	ActiveSetBonuses     []DriveDiscSetBonus `json:"active_set_bonuses"`     // The active 2-piece or 4-piece set bonuses.
	BaseStats            Stats               `json:"base_stats"`             // The agent's base combat stats before gear/buffs.
	Stats                Stats               `json:"stats"`                  // The agent's final combat stats including all gear/buffs.
	UIStats              UIStats             `json:"ui_stats"`               // UI-ready formatted combat stats panel with icons and breakdowns.
}

// MindscapeNode represents a single Mindscape Cinema level (1-6) for an Agent.
type MindscapeNode struct {
	Rank          int    `json:"rank"`                     // Cinema level (1 to 6).
	Name          string `json:"name"`                     // Localized name of the Mindscape Cinema.
	Description   string `json:"description"`              // Localized description of the effect.
	FormattedHTML string `json:"formatted_html,omitempty"` // Formatted HTML description with inline colors.
	Unlocked      bool   `json:"unlocked"`                 // True if unlocked (MindscapeCinema >= Rank).
}

// FormatHTML returns the Mindscape description formatted as HTML with inline CSS colors.
func (m MindscapeNode) FormatHTML() string {
	return formatHTML(m.Description)
}

// FormatPlainText returns the Mindscape description stripped of Rich Text formatting.
func (m MindscapeNode) FormatPlainText() string {
	return formatPlainText(m.Description)
}

// FormatMarkdown returns the Mindscape description formatted with Markdown syntax.
func (m MindscapeNode) FormatMarkdown() string {
	return formatMarkdown(m.Description)
}

// PotentialVision represents Potential Vision status and nodes for an Agent.
type PotentialVision struct {
	IsUnlocked bool                  `json:"is_unlocked"` // True if Potential Vision mechanic is unlocked.
	CurrentID  int                   `json:"current_id"`  // Current active Upgrade ID.
	Nodes      []PotentialVisionNode `json:"nodes"`       // All potential vision upgrade nodes.
}

// PotentialVisionNode represents a single Potential Vision upgrade node.
type PotentialVisionNode struct {
	ID            int    `json:"id"`                       // Upgrade node ID.
	Level         int    `json:"level"`                    // Level threshold (1 to 6).
	LevelName     string `json:"level_name"`               // Localized level title.
	Title         string `json:"title"`                    // Localized title.
	Description   string `json:"description"`              // Localized effect description.
	FormattedHTML string `json:"formatted_html,omitempty"` // Formatted HTML description with inline colors.
	IsActive      bool   `json:"is_active"`                // True if this node is active on the agent.
}

// FormatHTML returns the PotentialVisionNode description formatted as HTML with inline CSS colors.
func (p PotentialVisionNode) FormatHTML() string {
	return formatHTML(p.Description)
}

// FormatPlainText returns the PotentialVisionNode description stripped of Rich Text formatting.
func (p PotentialVisionNode) FormatPlainText() string {
	return formatPlainText(p.Description)
}

// FormatMarkdown returns the PotentialVisionNode description formatted with Markdown syntax.
func (p PotentialVisionNode) FormatMarkdown() string {
	return formatMarkdown(p.Description)
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

func getStatName(st store.MetadataStore, key string, lang Language) string {
	if st != nil {
		if val := st.Localize(key, string(lang)); val != "" && val != key {
			return val
		}
	}
	return key
}

// formatAgentUIStats generates a complete breakdown of base vs added stats for UI display
// using the provided MetadataStore and Language.
func formatAgentUIStats(a *Agent, s store.MetadataStore, lang Language) UIStats {
	if s == nil {
		return UIStats{}
	}

	hpName := getStatName(s, locKeyHP, lang)
	atkName := getStatName(s, locKeyATK, lang)
	defName := getStatName(s, locKeyDEF, lang)
	impactName := getStatName(s, locKeyImpact, lang)
	critRateName := getStatName(s, locKeyCritRate, lang)
	critDMGName := getStatName(s, locKeyCritDMG, lang)
	anomalyMasteryName := getStatName(s, locKeyAnomalyMastery, lang)
	anomalyProficiencyName := getStatName(s, locKeyAnomalyProficiency, lang)
	penRatioName := getStatName(s, locKeyPenRatio, lang)
	penFlatName := getStatName(s, locKeyPenFlat, lang)

	energyRegenKey := locKeyEnergyRegen
	energyRegenProp := PropBaseEnergyRegen
	if a.Specialty == SpecialtyRupture {
		energyRegenKey = locKeyRpRecover
		energyRegenProp = PropBaseRpRecover
	}
	energyRegenName := getStatName(s, energyRegenKey, lang)
	sheerForceName := getStatName(s, locKeySheerForce, lang)

	var attrDMGName string
	var attrDMGProp PropertyID
	switch a.Attribute.BaseAttribute() {
	case AttributePhysical:
		attrDMGName = getStatName(s, locKeyPhysicalDMGBonus, lang)
		attrDMGProp = PropPhysicalDMGBonus
	case AttributeFire:
		attrDMGName = getStatName(s, locKeyFireDMGBonus, lang)
		attrDMGProp = PropFireDMGBonus
	case AttributeIce:
		attrDMGName = getStatName(s, locKeyIceDMGBonus, lang)
		attrDMGProp = PropIceDMGBonus
	case AttributeElectric:
		attrDMGName = getStatName(s, locKeyElectricDMGBonus, lang)
		attrDMGProp = PropElectricDMGBonus
	case AttributeEther:
		attrDMGName = getStatName(s, locKeyEtherDMGBonus, lang)
		attrDMGProp = PropEtherDMGBonus
	case AttributeWind:
		attrDMGName = getStatName(s, locKeyWindDMGBonus, lang)
		attrDMGProp = PropWindDMGBonus
	}

	return UIStats{
		HP:                 formatFlatBreakdown(PropBaseHP, hpName, a.BaseStats.HP, a.Stats.HP, 0),
		ATK:                formatFlatBreakdown(PropBaseATK, atkName, a.BaseStats.ATK, a.Stats.ATK, 0),
		DEF:                formatFlatBreakdown(PropBaseDEF, defName, a.BaseStats.DEF, a.Stats.DEF, 0),
		Impact:             formatFlatBreakdown(PropBaseImpact, impactName, a.BaseStats.Impact, a.Stats.Impact, 0),
		CritRate:           formatPercentBreakdown(PropBaseCritRate, critRateName, a.BaseStats.CritRate, a.Stats.CritRate, 1),
		CritDMG:            formatPercentBreakdown(PropBaseCritDMG, critDMGName, a.BaseStats.CritDMG, a.Stats.CritDMG, 1),
		AttributeDMGBonus:  formatPercentBreakdown(attrDMGProp, attrDMGName, a.BaseStats.AttributeDMGBonus, a.Stats.AttributeDMGBonus, 1),
		AnomalyMastery:     formatFlatBreakdown(PropBaseAnomalyMastery, anomalyMasteryName, a.BaseStats.AnomalyMastery, a.Stats.AnomalyMastery, 0),
		AnomalyProficiency: formatFlatBreakdown(PropBaseAnomalyProficiency, anomalyProficiencyName, a.BaseStats.AnomalyProficiency, a.Stats.AnomalyProficiency, 0),
		PenRatio:           formatPercentBreakdown(PropBasePENRatio, penRatioName, a.BaseStats.PenRatio, a.Stats.PenRatio, 1),
		PenFlat:            formatFlatBreakdown(PropBasePENFlat, penFlatName, a.BaseStats.PenFlat, a.Stats.PenFlat, 0),
		EnergyRegen:        formatFlatBreakdown(energyRegenProp, energyRegenName, a.BaseStats.EnergyRegen, a.Stats.EnergyRegen, 2),
		SheerForce:         formatFlatBreakdown(PropBaseSheerForce, sheerForceName, a.BaseStats.SheerForce, a.Stats.SheerForce, 0),
	}
}

// groupAgentSkills categorizes an agent's skills into 6 distinct groups
// matching the in-game skill buttons:
// 1. Passives / Talents (Core Passive + Additional Ability)
// 2. Basic Attack
// 3. Dodge
// 4. Assist
// 5. Special Attack
// 6. Chain Attack & Ultimate
func groupAgentSkills(skills []Skill) []SkillGroup {
	order := []SkillType{
		SkillTypePassive,
		SkillTypeBasic,
		SkillTypeDodge,
		SkillTypeChain,
		SkillTypeAssist,
		SkillTypeSpecial,
	}

	groupsMap := make(map[SkillType]*SkillGroup)
	for _, sk := range skills {
		st := sk.Type
		if st == "" {
			st = SkillTypeBasic
		}
		grp, ok := groupsMap[st]
		if !ok {
			grp = &SkillGroup{
				Type:     st,
				TypeName: sk.TypeName,
				Level:    sk.Level,
				Skills:   make([]Skill, 0),
			}
			groupsMap[st] = grp
		}
		grp.Skills = append(grp.Skills, sk)
	}

	result := make([]SkillGroup, 0, len(order))
	for _, st := range order {
		if grp, ok := groupsMap[st]; ok {
			result = append(result, *grp)
		}
	}
	return result
}
