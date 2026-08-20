package fairy

import (
	"encoding/base64"

	"github.com/kirinyoku/fairy/internal/store"
)

// Attribute represents the elemental combat attribute (damage type) of an [Agent].
type Attribute string

// Supported elemental combat attributes in Zenless Zone Zero.
const (
	// AttributePhysical represents the Physical damage attribute.
	AttributePhysical Attribute = "Physical"
	// AttributeHonedEdge represents the Honed Edge attribute (Physical variant).
	AttributeHonedEdge Attribute = "HonedEdge"
	// AttributeFire represents the Fire damage attribute.
	AttributeFire Attribute = "Fire"
	// AttributeIce represents the Ice damage attribute.
	AttributeIce Attribute = "Ice"
	// AttributeFrost represents the Frost attribute (Ice variant).
	AttributeFrost Attribute = "Frost"
	// AttributeElectric represents the Electric damage attribute.
	AttributeElectric Attribute = "Electric"
	// AttributeEther represents the Ether damage attribute.
	AttributeEther Attribute = "Ether"
	// AttributeAuricInk represents the Auric Ink attribute (Ether variant).
	AttributeAuricInk Attribute = "AuricInk"
	// AttributeWind represents the Wind damage attribute.
	AttributeWind Attribute = "Wind"
	// AttributeLumiflux represents the Lumiflux damage attribute.
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

// AllAttributes returns a newly allocated slice containing all 10 supported [Attribute] constants.
// The returned slice is a defensive copy and can be safely mutated by the caller.
func AllAttributes() []Attribute {
	attrs := make([]Attribute, len(allAttributes))
	copy(attrs, allAttributes[:])
	return attrs
}

// IsValid reports whether the attribute is one of the recognized [Attribute] constants.
func (a Attribute) IsValid() bool {
	for _, attr := range allAttributes {
		if a == attr {
			return true
		}
	}
	return false
}

// Specialty represents the combat role or class of an [Agent].
type Specialty string

// Supported combat specialties in Zenless Zone Zero.
const (
	// SpecialtyAttack represents the Attack role.
	SpecialtyAttack Specialty = "Attack"
	// SpecialtyStun represents the Stun role.
	SpecialtyStun Specialty = "Stun"
	// SpecialtyAnomaly represents the Anomaly role.
	SpecialtyAnomaly Specialty = "Anomaly"
	// SpecialtySupport represents the Support role.
	SpecialtySupport Specialty = "Support"
	// SpecialtyDefense represents the Defense role.
	SpecialtyDefense Specialty = "Defense"
	// SpecialtyRupture represents the Rupture role.
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

// AllSpecialties returns a newly allocated slice containing all 6 supported [Specialty] constants.
// The returned slice is a defensive copy and can be safely mutated by the caller.
func AllSpecialties() []Specialty {
	specs := make([]Specialty, len(allSpecialties))
	copy(specs, allSpecialties[:])
	return specs
}

// IsValid reports whether the specialty is one of the recognized [Specialty] constants.
func (s Specialty) IsValid() bool {
	for _, spec := range allSpecialties {
		if s == spec {
			return true
		}
	}
	return false
}

// Rarity represents the rarity tier (rank) of an [Agent], [WEngine], or [DriveDisc].
type Rarity string

// Supported rarity ranks in Zenless Zone Zero.
const (
	// RarityS represents the S-Rank tier.
	RarityS Rarity = "S"
	// RarityA represents the A-Rank tier.
	RarityA Rarity = "A"
	// RarityB represents the B-Rank tier (used for W-Engines and Drive Discs).
	RarityB Rarity = "B"
)

var allRarities = [...]Rarity{
	RarityS,
	RarityA,
	RarityB,
}

// AllRarities returns a newly allocated slice containing all supported [Rarity] constants.
// The returned slice is a defensive copy and can be safely mutated by the caller.
func AllRarities() []Rarity {
	rarities := make([]Rarity, len(allRarities))
	copy(rarities, allRarities[:])
	return rarities
}

// IsValid reports whether the rarity is one of the recognized [Rarity] constants.
func (r Rarity) IsValid() bool {
	for _, rarity := range allRarities {
		if r == rarity {
			return true
		}
	}
	return false
}

// BaseAttribute returns the core elemental attribute that this attribute deals damage as.
// For example, [AttributeAuricInk] deals Ether DMG, [AttributeHonedEdge] deals Physical DMG,
// and [AttributeFrost] deals Ice DMG.
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

// SVG returns the raw inline SVG markup string for the attribute icon.
func (a Attribute) SVG() string {
	return attributeSVGMap[a]
}

// IconURL returns a base64-encoded Data URI string ("data:image/svg+xml;base64,...")
// containing the attribute's SVG icon for direct use in web frontend <img> tags.
func (a Attribute) IconURL() string {
	svg := a.SVG()
	if svg == "" {
		return ""
	}
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg))
}

// IconURL returns the official EnkaNetwork CDN icon URL for the specialty.
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

// IconURL returns the official EnkaNetwork CDN icon URL for the rarity tier.
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

// Skin represents an equipped cosmetic skin (outfit) of an [Agent].
type Skin struct {
	// ID is the internal numeric identifier of the skin.
	ID int `json:"id"`

	// Name is the localized name of the skin.
	Name string `json:"name"`

	// Description is the localized lore or description of the skin.
	Description string `json:"description"`

	// SplashArtURL is the absolute HTTPS URL pointing to the skin's splash art on the EnkaNetwork CDN.
	SplashArtURL string `json:"splash_art_url"`
}

// SkillType represents a categorized category of combat skill matching in-game button inputs.
type SkillType string

// Supported combat skill categories in Zenless Zone Zero.
const (
	// SkillTypeBasic represents Basic Attack combos.
	SkillTypeBasic SkillType = "basic"
	// SkillTypeDodge represents Dodge, Dash Attack, and Dodge Counter.
	SkillTypeDodge SkillType = "dodge"
	// SkillTypeAssist represents Quick Assist, Defensive Assist, and Evasive Assist.
	SkillTypeAssist SkillType = "assist"
	// SkillTypeSpecial represents Special Attack and EX Special Attack.
	SkillTypeSpecial SkillType = "special"
	// SkillTypeChain represents Chain Attack and Ultimate.
	SkillTypeChain SkillType = "chain"
	// SkillTypePassive represents Core Passive and Additional Ability.
	SkillTypePassive SkillType = "passive"
)

var allSkillTypes = [...]SkillType{
	SkillTypeBasic,
	SkillTypeDodge,
	SkillTypeAssist,
	SkillTypeSpecial,
	SkillTypeChain,
	SkillTypePassive,
}

// AllSkillTypes returns a newly allocated slice containing all 6 supported [SkillType] categories.
// The returned slice is a defensive copy and can be safely mutated by the caller.
func AllSkillTypes() []SkillType {
	types := make([]SkillType, len(allSkillTypes))
	copy(types, allSkillTypes[:])
	return types
}

// IsValid reports whether the skill type is one of the recognized [SkillType] categories.
func (st SkillType) IsValid() bool {
	for _, t := range allSkillTypes {
		if st == t {
			return true
		}
	}
	return false
}

// SkillParam represents a calculated numeric parameter or damage multiplier for a skill.
type SkillParam struct {
	// Name is the localized parameter label (e.g. "1-Hit DMG").
	Name string `json:"name"`

	// Value is the pre-formatted value string with level scaling evaluated (e.g. "124.5%").
	Value string `json:"value"`
}

// Skill represents an individual combat ability or passive effect of an [Agent].
type Skill struct {
	// Level is the current progression level of the skill.
	Level int `json:"level"`

	// Name is the localized display name of the skill.
	Name string `json:"name"`

	// Description is the localized description text with Unity Rich Text formatting and formula tags.
	Description string `json:"description"`

	// FormattedHTML is the web-ready HTML description with inline colors, icon tags, and evaluated formulas.
	FormattedHTML string `json:"formatted_html,omitempty"`

	// Type is the categorized [SkillType] of the skill (basic, dodge, assist, special, chain, passive).
	Type SkillType `json:"type"`

	// TypeName is the localized name of the skill category (e.g. "Basic Attack", "EX Special Attack").
	TypeName string `json:"type_name"`

	// Params is the list of calculated numeric parameters and multipliers evaluated for the skill's current level.
	Params []SkillParam `json:"params,omitempty"`
}

// SkillGroup represents a categorized group of skills matching the 6 in-game UI skill tabs/buttons.
type SkillGroup struct {
	// Type is the group category key ([SkillTypeBasic], [SkillTypeSpecial], etc.).
	Type SkillType `json:"type"`

	// TypeName is the localized category tab name.
	TypeName string `json:"type_name"`

	// Level is the progression level of the group (1–12 for active skills, 0–6 for core passives).
	Level int `json:"level"`

	// Skills is the list of individual [Skill] abilities belonging to this category group tab.
	Skills []Skill `json:"skills"`
}

// EvaluatedDescription returns the skill description with all dynamic scaling formulas ({CAL:...})
// evaluated for the skill's current level.
func (s Skill) EvaluatedDescription() string {
	return evaluateFormulas(s.Description, s.Level)
}

// FormatHTML returns the skill description formatted as HTML with inline CSS styling,
// embedded Enka CDN icon tags, and scaling formulas evaluated for the skill's current level.
func (s Skill) FormatHTML() string {
	return formatHTML(s.Description, s.Level)
}

// FormatPlainText returns the skill description as clean plain text with all Unity Rich Text tags stripped
// and scaling formulas evaluated for the skill's current level.
func (s Skill) FormatPlainText() string {
	return formatPlainText(s.Description, s.Level)
}

// FormatMarkdown returns the skill description formatted in Markdown (bold highlights for colored text)
// with scaling formulas evaluated for the skill's current level.
func (s Skill) FormatMarkdown() string {
	return formatMarkdown(s.Description, s.Level)
}

// Agent represents an enriched Agent showcased on a player's [Profile].
// A profile can showcase up to 6 agents. It aggregates combat metadata,
// equipped gear ([WEngine], [DriveDisc] entries), active set bonuses, and final combat [Stats].
type Agent struct {
	// ID is the internal numeric identifier of the Agent.
	ID int `json:"id"`

	// Name is the localized display name of the Agent (e.g. "Ellen", "Zhu Yuan", "Miyabi").
	Name string `json:"name"`

	// Level is the current level of the Agent (1–60).
	Level int `json:"level"`

	// Promotion is the current Promotion (ascension) phase of the Agent (0–5).
	Promotion int `json:"promotion"`

	// MindscapeCinema is the unlocked Mindscape Cinema level of the Agent (0–6).
	MindscapeCinema int `json:"mindscape_cinema"`

	// CoreSkillEnhancement is the Core Skill Enhancement level of the Agent (0–6 / Core A–F).
	CoreSkillEnhancement int `json:"core_skill_enhancement"`

	// Attribute is the elemental combat damage type of the Agent (e.g. [AttributeIce], [AttributeEther]).
	Attribute Attribute `json:"attribute"`

	// AttributeName is the localized display name of the Agent's elemental attribute (e.g. "Ice", "Ether").
	AttributeName string `json:"attribute_name"`

	// Specialty is the combat role of the Agent (e.g. [SpecialtyAttack], [SpecialtyStun]).
	Specialty Specialty `json:"specialty"`

	// SpecialtyName is the localized display name of the Agent's combat specialty (e.g. "Attack", "Stun").
	SpecialtyName string `json:"specialty_name"`

	// Rarity is the rarity rank of the Agent ([RarityS] or [RarityA]).
	Rarity Rarity `json:"rarity"`

	// Skin is the currently equipped cosmetic outfit. May be nil if default appearance is used.
	Skin *Skin `json:"skin"`

	// SplashArtURL is the absolute HTTPS URL pointing to the Agent's full splash art on the EnkaNetwork CDN.
	SplashArtURL string `json:"splash_art_url"`

	// Skills is the flat list of all individual combat abilities and passives.
	Skills []Skill `json:"skills"`

	// SkillGroups contains the Agent's skills categorized into 6 UI groups matching in-game skill tabs.
	SkillGroups []SkillGroup `json:"grouped_skills"`

	// Mindscapes contains all 6 [MindscapeNode] levels (Cinema 1–6) with their unlocked status.
	Mindscapes []MindscapeNode `json:"mindscapes"`

	// PotentialVision holds Potential Vision upgrade status and nodes (nil if Agent has no Potential Vision).
	PotentialVision *PotentialVision `json:"potential_vision"`

	// WEngine is the currently equipped W-Engine. May be nil if no weapon is equipped.
	WEngine *WEngine `json:"w_engine"`

	// DriveDiscs holds the equipped Drive Discs (slots 1–6) and active set bonuses.
	DriveDiscs DriveDiscs `json:"drive_discs"`

	// BaseStats contains the Agent's innate combat stats (Agent level growth + W-Engine Base ATK).
	BaseStats Stats `json:"base_stats"`

	// Stats contains the Agent's final calculated combat stats after applying all gear and buffs.
	Stats Stats `json:"stats"`

	// UIStats contains pre-formatted combat stat breakdowns (Base + Added = Total) with localized names and icons ready for frontend rendering.
	UIStats UIStats `json:"ui_stats"`
}

// MindscapeNode represents a single Mindscape Cinema level (M1–M6) for an [Agent].
type MindscapeNode struct {
	// Rank is the Mindscape Cinema level (1 to 6).
	Rank int `json:"rank"`

	// Name is the localized name of the Mindscape Cinema node.
	Name string `json:"name"`

	// Description is the localized description text of the Mindscape Cinema node effect.
	Description string `json:"description"`

	// FormattedHTML is the web-ready HTML description with inline CSS colors.
	FormattedHTML string `json:"formatted_html,omitempty"`

	// Unlocked indicates whether this Mindscape Cinema node is unlocked on the Agent (MindscapeCinema >= Rank).
	Unlocked bool `json:"unlocked"`
}

// FormatHTML returns the Mindscape Cinema node description formatted as HTML with inline CSS styling.
func (m MindscapeNode) FormatHTML() string {
	return formatHTML(m.Description)
}

// FormatPlainText returns the Mindscape Cinema node description stripped of Unity Rich Text formatting.
func (m MindscapeNode) FormatPlainText() string {
	return formatPlainText(m.Description)
}

// FormatMarkdown returns the Mindscape Cinema node description formatted with Markdown syntax (bold highlights).
func (m MindscapeNode) FormatMarkdown() string {
	return formatMarkdown(m.Description)
}

// PotentialVision holds Potential Vision upgrade mechanics and nodes for an [Agent].
type PotentialVision struct {
	// IsUnlocked indicates whether the Potential Vision mechanic is unlocked.
	IsUnlocked bool `json:"is_unlocked"`

	// CurrentID is the currently active upgrade node ID.
	CurrentID int `json:"current_id"`

	// Nodes contains all Potential Vision upgrade nodes for the Agent.
	Nodes []PotentialVisionNode `json:"nodes"`
}

// PotentialVisionNode represents a single Potential Vision upgrade node.
type PotentialVisionNode struct {
	// ID is the internal numeric identifier of the upgrade node.
	ID int `json:"id"`

	// Level is the level threshold of the node (1 to 6).
	Level int `json:"level"`

	// LevelName is the localized level title.
	LevelName string `json:"level_name"`

	// Title is the localized title of the upgrade effect.
	Title string `json:"title"`

	// Description is the localized effect description.
	Description string `json:"description"`

	// FormattedHTML is the web-ready HTML description with inline CSS colors.
	FormattedHTML string `json:"formatted_html,omitempty"`

	// IsActive indicates whether this upgrade node is currently active on the Agent.
	IsActive bool `json:"is_active"`
}

// FormatHTML returns the PotentialVisionNode description formatted as HTML with inline CSS styling.
func (p PotentialVisionNode) FormatHTML() string {
	return formatHTML(p.Description)
}

// FormatPlainText returns the PotentialVisionNode description stripped of Unity Rich Text formatting.
func (p PotentialVisionNode) FormatPlainText() string {
	return formatPlainText(p.Description)
}

// FormatMarkdown returns the PotentialVisionNode description formatted with Markdown syntax.
func (p PotentialVisionNode) FormatMarkdown() string {
	return formatMarkdown(p.Description)
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
