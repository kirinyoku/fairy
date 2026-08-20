package fairy

const (
	// EnkaAssetBaseURL is the base HTTPS URL for Enka.Network Zenless Zone Zero CDN UI assets
	// (such as agent splash arts, icons, badges, and W-Engine textures).
	EnkaAssetBaseURL = "https://enka.network/ui/zzz/"

	// statModifierScale is the fixed divisor used by ZZZ data tables to convert integer stat growth/modifiers into float multipliers (10000 = 100%).
	statModifierScale = 10000.0

	// defaultBaseCritRate is the default innate base critical hit rate for all agents (5.0%).
	defaultBaseCritRate = 0.05

	// defaultBaseCritDMG is the default innate base critical hit damage for all agents (50.0%).
	defaultBaseCritDMG = 0.50

	// defaultBasePenRatio is the default innate base penetration ratio for all agents (0.0%).
	defaultBasePenRatio = 0.0

	// defaultBasePenFlat is the default innate base flat penetration for all agents (0).
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
