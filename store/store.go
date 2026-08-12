// Package store provides an abstraction over the Zenless Zone Zero datamined game data.
// It serves as a unified repository for all stat scaling formulas, localization text,
// and item definitions needed to reconstruct a full player profile from raw API IDs.
package store

// MetadataStore defines the interface for fetching game metadata and localizations.
// It acts as an abstraction layer so the underlying data source (embedded JSONs, remote CDN, etc.)
// can be easily swapped or mocked.
type MetadataStore interface {
	// Localize translates a text hash into the specified language.
	Localize(hash string, lang string) string

	// AvatarMeta returns the metadata for a specific agent (avatar) by its internal ID.
	// Returns false if the agent is not found.
	AvatarMeta(id int) (AvatarMeta, bool)

	// AvatarSkillsMeta returns the skills metadata for a specific agent by its internal ID.
	AvatarSkillsMeta(avatarID int) ([]SkillMeta, bool)

	// AvatarMindscapesMeta returns the Mindscape Cinema metadata for a specific agent by its internal ID.
	AvatarMindscapesMeta(avatarID int) ([]MindscapeMeta, bool)

	// WeaponMeta returns the metadata for a specific W-Engine by its internal ID.
	WeaponMeta(id int) (WeaponMeta, bool)

	// EquipmentMeta returns the metadata for a specific Drive Disc by its internal ID.
	EquipmentMeta(id int) (EquipmentMeta, bool)

	// EquipmentSuitMeta returns the metadata for a Drive Disc Set (suit) by its ID.
	EquipmentSuitMeta(suitID int) (EquipmentSuitMeta, bool)

	// PropertyMeta returns the metadata for a combat property (e.g., ATK, CRIT) by its ID.
	PropertyMeta(id int) (PropertyMeta, bool)

	// MedalMeta returns the metadata for a profile medal/badge by its ID.
	MedalMeta(id int) (MedalMeta, bool)

	// TitleMeta returns the metadata for a player title by its ID.
	TitleMeta(id int) (TitleMeta, bool)

	// NamecardMeta returns the metadata for a profile namecard by its ID.
	NamecardMeta(id int) (NamecardMeta, bool)

	// PfpMeta returns the metadata for a profile picture (avatar icon) by its ID.
	PfpMeta(id int) (PfpMeta, bool)

	// SkinMeta returns the metadata for an agent skin by its ID.
	SkinMeta(id int) (SkinMeta, bool)

	// DefaultSkinMeta returns the default skin metadata for an agent by their avatar ID.
	DefaultSkinMeta(avatarID int) (SkinMeta, bool)

	// WeaponLevelTemplate returns the multiplier template for a weapon's main stat at a given level.
	WeaponLevelTemplate(rarity, level int) (WeaponLevelTemplate, bool)

	// WeaponStarTemplate returns the multiplier template for a weapon's stats based on refinement phase.
	WeaponStarTemplate(rarity, phase int) (WeaponStarTemplate, bool)

	// EquipmentLevelTemplate returns the multiplier template for a drive disc's main stat.
	EquipmentLevelTemplate(rarity, level int) (EquipmentLevelTemplate, bool)
}
