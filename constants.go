package fairy

const (
	// EnkaAssetBaseURL is the base URL for Enka Network Zenless Zone Zero UI assets.
	EnkaAssetBaseURL = "https://enka.network/ui/zzz/"

	// statModifierScale is the fixed divisor used by ZZZ data tables to convert integer stat growth/modifiers into float multipliers.
	statModifierScale = 10000.0

	// defaultBaseCritRate is the default base critical hit rate for all agents (5%).
	defaultBaseCritRate = 0.05

	// defaultBaseCritDMG is the default base critical hit damage for all agents (50%).
	defaultBaseCritDMG = 0.50

	// defaultBasePenRatio is the default base penetration ratio for all agents (0%).
	defaultBasePenRatio = 0.0

	// defaultBasePenFlat is the default base flat penetration for all agents (0).
	defaultBasePenFlat = 0.0
)

// Skill slot index constants matching ZZZ character skill mapping slots.
const (
	skillIndexBasic    = 0 // Basic Attack slot
	skillIndexDodge    = 1 // Dodge slot
	skillIndexAssist   = 2 // Assist slot
	skillIndexSpecial  = 3 // Special Attack slot
	skillIndexChainAlt = 4 // Alternate Chain Attack slot
	skillIndexChain    = 6 // Primary Chain Attack slot
)
