package fairy

import (
	"encoding/base64"
	"fmt"
)

// PropertyID represents a strongly-typed numeric identifier for combat attributes and stats in Zenless Zone Zero.
//
// The naming conventions categorize stats into three distinct layers:
//   - Base: The innate foundational stat before gear multipliers (e.g. [PropBaseATK], [PropBaseHP]).
//   - Percent: A percentage modifier applied to the base stat (e.g. [PropATKPercent], [PropHPPercent]).
//   - Flat: A direct numerical addition applied after percentage scaling (e.g. [PropATKFlat], [PropHPFlat]).
type PropertyID int

// Strongly-typed PropertyID constants for all combat attributes.
const (
	// PropBaseHP represents the base Health Points stat.
	PropBaseHP PropertyID = 11101
	// PropHPPercent represents a percentage increase to Health Points.
	PropHPPercent PropertyID = 11102
	// PropHPFlat represents a flat increase to Health Points.
	PropHPFlat PropertyID = 11103
	// PropHPPercentBonus represents an additional percentage bonus to Health Points.
	PropHPPercentBonus PropertyID = 11104
	// PropHPFlatBonus represents an additional flat bonus to Health Points.
	PropHPFlatBonus PropertyID = 11105

	// PropBaseATK represents the base Attack stat.
	PropBaseATK PropertyID = 12101
	// PropATKPercent represents a percentage increase to Attack.
	PropATKPercent PropertyID = 12102
	// PropATKFlat represents a flat increase to Attack.
	PropATKFlat PropertyID = 12103

	// PropBaseDEF represents the base Defense stat.
	PropBaseDEF PropertyID = 13101
	// PropDEFPercent represents a percentage increase to Defense.
	PropDEFPercent PropertyID = 13102
	// PropDEFFlat represents a flat increase to Defense.
	PropDEFFlat PropertyID = 13103

	// PropBaseImpact represents the base Impact stat (influences Daze build-up).
	PropBaseImpact PropertyID = 12201
	// PropImpactPercent represents a percentage increase to Impact.
	PropImpactPercent PropertyID = 12202
	// PropImpactFlat represents a flat increase to Impact.
	PropImpactFlat PropertyID = 12203

	// PropBaseCritRate represents the base Critical Rate stat.
	PropBaseCritRate PropertyID = 20101
	// PropCritRate represents an increase to Critical Rate.
	PropCritRate PropertyID = 20103
	// PropBaseCritDMG represents the base Critical Damage stat.
	PropBaseCritDMG PropertyID = 21101
	// PropCritDMG represents an increase to Critical Damage.
	PropCritDMG PropertyID = 21103

	// PropBasePENRatio represents the base Penetration Ratio stat.
	PropBasePENRatio PropertyID = 23101
	// PropPENRatio represents an increase to Penetration Ratio.
	PropPENRatio PropertyID = 23103
	// PropBasePENFlat represents the base flat Penetration stat.
	PropBasePENFlat PropertyID = 23201
	// PropPENFlat represents an increase to flat Penetration.
	PropPENFlat PropertyID = 23203

	// PropBaseEnergyRegen represents the base Energy Regeneration stat.
	PropBaseEnergyRegen PropertyID = 30501
	// PropEnergyRegenPercent represents a percentage increase to Energy Regeneration.
	PropEnergyRegenPercent PropertyID = 30502
	// PropEnergyRegen represents an increase to Energy Regeneration.
	PropEnergyRegen PropertyID = 30503

	// PropBaseAnomalyProficiency represents the base Anomaly Proficiency stat (scales Anomaly damage).
	PropBaseAnomalyProficiency PropertyID = 31201
	// PropAnomalyProficiencyPercent represents a percentage increase to Anomaly Proficiency.
	PropAnomalyProficiencyPercent PropertyID = 31202
	// PropAnomalyProficiency represents an increase to Anomaly Proficiency.
	PropAnomalyProficiency PropertyID = 31203

	// PropBaseAnomalyMastery represents the base Anomaly Mastery stat (scales Anomaly build-up speed).
	PropBaseAnomalyMastery PropertyID = 31401
	// PropAnomalyMasteryPercent represents a percentage increase to Anomaly Mastery.
	PropAnomalyMasteryPercent PropertyID = 31402
	// PropAnomalyMastery represents an increase to Anomaly Mastery.
	PropAnomalyMastery PropertyID = 31403

	// PropBaseSheerForce represents the base Sheer Force stat (for Rupture agents).
	PropBaseSheerForce PropertyID = 12301
	// PropSheerForce represents an increase to Sheer Force.
	PropSheerForce PropertyID = 12303

	// PropBaseRpRecover represents the base Adrenaline Auto-Accumulation stat (for Rupture agents).
	PropBaseRpRecover PropertyID = 32001
	// PropRpRecoverPercent represents a percentage increase to Adrenaline Auto-Accumulation.
	PropRpRecoverPercent PropertyID = 32002
	// PropRpRecover represents an increase to Adrenaline Auto-Accumulation.
	PropRpRecover PropertyID = 32003

	// PropPhysicalDMGBonus represents the Physical Damage Bonus stat.
	PropPhysicalDMGBonus PropertyID = 31505
	// PropFireDMGBonus represents the Fire Damage Bonus stat.
	PropFireDMGBonus PropertyID = 31605
	// PropIceDMGBonus represents the Ice Damage Bonus stat.
	PropIceDMGBonus PropertyID = 31705
	// PropElectricDMGBonus represents the Electric Damage Bonus stat.
	PropElectricDMGBonus PropertyID = 31805
	// PropEtherDMGBonus represents the Ether Damage Bonus stat.
	PropEtherDMGBonus PropertyID = 31905
	// PropWindDMGBonus represents the Wind Damage Bonus stat.
	PropWindDMGBonus PropertyID = 32305
)

// Property group base prefixes for damage bonuses.
const (
	propGroupGeneralDMG  PropertyID = 31000
	propGroupPhysicalDMG PropertyID = 31500
	propGroupFireDMG     PropertyID = 31600
	propGroupIceDMG      PropertyID = 31700
	propGroupElectricDMG PropertyID = 31800
	propGroupEtherDMG    PropertyID = 31900
	propGroupWindDMG     PropertyID = 32300
)

// Datamined property localization keys in locs.json.
const (
	locKeyHP                 = "HpMax"
	locKeyATK                = "Atk"
	locKeyDEF                = "Def"
	locKeyImpact             = "BreakStun"
	locKeyCritRate           = "Crit"
	locKeyCritDMG            = "CritDmg"
	locKeyAnomalyMastery     = "ElementAbnormalPower"
	locKeyAnomalyProficiency = "ElementMystery"
	locKeyPenRatio           = "PenRatio"
	locKeyPenFlat            = "PenDelta"
	locKeyEnergyRegen        = "SpRecover"
	locKeyRpRecover          = "RpRecover"
	locKeySheerForce         = "SkipDefAtk"

	locKeyGeneralDMGBonus  = "AddedDamageRatio"
	locKeyPhysicalDMGBonus = "AddedDamageRatio_Physics"
	locKeyFireDMGBonus     = "AddedDamageRatio_Fire"
	locKeyIceDMGBonus      = "AddedDamageRatio_Ice"
	locKeyElectricDMGBonus = "AddedDamageRatio_Elec"
	locKeyEtherDMGBonus    = "AddedDamageRatio_Ether"
	locKeyWindDMGBonus     = "AddedDamageRatio_Wind"
)

// StatValue represents a single combat stat or sub-stat entry (such as a Drive Disc main stat or substat roll).
type StatValue struct {
	// PropertyID is the internal strongly-typed [PropertyID] (e.g. [PropATKPercent]).
	PropertyID PropertyID `json:"property_id"`

	// Name is the localized display name of the stat (e.g. "ATK", "CRIT Rate").
	Name string `json:"name"`

	// Value is the calculated numerical value of the stat. Percentage stats are stored in decimal format (e.g. 0.048 for 4.8%).
	Value float64 `json:"value"`

	// IsPercent indicates whether the stat represents a percentage value.
	IsPercent bool `json:"is_percent"`

	// Rolls is the number of times this stat was upgraded (1 for base roll, up to 5 with upgrades).
	Rolls int `json:"rolls"`

	// IconURL is the base64 Data URI string ("data:image/svg+xml;base64,...") containing the stat's SVG icon.
	IconURL string `json:"icon_url"`
}

// SVG returns the raw inline SVG markup string for the stat property.
func (p PropertyID) SVG() string {
	if svg, ok := propertySVGMap[p]; ok {
		return svg
	}
	return ""
}

// IconURL returns a base64-encoded Data URI string ("data:image/svg+xml;base64,...")
// containing the stat property's SVG icon for direct use in frontend <img> tags.
func (p PropertyID) IconURL() string {
	svg := p.SVG()
	if svg == "" {
		return ""
	}
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg))
}

// SVG returns the raw inline SVG markup string for the stat value.
func (s StatValue) SVG() string {
	return s.PropertyID.SVG()
}

// DisplayValue returns the stat's value formatted as a human-readable string (e.g. "4.8%" or "310").
func (s StatValue) DisplayValue() string {
	if s.IsPercent {
		return fmt.Sprintf("%.1f%%", s.Value*100)
	}
	return fmt.Sprintf("%.0f", s.Value)
}

// Stats represents the complete aggregated numerical combat stats of an [Agent].
//
// Value representations:
//   - Percentage stats ([Stats.CritRate], [Stats.CritDMG], [Stats.AttributeDMGBonus], [Stats.PenRatio]) are stored as decimals (e.g. 0.05 = 5%, 1.50 = 150%).
//   - [Stats.EnergyRegen] is stored as energy recovered per second (e.g. 1.20).
//   - Flat stats ([Stats.HP], [Stats.ATK], [Stats.DEF], [Stats.Impact], [Stats.AnomalyMastery], [Stats.AnomalyProficiency], [Stats.PenFlat], [Stats.SheerForce]) are stored as raw floating-point numbers.
//
// Use [Stats.Formatted] or [Agent.UIStats] to retrieve pre-formatted string representations.
type Stats struct {
	// HP is the total Health Points.
	HP float64 `json:"hp"`

	// ATK is the total Attack.
	ATK float64 `json:"atk"`

	// DEF is the total Defense.
	DEF float64 `json:"def"`

	// Impact is the Impact stat (influences Daze accumulation rate).
	Impact float64 `json:"impact"`

	// CritRate is the Critical Hit Rate as a decimal fraction (e.g. 0.05 for 5%).
	CritRate float64 `json:"crit_rate"`

	// CritDMG is the Critical Hit Damage as a decimal fraction (e.g. 1.50 for 150%).
	CritDMG float64 `json:"crit_dmg"`

	// AttributeDMGBonus is the matching elemental Damage Bonus as a decimal fraction (e.g. 0.30 for 30%).
	AttributeDMGBonus float64 `json:"attribute_dmg_bonus"`

	// AnomalyMastery is the Anomaly Mastery stat (influences Anomaly Buildup speed).
	AnomalyMastery float64 `json:"anomaly_mastery"`

	// AnomalyProficiency is the Anomaly Proficiency stat (scales Anomaly damage).
	AnomalyProficiency float64 `json:"anomaly_proficiency"`

	// PenRatio is the Penetration Ratio as a decimal fraction (ignores a percentage of enemy DEF).
	PenRatio float64 `json:"pen_ratio"`

	// PenFlat is the flat Penetration stat (ignores a flat amount of enemy DEF).
	PenFlat float64 `json:"pen_flat"`

	// EnergyRegen is the Energy Regeneration rate per second (e.g. 1.20).
	EnergyRegen float64 `json:"energy_regen"`

	// SheerForce is the Sheer Force stat (damage multiplier for Rupture agents, ignoring DEF).
	SheerForce float64 `json:"sheer_force"`
}

// FormattedStats contains the agent's combat stats pre-formatted as human-readable strings.
// This is convenient for frontend rendering where formatted values (e.g. "50.0%", "3,120") are required directly.
type FormattedStats struct {
	// HP is the formatted Health Points string.
	HP string `json:"hp"`

	// ATK is the formatted Attack string.
	ATK string `json:"atk"`

	// DEF is the formatted Defense string.
	DEF string `json:"def"`

	// Impact is the formatted Impact string.
	Impact string `json:"impact"`

	// CritRate is the formatted Critical Rate percentage string (e.g. "50.0%").
	CritRate string `json:"crit_rate"`

	// CritDMG is the formatted Critical Damage percentage string (e.g. "150.0%").
	CritDMG string `json:"crit_dmg"`

	// AttributeDMGBonus is the formatted Attribute Damage Bonus percentage string (e.g. "30.0%").
	AttributeDMGBonus string `json:"attribute_dmg_bonus"`

	// AnomalyMastery is the formatted Anomaly Mastery string.
	AnomalyMastery string `json:"anomaly_mastery"`

	// AnomalyProficiency is the formatted Anomaly Proficiency string.
	AnomalyProficiency string `json:"anomaly_proficiency"`

	// PenRatio is the formatted Penetration Ratio percentage string (e.g. "24.0%").
	PenRatio string `json:"pen_ratio"`

	// PenFlat is the formatted Flat Penetration string.
	PenFlat string `json:"pen_flat"`

	// EnergyRegen is the formatted Energy Regeneration string (e.g. "1.20").
	EnergyRegen string `json:"energy_regen"`

	// SheerForce is the formatted Sheer Force string.
	SheerForce string `json:"sheer_force"`
}

// FormattedStatBreakdown represents a single combat stat broken down into its base and added components,
// pre-formatted as human-readable strings ready for UI display.
type FormattedStatBreakdown struct {
	// PropertyID is the strongly-typed [PropertyID] of the stat.
	PropertyID PropertyID `json:"property_id"`

	// Name is the localized display name of the stat.
	Name string `json:"name"`

	// Base is the pre-formatted innate base stat value.
	Base string `json:"base"`

	// Added is the pre-formatted added stat value from W-Engines, Drive Discs, and set bonuses.
	Added string `json:"added"`

	// Total is the pre-formatted final combat stat value.
	Total string `json:"total"`

	// IconURL is the base64 Data URI string containing the stat's SVG icon.
	IconURL string `json:"icon_url"`
}

// UIStats contains all combat stats broken down into Base + Added = Total components,
// with localized names and icon URLs, structured for frontend profile and agent inspect panels.
type UIStats struct {
	// HP is the Health Points breakdown.
	HP FormattedStatBreakdown `json:"hp"`

	// ATK is the Attack breakdown.
	ATK FormattedStatBreakdown `json:"atk"`

	// DEF is the Defense breakdown.
	DEF FormattedStatBreakdown `json:"def"`

	// Impact is the Impact breakdown.
	Impact FormattedStatBreakdown `json:"impact"`

	// CritRate is the Critical Rate breakdown.
	CritRate FormattedStatBreakdown `json:"crit_rate"`

	// CritDMG is the Critical Damage breakdown.
	CritDMG FormattedStatBreakdown `json:"crit_dmg"`

	// AttributeDMGBonus is the matching elemental Damage Bonus breakdown.
	AttributeDMGBonus FormattedStatBreakdown `json:"attribute_dmg_bonus"`

	// AnomalyMastery is the Anomaly Mastery breakdown.
	AnomalyMastery FormattedStatBreakdown `json:"anomaly_mastery"`

	// AnomalyProficiency is the Anomaly Proficiency breakdown.
	AnomalyProficiency FormattedStatBreakdown `json:"anomaly_proficiency"`

	// PenRatio is the Penetration Ratio breakdown.
	PenRatio FormattedStatBreakdown `json:"pen_ratio"`

	// PenFlat is the Flat Penetration breakdown.
	PenFlat FormattedStatBreakdown `json:"pen_flat"`

	// EnergyRegen is the Energy Regeneration breakdown.
	EnergyRegen FormattedStatBreakdown `json:"energy_regen"`

	// SheerForce is the Sheer Force breakdown (for Rupture agents).
	SheerForce FormattedStatBreakdown `json:"sheer_force"`
}

// List returns all combat stat breakdowns as a slice in the canonical in-game display order.
func (u UIStats) List() []FormattedStatBreakdown {
	list := make([]FormattedStatBreakdown, 0, 13)
	list = append(list,
		u.HP,
		u.ATK,
		u.DEF,
		u.Impact,
		u.CritRate,
		u.CritDMG,
		u.AnomalyMastery,
		u.AnomalyProficiency,
		u.PenRatio,
		u.PenFlat,
		u.EnergyRegen,
	)
	// Skip Lumiflux
	if u.AttributeDMGBonus.PropertyID != 0 && u.AttributeDMGBonus.Name != "" {
		list = append(list, u.AttributeDMGBonus)
	}
	list = append(list,
		u.SheerForce,
	)
	return list
}

// Formatted returns a new [FormattedStats] struct where all numerical stats
// are converted into precise, human-readable strings (e.g. "50.0%" instead of 0.5, "3120" instead of 3120.0).
func (s *Stats) Formatted() FormattedStats {
	return FormattedStats{
		HP:                 fmt.Sprintf("%.0f", s.HP),
		ATK:                fmt.Sprintf("%.0f", s.ATK),
		DEF:                fmt.Sprintf("%.0f", s.DEF),
		Impact:             fmt.Sprintf("%.0f", s.Impact),
		CritRate:           fmt.Sprintf("%.1f%%", s.CritRate*100),
		CritDMG:            fmt.Sprintf("%.1f%%", s.CritDMG*100),
		AttributeDMGBonus:  fmt.Sprintf("%.1f%%", s.AttributeDMGBonus*100),
		AnomalyMastery:     fmt.Sprintf("%.0f", s.AnomalyMastery),
		AnomalyProficiency: fmt.Sprintf("%.0f", s.AnomalyProficiency),
		PenRatio:           fmt.Sprintf("%.1f%%", s.PenRatio*100),
		PenFlat:            fmt.Sprintf("%.0f", s.PenFlat),
		EnergyRegen:        fmt.Sprintf("%.2f", s.EnergyRegen),
		SheerForce:         fmt.Sprintf("%.0f", s.SheerForce),
	}
}

// List returns all numeric combat stats as a slice of [StatValue] in canonical in-game display order.
func (s Stats) List() []StatValue {
	return []StatValue{
		{PropertyID: PropBaseHP, Value: s.HP, IsPercent: false, IconURL: PropBaseHP.IconURL()},
		{PropertyID: PropBaseATK, Value: s.ATK, IsPercent: false, IconURL: PropBaseATK.IconURL()},
		{PropertyID: PropBaseDEF, Value: s.DEF, IsPercent: false, IconURL: PropBaseDEF.IconURL()},
		{PropertyID: PropBaseImpact, Value: s.Impact, IsPercent: false, IconURL: PropBaseImpact.IconURL()},
		{PropertyID: PropBaseCritRate, Value: s.CritRate, IsPercent: true, IconURL: PropBaseCritRate.IconURL()},
		{PropertyID: PropBaseCritDMG, Value: s.CritDMG, IsPercent: true, IconURL: PropBaseCritDMG.IconURL()},
		{PropertyID: PropBaseAnomalyMastery, Value: s.AnomalyMastery, IsPercent: false, IconURL: PropBaseAnomalyMastery.IconURL()},
		{PropertyID: PropBaseAnomalyProficiency, Value: s.AnomalyProficiency, IsPercent: false, IconURL: PropBaseAnomalyProficiency.IconURL()},
		{PropertyID: PropBasePENRatio, Value: s.PenRatio, IsPercent: true, IconURL: PropBasePENRatio.IconURL()},
		{PropertyID: PropBasePENFlat, Value: s.PenFlat, IsPercent: false, IconURL: PropBasePENFlat.IconURL()},
		{PropertyID: PropBaseEnergyRegen, Value: s.EnergyRegen, IsPercent: false, IconURL: PropBaseEnergyRegen.IconURL()},
		{PropertyID: propGroupGeneralDMG, Value: s.AttributeDMGBonus, IsPercent: true, IconURL: propGroupGeneralDMG.IconURL()},
		{PropertyID: PropBaseSheerForce, Value: s.SheerForce, IsPercent: false, IconURL: PropBaseSheerForce.IconURL()},
	}
}

// formatFlatBreakdown is a helper to format a flat stat breakdown.
func formatFlatBreakdown(propID PropertyID, name string, base, total float64, precision int) FormattedStatBreakdown {
	added := total - base
	return FormattedStatBreakdown{
		PropertyID: propID,
		Name:       name,
		Base:       fmt.Sprintf("%.*f", precision, base),
		Added:      fmt.Sprintf("%.*f", precision, added),
		Total:      fmt.Sprintf("%.*f", precision, total),
		IconURL:    propID.IconURL(),
	}
}

// formatPercentBreakdown is a helper to format a percentage stat breakdown.
func formatPercentBreakdown(propID PropertyID, name string, base, total float64, precision int) FormattedStatBreakdown {
	added := total - base
	return FormattedStatBreakdown{
		PropertyID: propID,
		Name:       name,
		Base:       fmt.Sprintf("%.*f%%", precision, base*100),
		Added:      fmt.Sprintf("%.*f%%", precision, added*100),
		Total:      fmt.Sprintf("%.*f%%", precision, total*100),
		IconURL:    propID.IconURL(),
	}
}
