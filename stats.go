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
	locKeySheerForce         = "SkipDefAtk"

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
	Name       string     `json:"name"`        // Localized stat display name (e.g., "HP", "ATK" / "Сила атаки").
	Base       string     `json:"base"`        // Base stat value (pre-formatted string).
	Added      string     `json:"added"`       // Added stat value from gear/buffs (pre-formatted string).
	Total      string     `json:"total"`       // Final combat stat value (pre-formatted string).
	IconURL    string     `json:"icon_url"`    // The base64 Data URI string containing the stat's SVG icon.
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
	AttributeDMGBonus  FormattedStatBreakdown `json:"attribute_dmg_bonus"`
	AnomalyMastery     FormattedStatBreakdown `json:"anomaly_mastery"`
	AnomalyProficiency FormattedStatBreakdown `json:"anomaly_proficiency"`
	PenRatio           FormattedStatBreakdown `json:"pen_ratio"`
	PenFlat            FormattedStatBreakdown `json:"pen_flat"`
	EnergyRegen        FormattedStatBreakdown `json:"energy_regen"`
	SheerForce         FormattedStatBreakdown `json:"sheer_force"`
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
		EnergyRegen:        fmt.Sprintf("%.2f", s.EnergyRegen), // usually something like 1.20
		SheerForce:         fmt.Sprintf("%.0f", s.SheerForce),
	}
}

// formatBreakdown is an internal helper to format a single stat breakdown.
// precisionFlat and precisionPercent determine the number of decimal places used.
func formatBreakdown(propID PropertyID, name string, base, total float64, isPercent bool, precisionFlat int, precisionPercent int) FormattedStatBreakdown {
	added := total - base
	var baseStr, addedStr, totalStr string
	if isPercent {
		format := fmt.Sprintf("%%.%df%%%%", precisionPercent)
		baseStr = fmt.Sprintf(format, base*100)
		addedStr = fmt.Sprintf(format, added*100)
		totalStr = fmt.Sprintf(format, total*100)
	} else {
		format := fmt.Sprintf("%%.%df", precisionFlat)
		baseStr = fmt.Sprintf(format, base)
		addedStr = fmt.Sprintf(format, added)
		totalStr = fmt.Sprintf(format, total)
	}

	return FormattedStatBreakdown{
		PropertyID: propID,
		Name:       name,
		Base:       baseStr,
		Added:      addedStr,
		Total:      totalStr,
		IconURL:    propID.IconURL(),
	}
}

const (
	svgHP                 = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon HpMax" viewBox="0 0 14 14"><path d="M69.609 143.14c-.99 1.034-5.872 6.293-5.86 8.987.023 4.992 4.376 7.404 8.375 5.914.7-.261 1.294-.675 1.825-1.198 3.243-3.198 1.972-7.063-.222-9.242l-3.806-4.452a.207.207 0 0 0-.312-.01zm3.304 6.883c-2.947-.67-2.915 2.468-5.26 2.04-.76-.138-1.242-.545-1.445-1.126.128-.588 2.736-3.688 3.407-4.49a.217.217 0 0 1 .328-.016c.76.774 2.854 3.093 2.97 3.592" style="fill:#ffffff;stroke:none" transform="translate(-54.98 -126.82)scale(.88755)"/></svg>`
	svgATK                = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon Atk" viewBox="0 0 14 14"><path d="M5.096.49v5.167a1.08 1.08 0 0 0 1.08 1.08h1.667a1.08 1.08 0 0 0 1.08-1.08V.491ZM1.4 1.182a.74.74 0 0 0-.742.742v3.362A1.45 1.45 0 0 0 2.11 6.737h2.07a.305.305 0 0 0 .305-.305V1.436a.255.255 0 0 0-.255-.255Zm8.475 0a.34.34 0 0 0-.34.341v4.875a.34.34 0 0 0 .34.34h2.034a1.45 1.45 0 0 0 1.452-1.452V1.923a.74.74 0 0 0-.742-.742Zm-7.76 6.345A1.476 1.476 0 0 0 .64 9.002v3.03a1.476 1.476 0 0 0 1.476 1.477h9.77a1.476 1.476 0 0 0 1.476-1.476v-3.03a1.476 1.476 0 0 0-1.476-1.477H6.642a.81.81 0 0 0-.81.809v.782a.743.743 0 0 0 .744.743h2.258a.483.483 0 0 1 .484.483v.339a.394.394 0 0 1-.395.395H5.4a.717.717 0 0 1-.717-.717V7.864a.34.34 0 0 0-.338-.338Z" style="fill:#ffffff;stroke-width:0.909407;stroke-linecap:square;stroke-dashoffset:3.77952;stop-color:#ffffff"/></svg>`
	svgDEF                = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon Def" viewBox="0 0 14 14"><path d="M6.963.207S5.672.655 5.086.824a22 22 0 0 1-3.167.677c-.38.052-1.712.174-1.712.174s.026.848.03 1.106c.014.782.13 1.614.277 2.383.549 2.848 1.644 5.718 4.077 7.475.67.484 1.697 1.234 2.573 1.147.822-.08 1.729-.712 2.374-1.188 2.346-1.731 3.473-4.556 3.997-7.335.154-.815.258-1.653.258-2.482 0-.265-.041-1.133-.041-1.133s-1.34-.1-1.74-.166A30 30 0 0 1 8.747.753C8.204.6 6.963.207 6.963.207m-4.45 3.268c1.517-.149 3.49-.565 4.453-.951 1.024.456 4.566.94 4.566.94.037.73-.075 1.73-.23 2.446-.417 1.923-1.216 4.05-3.3 5.232-.33.187-1.032.423-1.032.423s-.905-.15-2.156-1.296C3.02 8.516 2.513 5.893 2.513 3.476m4.453.099c-1.117.38-2.248.675-3.427.717.152 2.067.374 4.056 2.24 5.565.358.29.64.483 1.214.565.595-.074 1.1-.4 1.455-.69 1.082-.885 1.701-2.306 1.882-3.674.054-.412.135-1.735.135-1.735-1.137-.02-2.355-.398-3.5-.748" style="fill:#ffffff;stroke:none;stroke-width:0.374906"/></svg>`
	svgImpact             = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon BreakStun" viewBox="0 0 14 14"><path d="m.336.336 4.613 7.688-2.563.768 2.819 1.794-4.613 3.075h7.944c1.506 0 2.923.11 4.079-1.046s1.046-2.573 1.046-4.079V.592c-1.236 1.654-2.164 3.17-3.075 4.613L9.17 2.32 8.024 4.95zm6.919 6.663 2.819 1.793 1.281-1.537c0 1.153.26 2.724-.69 3.507-.949.783-2.443.575-3.41.593l1.537-1.281z" style="fill:#ffffff;stroke:none;stroke-width:0.256258"/></svg>`
	svgCritRate           = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon Crit" viewBox="0 0 14 14"><path d="M3.324 1.067C3.546 3.702 2.304 5.637.13 7c2.415 1.075 3.195 3.484 3.195 5.934 0 0 1.592-1.104 3.651-1.104 2.06 0 3.651 1.104 3.651 1.104-.362-2.61 1.078-4.61 3.195-5.934-2.415-1.076-3.195-3.484-3.195-5.933 0 0-1.627 1.102-3.65 1.102-2.025 0-3.652-1.102-3.652-1.102m9.53 1.111a1.12 1.12 0 1 0 0 2.241 1.12 1.12 0 0 0 0-2.24m-3.66 1.494c.296 1.24.846 2.43 1.76 3.41-.865.906-1.531 1.979-1.866 3.219a8.4 8.4 0 0 0-2.035-.274c-.7 0-1.437.115-2.14.303-.298-1.241-.85-2.432-1.763-3.41.865-.92 1.491-1.99 1.82-3.232a8.3 8.3 0 0 0 2.083.293 8.3 8.3 0 0 0 2.14-.309" style="fill:#ffffff;stroke:none;stroke-width:0.228206"/><path d="M7 4.988s.32.985.674 1.338C8.027 6.68 9.012 7 9.012 7s-.985.32-1.338.674C7.32 8.027 7 9.012 7 9.012s-.32-.985-.674-1.338C5.973 7.32 4.988 7 4.988 7s.985-.32 1.338-.674C6.68 5.973 7 4.988 7 4.988" style="fill:#ffffff;stroke:none;stroke-width:0.166276"/></svg>`
	svgCritDMG            = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon CritDmg" viewBox="0 0 14 14"><path d="M3.267.933C3.493 3.628 2.224 5.607 0 7c2.47 1.1 3.267 3.562 3.267 6.067 0 0 1.627-1.128 3.733-1.128s3.733 1.128 3.733 1.128C10.363 10.398 11.836 8.355 14 7c-2.47-1.1-3.267-3.562-3.267-6.067 0 0-1.663 1.128-3.733 1.128S3.267.933 3.267.933M7 4.177s.45 1.382.945 1.878C8.441 6.55 9.823 7 9.823 7s-1.382.45-1.878.945C7.45 8.441 7 9.823 7 9.823s-.45-1.382-.945-1.878C5.559 7.45 4.177 7 4.177 7s1.382-.45 1.878-.945C6.55 5.559 7 4.177 7 4.177" style="fill:#ffffff;stroke:none;stroke-width:0.233336"/></svg>`
	svgPENRatio           = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon PenRatio" viewBox="0 0 14 14"><path d="M12.207.13a1.305 1.305 0 0 0-1.304 1.304 1.305 1.305 0 0 0 1.304 1.304 1.305 1.305 0 0 0 1.305-1.304A1.305 1.305 0 0 0 12.207.129M9.331.412 6.814 3.039.628.857l2.324 6.051L.36 9.463l9.413.269Zm1.52 3.21.13 2.764a3.46 3.46 0 0 1 .645 2.006 3.46 3.46 0 0 1-3.462 3.463A3.46 3.46 0 0 1 5.891 11l-2.578-.073a5.48 5.48 0 0 0 4.851 2.943 5.477 5.477 0 0 0 5.477-5.477 5.48 5.48 0 0 0-2.79-4.77Z" style="fill:#ffffff;stroke:none;stroke-width:0.268938"/></svg>`
	svgPENFlat            = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon PenDelta" viewBox="0 0 14 14"><path d="M9.388.105 6.81 2.795.47.56l2.382 6.2L.195 9.38l9.646.275Zm1.558 3.289.134 2.833c.428.6.66 1.318.66 2.056a3.547 3.547 0 0 1-3.547 3.547 3.55 3.55 0 0 1-2.33-.875l-2.641-.076a5.61 5.61 0 0 0 4.971 3.016 5.612 5.612 0 0 0 2.753-10.5Z" style="fill:#ffffff;stroke:none;stroke-width:0.275589"/></svg>`
	svgEnergyRegen        = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon SpRecover" viewBox="0 0 14 14"><path d="m6.997-.003-2 3.972A2.33 2.33 0 0 1 3.965 5L-.007 7l3.972 2c.446.224.807.585 1.032 1.03l2 3.973 2-3.972A2.33 2.33 0 0 1 10.026 9L14 7l-3.973-2a2.33 2.33 0 0 1-1.03-1.031Zm5.126.53-1.884 2.49h1.224v1.179h1.32v-1.18h1.224zM6.997 4.865q.032.001.06.071c.117.324.392 1.02.688 1.315.296.297.992.572 1.316.69.094.033.094.085 0 .119-.324.118-1.02.393-1.316.69s-.571.991-.689 1.315c-.034.094-.085.094-.12 0-.117-.325-.392-1.02-.688-1.316s-.992-.572-1.316-.689c-.094-.034-.094-.085 0-.12.324-.117 1.02-.393 1.316-.689s.572-.992.689-1.315q.026-.072.06-.07" style="fill:#ffffff;stroke-width:0.894903;stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`
	svgAnomalyProficiency = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon ElementMystery" viewBox="0 0 14 14"><path d="M73.284 143.419a7.267 7.267 0 1 0 0 14.534 7.267 7.267 0 0 0 0-14.534m-.664 1.878a1.66 1.66 0 0 0-1.004 1.51c.001.207.041.412.118.605h-.004l.018.039q.024.053.05.105l.873 1.852a.9.9 77.387 0 1 .084.376v.14a.24.24 135 0 1-.24.24h-.142a.9.9 12.615 0 1-.376-.084l-1.796-.846a2 2 0 0 0-.194-.092l-.005-.002a1.7 1.7 0 0 0-.604-.114 1.66 1.66 0 0 0-1.502.984 5.44 5.44 0 0 1 4.724-4.713m1.349.002a5.44 5.44 0 0 1 4.703 4.717 1.66 1.66 0 0 0-1.518-.99 1.7 1.7 0 0 0-.604.117v-.004l-.04.019q-.052.023-.105.05l-1.852.872a.9.9 167.387 0 1-.375.084h-.14a.24.24 45 0 1-.24-.24v-.142a.9.9 102.615 0 1 .084-.375l.847-1.797q.051-.094.091-.194l.002-.005a1.66 1.66 0 0 0-.853-2.113zm.068 5.908h.142a.9.9 12.615 0 1 .376.084l1.797.847q.093.051.194.091l.005.002c.192.075.397.114.604.114v.001a1.66 1.66 0 0 0 1.517-.99 5.44 5.44 0 0 1-4.704 4.717c.59-.27.969-.86.969-1.51a1.7 1.7 0 0 0-.118-.603h.004l-.019-.04-.05-.105-.872-1.852a.9.9 77.386 0 1-.084-.375v-.141a.24.24 135 0 1 .24-.24m-1.663 0h.141a.24.24 45 0 1 .24.24v.142a.9.9 102.615 0 1-.085.376l-.846 1.796a2 2 0 0 0-.092.194l-.002.005h.001a1.7 1.7 0 0 0-.115.604 1.66 1.66 0 0 0 1.006 1.511 5.44 5.44 0 0 1-4.726-4.713 1.66 1.66 0 0 0 1.502.984c.207 0 .412-.04.604-.117v.004l.04-.019q.052-.023.105-.05l1.852-.872a.9.9 167.385 0 1 .375-.084" style="fill:#ffffff;stroke-width:1;stroke-linecap:square;stroke-dashoffset:3.77952" transform="translate(-60.462 -131.715)scale(.92056)"/></svg>`
	svgAnomalyMastery     = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon ElementAbnormalPower" viewBox="0 0 14 14"><path d="M7 .287a6.69 6.69 0 0 0-6.69 6.69A6.69 6.69 0 0 0 7 13.667a6.69 6.69 0 0 0 6.69-6.69A6.69 6.69 0 0 0 7 .287m0 1.739a4.95 4.95 0 0 1 4.952 4.951A4.95 4.95 0 0 1 7 11.93a4.95 4.95 0 0 1-4.952-4.952A4.95 4.95 0 0 1 7 2.026m-2.175.98a.2.2 0 0 0-.204.199v1.733a.3.3 0 0 0 .145.26l2.089 1.296a.27.27 0 0 0 .29 0l2.09-1.296a.3.3 0 0 0 .144-.26V3.205a.2.2 0 0 0-.304-.17L7 4.323 4.925 3.035a.2.2 0 0 0-.1-.029m-.632 2.888a.3.3 0 0 0-.148.04L2.544 6.8a.2.2 0 0 0 .005.348l2.153 1.154-.078 2.44a.2.2 0 0 0 .299.179l1.5-.867a.3.3 0 0 0 .153-.254l.078-2.46a.27.27 0 0 0-.144-.25L4.342 5.93a.3.3 0 0 0-.15-.035Zm5.614 0a.3.3 0 0 0-.149.034L7.49 7.09a.27.27 0 0 0-.144.25l.077 2.46a.3.3 0 0 0 .153.254l1.5.867a.2.2 0 0 0 .3-.18l-.078-2.44 2.153-1.153a.2.2 0 0 0 .005-.349l-1.5-.866a.3.3 0 0 0-.149-.04Z" style="fill:#ffffff;stroke-width:1.21664;stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`
	svgSheerForce         = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon SkipDefAtk" viewBox="0 0 14 14"><path d="M.31.31s-.454 7.384 2.418 12.998c.646.581 2.901-1.92 2.901-1.92s5.862 2.415 5.951 2.325c.09-.09-4.287-5.055-4.287-5.055s-2.758.06-2.965.021C2.79 5.98 3.283 3.283 3.283 3.283S5.98 2.79 8.679 4.328c.038.207-.02 2.965-.02 2.965s4.965 4.377 5.054 4.287-2.325-5.95-2.325-5.95 2.501-2.256 1.92-2.902C7.694-.144.31.31.31.31" style="fill:#ffffff;fill-opacity:1;stroke-width:1.16745;stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`
	svgPhysicalDMG        = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon AddedDamageRatio_Physics" viewBox="0 0 14 14"><defs><linearGradient href="#zzz-AddedDamageRatio_Physics__a" id="zzz-AddedDamageRatio_Physics__b" x1="12.046" x2="12.046" y1="278.603" y2="299.007" gradientTransform="translate(-1.265 -191.16)scale(.68614)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_Physics__a"><stop offset="0" style="stop-color:#e78801;stop-opacity:1"/><stop offset="1" style="stop-color:#efd400;stop-opacity:1"/></linearGradient></defs><path d="M9.013.217c-.87 1.096-2.116 2.47-3.595 2.36-.822-.06-2.224-.78-2.348-.663-.119.148.464 1.196.444 1.806C3.442 5.958 0 6.851 0 6.971c0 .127 1.243.482 1.805.865 1.862 1.266.612 2.894.182 4.554 1.503-.514 2.918-1.891 4.504-.761.8.57 1.173 1.273 1.817 2.186.755-2.82.996-3.654 2.71-4.18.897-.277 2.053-.178 2.982-.178-.45-.435-.988-1.001-1.404-1.583-1.44-2.013.221-3.08 1.311-4.696-1.88.059-3.876.534-4.562-1.806-.106-.36-.128-.776-.14-1.148 0 0-.058-.037-.092-.039-.034-.001-.1.034-.1.034zm-.986 3c.213 1.86 1.299 1.586 2.449 1.587 0 0-.813.657-.813 1.625s.813 1.626.813 1.626c-1.877-.12-2.618.594-2.752 2.273C6.78 8.834 6.3 8.476 4.515 9.319c.848-1.776.097-2.214-1.344-2.606 1.176-.488 2.515-1.09 1.886-2.632 1.565.642 2.184.209 2.97-.863" style="fill:url(#zzz-AddedDamageRatio_Physics__b);stroke:none;stroke-width:0.686135"/></svg>`
	svgFireDMG            = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon AddedDamageRatio_Fire" viewBox="0 0 14 14"><defs><linearGradient href="#zzz-AddedDamageRatio_Fire__a" id="zzz-AddedDamageRatio_Fire__b" x1="12.182" x2="12.182" y1="302.124" y2="315.03" gradientTransform="translate(-4.909 -294.668)scale(.9776)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_Fire__a"><stop offset="0" style="stop-color:#ea1503;stop-opacity:1"/><stop offset="1" style="stop-color:#f3741a;stop-opacity:1"/></linearGradient></defs><path d="M1.972 11.944c1.112 1.072 2.998 1.314 4.368 2.027.007-1.742-1.044-3.54-3.158-3.855 2.096-.47 3.1-1.295 3.665-3.693.572 2.723 1.85 3.219 3.91 3.678-2.092.376-3.42 1.696-3.159 3.899 1.55-.806 3.624-1.29 4.76-2.677 2.117-2.58.435-6.764-2.47-7.9.949 2.315-1.854 2.066-2.105-.375C7.666 1.909 8.886.993 9.77.47 7.84-.68 4.008.404 4.242 2.366c.191 1.6 1.52 3.367-.155 3.766-1.715.409-1.2-2.104-1.2-2.104C.665 5.06-.254 9.8 1.973 11.945" style="fill:url(#zzz-AddedDamageRatio_Fire__b);fill-opacity:1;stroke:none;stroke-width:0.977604"/></svg>`
	svgIceDMG             = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon AddedDamageRatio_Ice" viewBox="0 0 14 14"><defs><linearGradient href="#zzz-AddedDamageRatio_Ice__a" id="zzz-AddedDamageRatio_Ice__b" x1="12.046" x2="12.046" y1="318.508" y2="331.879" gradientTransform="translate(-3.923 -287.884)scale(.9068)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_Ice__a"><stop offset="0" style="stop-color:#04c2c8;stop-opacity:1"/><stop offset="1" style="stop-color:#83f4f0;stop-opacity:1"/></linearGradient></defs><path d="M7 0 5.166 3.824.938 3.5 3.332 7 .938 10.5l4.228-.323L7 14l1.834-3.823 4.228.323L10.668 7l2.394-3.5-4.228.324ZM5.18 5.06a.1.1 0 0 1 .06.015L7 5.99l1.76-.914a.123.123 0 0 1 .166.167l-.914 1.76.914 1.759a.123.123 0 0 1-.166.166L7 8.013l-1.76.914a.123.123 0 0 1-.166-.166L5.988 7l-.914-1.76a.123.123 0 0 1 .106-.18" style="fill:url(#zzz-AddedDamageRatio_Ice__b);stroke-width:0.906792;stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`
	svgElectricDMG        = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon AddedDamageRatio_Elec" viewBox="0 0 14 14"><defs><linearGradient href="#zzz-AddedDamageRatio_Elec__a" id="zzz-AddedDamageRatio_Elec__b" x1="12.046" x2="12.046" y1="334.813" y2="349.957" gradientTransform="translate(-4.136 -309.526)scale(.92447)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_Elec__a"><stop offset="0" style="stop-color:#0075ff;stop-opacity:1"/><stop offset="1" style="stop-color:#3decff;stop-opacity:1"/></linearGradient></defs><path d="m9.822.624-5.237.624-1.573 5.236 2.143-.054-1.736 6.946 6.783-8.628-2.768.217ZM1.6 5.779 0 7.705l.678 2.578.652-.57.949 2.496.597-4.531-1.167.814Zm9.17 1.004L8.14 9.225l1.166.298-2.7 3.392 6.065-2.768-1.764-.786L14 7Z" style="fill:url(#zzz-AddedDamageRatio_Elec__b);stroke-width:0.924468;stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`
	svgEtherDMG           = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon AddedDamageRatio_Ether" viewBox="0 0 14 14"><defs><linearGradient href="#zzz-AddedDamageRatio_Ether__a" id="zzz-AddedDamageRatio_Ether__b" x1="16.461" x2="7.36" y1="354.589" y2="367.587" gradientTransform="translate(-2.93 -294.054)scale(.83374)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_Ether__a"><stop offset="0" style="stop-color:#ff0a1a;stop-opacity:1"/><stop offset=".171" style="stop-color:#ff0626;stop-opacity:1"/><stop offset=".5" style="stop-color:#b338dd;stop-opacity:1"/><stop offset=".85" style="stop-color:#2a6bea;stop-opacity:1"/><stop offset="1" style="stop-color:#2a6bea;stop-opacity:1"/></linearGradient></defs><path d="M6.52 0 5.036 3.715a1.68 1.68 0 0 1-.935.936L.385 6.135l3.716 1.483a1.47 1.47 0 0 1 .866.96L6.52 14l1.553-5.422a1.47 1.47 0 0 1 .866-.96l3.716-1.483L8.94 4.65a1.68 1.68 0 0 1-.935-.936Zm4.568 7.83-.57 1.403a1 1 0 0 1-.554.553l-1.403.57 1.403.57a1 1 0 0 1 .554.554l.57 1.404.57-1.404a1 1 0 0 1 .554-.553l1.403-.57-1.403-.57a1 1 0 0 1-.554-.554Z" style="fill:url(#zzz-AddedDamageRatio_Ether__b);stroke-width:0.833738;stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`
	svgWindDMG            = `<svg xmlns="http://www.w3.org/2000/svg" class="SvgIcon AddedDamageRatio_Wind" viewBox="0 0 14 14"><defs><linearGradient id="zzz-AddedDamageRatio_Wind__a"><stop offset="0" style="stop-color:#61a3ff;stop-opacity:1"/><stop offset="1" style="stop-color:#97e3fa;stop-opacity:1"/></linearGradient><linearGradient href="#zzz-AddedDamageRatio_Wind__a" id="zzz-AddedDamageRatio_Wind__b" x1="1862.372" x2="1908.889" y1="1828.644" y2="1900.273" gradientUnits="userSpaceOnUse"/><linearGradient href="#zzz-AddedDamageRatio_Wind__a" id="zzz-AddedDamageRatio_Wind__c" x1="1862.372" x2="1908.889" y1="1828.644" y2="1900.273" gradientUnits="userSpaceOnUse"/></defs><g style="fill:url(#zzz-AddedDamageRatio_Wind__b)"><path fill="#FFF" d="M1870.45 1829q10.4-.95 17.85.6 6.65 1.4 10.4 4.6-12.85-5.45-27.25-1.45-15.65 4.35-24.5 18.15-5.25 8.15-3.65 18.25 2-7.9 6.4-13.65 4.9-6.45 12.35-9.65-23.2 18.7-11.3 39.95 7.95 14.1 25.85 16.75 16.85 2.45 29.45-6.4-13.55 1.95-22.2.05-10.3-2.25-16.55-10.35 9.2 7.1 22 6.35 11.7-.7 22.2-7.5 10.55-6.85 14.7-16.65 4.55-10.75-.6-21.1-.789 17.455-14.05 25 9.75-10.45 9.55-20.85-.2-9.35-8.05-16.3-7.6-6.65-19.1-8.4-12.1-1.9-23.5 2.6m26.35 43.15q-7.2 5.35-15.8 5.35-8.55 0-13.55-5.35-4.9-5.4-3.25-12.9 1.65-7.55 8.95-13 7.25-5.3 15.85-5.3 8.55 0 13.45 5.3 5 5.45 3.35 13-1.65 7.5-9 12.9" style="fill:url(#zzz-AddedDamageRatio_Wind__c)" transform="translate(-299.382 -295.987)scale(.16249)"/></g></svg>`
)

var propertySVGMap = map[PropertyID]string{
	PropBaseHP:         svgHP,
	PropHPPercent:      svgHP,
	PropHPFlat:         svgHP,
	PropHPPercentBonus: svgHP,
	PropHPFlatBonus:    svgHP,

	PropBaseATK:    svgATK,
	PropATKPercent: svgATK,
	PropATKFlat:    svgATK,

	PropBaseDEF:    svgDEF,
	PropDEFPercent: svgDEF,
	PropDEFFlat:    svgDEF,

	PropBaseImpact:    svgImpact,
	PropImpactPercent: svgImpact,
	PropImpactFlat:    svgImpact,

	PropBaseCritRate: svgCritRate,
	PropCritRate:     svgCritRate,

	PropBaseCritDMG: svgCritDMG,
	PropCritDMG:     svgCritDMG,

	PropBasePENRatio: svgPENRatio,
	PropPENRatio:     svgPENRatio,

	PropBasePENFlat: svgPENFlat,
	PropPENFlat:     svgPENFlat,

	PropBaseEnergyRegen:    svgEnergyRegen,
	PropEnergyRegenPercent: svgEnergyRegen,
	PropEnergyRegen:        svgEnergyRegen,

	PropBaseAnomalyProficiency:    svgAnomalyProficiency,
	PropAnomalyProficiencyPercent: svgAnomalyProficiency,
	PropAnomalyProficiency:        svgAnomalyProficiency,

	PropBaseAnomalyMastery:    svgAnomalyMastery,
	PropAnomalyMasteryPercent: svgAnomalyMastery,
	PropAnomalyMastery:        svgAnomalyMastery,

	PropBaseSheerForce: svgSheerForce,
	PropSheerForce:     svgSheerForce,

	PropPhysicalDMGBonus: svgPhysicalDMG,
	PropFireDMGBonus:     svgFireDMG,
	PropIceDMGBonus:      svgIceDMG,
	PropElectricDMGBonus: svgElectricDMG,
	PropEtherDMGBonus:    svgEtherDMG,
	PropWindDMGBonus:     svgWindDMG,
}
