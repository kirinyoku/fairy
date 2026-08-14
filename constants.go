package fairy

const (
	// EnkaAssetBaseURL is the base URL for Enka Network Zenless Zone Zero UI assets.
	EnkaAssetBaseURL = "https://enka.network/ui/zzz/"

	// StatModifierScale is the fixed divisor used by ZZZ data tables to convert integer stat growth/modifiers into float multipliers.
	StatModifierScale = 10000.0

	// DefaultBaseCritRate is the default base critical hit rate for all agents (5%).
	DefaultBaseCritRate = 0.05

	// DefaultBaseCritDMG is the default base critical hit damage for all agents (50%).
	DefaultBaseCritDMG = 0.50

	// DefaultBasePenRatio is the default base penetration ratio for all agents (0%).
	DefaultBasePenRatio = 0.0

	// DefaultBasePenFlat is the default base flat penetration for all agents (0).
	DefaultBasePenFlat = 0.0
)

// Skill slot index constants matching ZZZ character skill mapping slots.
const (
	SkillIndexBasic    = 0 // Basic Attack slot
	SkillIndexDodge    = 1 // Dodge slot
	SkillIndexAssist   = 2 // Assist slot
	SkillIndexSpecial  = 3 // Special Attack slot
	SkillIndexChainAlt = 4 // Alternate Chain Attack slot
	SkillIndexChain    = 6 // Primary Chain Attack slot
)
