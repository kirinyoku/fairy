package fairy

import (
	"encoding/base64"
	"fmt"
)

// PropertyID represents a strongly-typed ID for combat properties.
// The naming convention follows:
// - Base: The character's innate foundational stat.
// - Percent / PercentBonus: A percentage modifier applied to the base stat.
// - Flat / FlatBonus: A direct numerical addition applied after percentages.
type PropertyID int

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

	// PropBaseImpact represents the base Impact stat.
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

	// PropBaseAnomalyProficiency represents the base Anomaly Proficiency stat.
	PropBaseAnomalyProficiency PropertyID = 31201
	// PropAnomalyProficiencyPercent represents a percentage increase to Anomaly Proficiency.
	PropAnomalyProficiencyPercent PropertyID = 31202
	// PropAnomalyProficiency represents an increase to Anomaly Proficiency.
	PropAnomalyProficiency PropertyID = 31203

	// PropBaseAnomalyMastery represents the base Anomaly Mastery stat.
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

	// PropertyIDs for elemental damage bonuses.
	PropPhysicalDMGBonus PropertyID = 31505
	PropFireDMGBonus     PropertyID = 31605
	PropIceDMGBonus      PropertyID = 31705
	PropElectricDMGBonus PropertyID = 31805
	PropEtherDMGBonus    PropertyID = 31905
	PropWindDMGBonus     PropertyID = 32305
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

// StatValue represents a single combat stat (main or sub stat).
type StatValue struct {
	PropertyID PropertyID `json:"property_id"` // The internal property ID (e.g., 12102 for ATK%).
	Name       string     `json:"name"`        // The localized name of the stat.
	Value      float64    `json:"value"`       // The final calculated value of the stat.
	IsPercent  bool       `json:"is_percent"`  // Indicates if the stat is a percentage.
	Rolls      int        `json:"rolls"`       // The number of times this stat was upgraded (1 for base, up to 5 for max upgrades).
	IconURL    string     `json:"icon_url"`    // The base64 Data URI string containing the stat's SVG icon.
}

// SVG returns the raw inline SVG markup string for the stat property.
func (p PropertyID) SVG() string {
	if svg, ok := propertySVGMap[p]; ok {
		return svg
	}
	return ""
}

// IconURL returns the base64-encoded Data URI string containing the stat property's SVG icon.
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

// DisplayValue returns the stat's value formatted as a human-readable string.
// Percentages are multiplied by 100 and formatted with a '%' sign.
func (s StatValue) DisplayValue() string {
	if s.IsPercent {
		return fmt.Sprintf("%.1f%%", s.Value*100)
	}
	return fmt.Sprintf("%.0f", s.Value)
}

// Stats represents the aggregated combat stats of an agent.
// This structure keeps fields minimal and flat for easy access.
// Internal representation of percentages is in decimal form (e.g., CritRate 0.05 = 5%).
// EnergyRegen is also represented as its divided final value (e.g., 1.20).
// Precise calculations and final formulas are opt-in via the calc package.
type Stats struct {
	HP                 float64 `json:"hp"`                  // Total Health Points.
	ATK                float64 `json:"atk"`                 // Total Attack.
	DEF                float64 `json:"def"`                 // Total Defense.
	Impact             float64 `json:"impact"`              // Impact (influences Daze build-up).
	CritRate           float64 `json:"crit_rate"`           // Critical Hit Rate (as a decimal, e.g., 0.05 for 5%).
	CritDMG            float64 `json:"crit_dmg"`            // Critical Hit Damage (as a decimal, e.g., 1.50 for 150%).
	AttributeDMGBonus  float64 `json:"attribute_dmg_bonus"` // Attribute DMG Bonus (as a decimal, e.g., 0.30 for 30%).
	AnomalyMastery     float64 `json:"anomaly_mastery"`     // Anomaly Mastery (influences Anomaly Buildup rate).
	AnomalyProficiency float64 `json:"anomaly_proficiency"` // Anomaly Proficiency (influences Anomaly Damage).
	PenRatio           float64 `json:"pen_ratio"`           // Penetration Ratio (ignores a percentage of enemy DEF).
	PenFlat            float64 `json:"pen_flat"`            // Flat Penetration (ignores a flat amount of enemy DEF).
	EnergyRegen        float64 `json:"energy_regen"`        // Energy Regeneration rate (as a decimal, e.g., 1.20).
	SheerForce         float64 `json:"sheer_force"`         // Sheer Force (damage multiplier for Rupture agents, ignoring DEF).
}

// FormattedStats contains the agent's combat stats pre-formatted as human-readable strings.
// This is extremely useful for UI/Frontend developers who just want to display the values.
type FormattedStats struct {
	HP                 string `json:"hp"`
	ATK                string `json:"atk"`
	DEF                string `json:"def"`
	Impact             string `json:"impact"`
	CritRate           string `json:"crit_rate"`
	CritDMG            string `json:"crit_dmg"`
	AttributeDMGBonus  string `json:"attribute_dmg_bonus"`
	AnomalyMastery     string `json:"anomaly_mastery"`
	AnomalyProficiency string `json:"anomaly_proficiency"`
	PenRatio           string `json:"pen_ratio"`
	PenFlat            string `json:"pen_flat"`
	EnergyRegen        string `json:"energy_regen"`
	SheerForce         string `json:"sheer_force"`
}

// FormattedStatBreakdown represents a single stat broken down into its base and added components,
// pre-formatted as human-readable strings for UI display.
type FormattedStatBreakdown struct {
	PropertyID PropertyID `json:"property_id"` // Strongly-typed PropertyID of the stat.
	Name       string     `json:"name"`        // Localized stat display name.
	Base       string     `json:"base"`        // Base stat value (pre-formatted string).
	Added      string     `json:"added"`       // Added stat value from gear/buffs (pre-formatted string).
	Total      string     `json:"total"`       // Final combat stat value (pre-formatted string).
	IconURL    string     `json:"icon_url"`    // The base64 Data URI string containing the stat's SVG icon.
}

// UIStats contains all combat stats broken down into base and added components,
// ready to be displayed on a frontend.
type UIStats struct {
	HP                 FormattedStatBreakdown `json:"hp"`
	ATK                FormattedStatBreakdown `json:"atk"`
	DEF                FormattedStatBreakdown `json:"def"`
	Impact             FormattedStatBreakdown `json:"impact"`
	CritRate           FormattedStatBreakdown `json:"crit_rate"`
	CritDMG            FormattedStatBreakdown `json:"crit_dmg"`
	AttributeDMGBonus  FormattedStatBreakdown `json:"attribute_dmg_bonus"`
	AnomalyMastery     FormattedStatBreakdown `json:"anomaly_mastery"`
	AnomalyProficiency FormattedStatBreakdown `json:"anomaly_proficiency"`
	PenRatio           FormattedStatBreakdown `json:"pen_ratio"`
	PenFlat            FormattedStatBreakdown `json:"pen_flat"`
	EnergyRegen        FormattedStatBreakdown `json:"energy_regen"`
	SheerForce         FormattedStatBreakdown `json:"sheer_force"`
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

// Formatted returns a new FormattedStats struct where all numerical stats
// are converted into precise, human-readable strings (e.g. "50.0%" instead of 0.5).
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

// List returns all numeric combat stats as a slice of StatValue in canonical in-game display order.
func (s Stats) List() []StatValue {
	return []StatValue{
		{PropertyID: PropBaseHP, Value: s.HP, IsPercent: false, IconURL: PropBaseHP.IconURL()},
		{PropertyID: PropBaseATK, Value: s.ATK, IsPercent: false, IconURL: PropBaseATK.IconURL()},
		{PropertyID: PropBaseDEF, Value: s.DEF, IsPercent: false, IconURL: PropBaseDEF.IconURL()},
		{PropertyID: PropBaseImpact, Value: s.Impact, IsPercent: false, IconURL: PropBaseImpact.IconURL()},
		{PropertyID: PropBaseCritRate, Value: s.CritRate, IsPercent: true, IconURL: PropBaseCritRate.IconURL()},
		{PropertyID: PropBaseCritDMG, Value: s.CritDMG, IsPercent: true, IconURL: PropBaseCritDMG.IconURL()},
		{PropertyID: propGroupGeneralDMG, Value: s.AttributeDMGBonus, IsPercent: true, IconURL: propGroupGeneralDMG.IconURL()},
		{PropertyID: PropBaseAnomalyMastery, Value: s.AnomalyMastery, IsPercent: false, IconURL: PropBaseAnomalyMastery.IconURL()},
		{PropertyID: PropBaseAnomalyProficiency, Value: s.AnomalyProficiency, IsPercent: false, IconURL: PropBaseAnomalyProficiency.IconURL()},
		{PropertyID: PropBasePENRatio, Value: s.PenRatio, IsPercent: true, IconURL: PropBasePENRatio.IconURL()},
		{PropertyID: PropBasePENFlat, Value: s.PenFlat, IsPercent: false, IconURL: PropBasePENFlat.IconURL()},
		{PropertyID: PropBaseEnergyRegen, Value: s.EnergyRegen, IsPercent: false, IconURL: PropBaseEnergyRegen.IconURL()},
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
