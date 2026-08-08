package fairy

import "fmt"

// PropertyID represents a strongly-typed ID for combat properties.
// The naming convention follows:
// - Base: The character's innate foundational stat (e.g., PropBaseHP).
// - Percent / PercentBonus: A percentage modifier applied to the base stat.
// - Flat / FlatBonus: A direct numerical addition applied after percentages.
type PropertyID int

const (
	// Health properties
	PropBaseHP         PropertyID = 11101
	PropHPPercent      PropertyID = 11102
	PropHPFlat         PropertyID = 11103
	PropHPPercentBonus PropertyID = 11104
	PropHPFlatBonus    PropertyID = 11105

	// Attack properties
	PropBaseATK    PropertyID = 12101
	PropATKPercent PropertyID = 12102
	PropATKFlat    PropertyID = 12103

	// Defense properties
	PropBaseDEF    PropertyID = 13101
	PropDEFPercent PropertyID = 13102
	PropDEFFlat    PropertyID = 13103

	// Impact properties
	PropBaseImpact    PropertyID = 12201
	PropImpactPercent PropertyID = 12202
	PropImpactFlat    PropertyID = 12203

	// Critical properties
	PropBaseCritRate PropertyID = 20101
	PropCritRate     PropertyID = 20103
	PropBaseCritDMG  PropertyID = 21101
	PropCritDMG      PropertyID = 21103

	// Penetration properties
	PropBasePENRatio PropertyID = 23101
	PropPENRatio     PropertyID = 23103
	PropBasePENFlat  PropertyID = 23201
	PropPENFlat      PropertyID = 23203

	// Energy properties
	PropBaseEnergyRegen    PropertyID = 30501
	PropEnergyRegenPercent PropertyID = 30502
	PropEnergyRegen        PropertyID = 30503

	// Anomaly properties
	PropBaseAnomalyMastery        PropertyID = 31201
	PropAnomalyMastery            PropertyID = 31203
	PropBaseAnomalyProficiency    PropertyID = 31401
	PropAnomalyProficiencyPercent PropertyID = 31402
	PropAnomalyProficiency        PropertyID = 31403
)

// StatValue represents a single combat stat (main or sub stat).
type StatValue struct {
	PropertyID PropertyID `json:"property_id"` // The internal property ID (e.g., 12102 for ATK%).
	Name       string     `json:"name"`        // The localized name of the stat.
	Value      float64    `json:"value"`       // The final calculated value of the stat.
	IsPercent  bool       `json:"is_percent"`  // Indicates if the stat is a percentage.
	Rolls      int        `json:"rolls"`       // The number of times this stat was upgraded (1 for base, up to 5 for max upgrades).
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
// NOTE: Internal representation of percentages is in decimal form (e.g., CritRate 0.05 = 5%).
// EnergyRegen is also represented as its divided final value (e.g., 1.20).
// Precise calculations and final formulas are opt-in via the calc package.
type Stats struct {
	HP                 float64 `json:"hp"`                  // Total Health Points.
	ATK                float64 `json:"atk"`                 // Total Attack.
	DEF                float64 `json:"def"`                 // Total Defense.
	Impact             float64 `json:"impact"`              // Impact (influences Daze build-up).
	CritRate           float64 `json:"crit_rate"`           // Critical Hit Rate (as a decimal, e.g., 0.05 for 5%).
	CritDMG            float64 `json:"crit_dmg"`            // Critical Hit Damage (as a decimal, e.g., 1.50 for 150%).
	AnomalyMastery     float64 `json:"anomaly_mastery"`     // Anomaly Mastery (influences Anomaly Buildup rate).
	AnomalyProficiency float64 `json:"anomaly_proficiency"` // Anomaly Proficiency (influences Anomaly Damage).
	PenRatio           float64 `json:"pen_ratio"`           // Penetration Ratio (ignores a percentage of enemy DEF).
	PenFlat            float64 `json:"pen_flat"`            // Flat Penetration (ignores a flat amount of enemy DEF).
	EnergyRegen        float64 `json:"energy_regen"`        // Energy Regeneration rate (as a decimal, e.g., 1.20).
}

// FormattedStats contains the agent's combat stats pre-formatted as human-readable strings.
// This is extremely useful for UI/Frontend developers who just want to display the values.
// Note: JSON tags intentionally match those in UIStats for frontend consistency.
type FormattedStats struct {
	HP                 string `json:"hp"`
	ATK                string `json:"atk"`
	DEF                string `json:"def"`
	Impact             string `json:"impact"`
	CritRate           string `json:"crit_rate"`
	CritDMG            string `json:"crit_dmg"`
	AnomalyMastery     string `json:"anomaly_mastery"`
	AnomalyProficiency string `json:"anomaly_proficiency"`
	PenRatio           string `json:"pen_ratio"`
	PenFlat            string `json:"pen_flat"`
	EnergyRegen        string `json:"energy_regen"`
}

// FormattedStatBreakdown represents a single stat broken down into its base and added components,
// pre-formatted as human-readable strings for UI display.
type FormattedStatBreakdown struct {
	Base  string `json:"base"`
	Added string `json:"added"`
	Total string `json:"total"`
}

// UIStats contains all combat stats broken down into base and added components,
// ready to be displayed on a frontend profile page (like Enka.Network).
type UIStats struct {
	HP                 FormattedStatBreakdown `json:"hp"`
	ATK                FormattedStatBreakdown `json:"atk"`
	DEF                FormattedStatBreakdown `json:"def"`
	Impact             FormattedStatBreakdown `json:"impact"`
	CritRate           FormattedStatBreakdown `json:"crit_rate"`
	CritDMG            FormattedStatBreakdown `json:"crit_dmg"`
	AnomalyMastery     FormattedStatBreakdown `json:"anomaly_mastery"`
	AnomalyProficiency FormattedStatBreakdown `json:"anomaly_proficiency"`
	PenRatio           FormattedStatBreakdown `json:"pen_ratio"`
	PenFlat            FormattedStatBreakdown `json:"pen_flat"`
	EnergyRegen        FormattedStatBreakdown `json:"energy_regen"`
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
		AnomalyMastery:     fmt.Sprintf("%.0f", s.AnomalyMastery),
		AnomalyProficiency: fmt.Sprintf("%.0f", s.AnomalyProficiency),
		PenRatio:           fmt.Sprintf("%.1f%%", s.PenRatio*100),
		PenFlat:            fmt.Sprintf("%.0f", s.PenFlat),
		EnergyRegen:        fmt.Sprintf("%.2f", s.EnergyRegen), // usually something like 1.20
	}
}

// formatBreakdown is an internal helper to format a single stat breakdown.
// precisionFlat and precisionPercent determine the number of decimal places used.
func formatBreakdown(base, total float64, isPercent bool, precisionFlat int, precisionPercent int) FormattedStatBreakdown {
	added := total - base
	if isPercent {
		format := fmt.Sprintf("%%.%df%%%%", precisionPercent)
		return FormattedStatBreakdown{
			Base:  fmt.Sprintf(format, base*100),
			Added: fmt.Sprintf(format, added*100),
			Total: fmt.Sprintf(format, total*100),
		}
	}
	format := fmt.Sprintf("%%.%df", precisionFlat)
	return FormattedStatBreakdown{
		Base:  fmt.Sprintf(format, base),
		Added: fmt.Sprintf(format, added),
		Total: fmt.Sprintf(format, total),
	}
}
