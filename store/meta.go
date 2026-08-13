package store

// AvatarMeta represents immutable, normalized avatar metadata.
// It contains the character's base definitions and stat scaling arrays.
type AvatarMeta struct {
	// Name is a localization hash key (e.g., "AvatarName_1191"), NOT the actual localized name string.
	Name string
	// Rarity indicates the base tier (4=S-rank, 3=A-rank, 2=B-rank).
	Rarity         int
	ProfessionType string
	Camp           string
	ElementTypes   []string
	// BaseProps maps PropertyID to the base stat value at Level 1.
	BaseProps map[int]int
	// GrowthProps maps PropertyID to the stat growth coefficient per level.
	GrowthProps map[int]int
	// PromotionProps contains flat stat additions for ascension tiers.
	// It is indexed by (promotion level - 1).
	PromotionProps []map[int]int
	// CoreEnhancementProps contains flat stat additions for core skill upgrades.
	// It is indexed by core enhancement level (0-based). Level 0 means "no enhancement".
	CoreEnhancementProps []map[int]int
	Image                string
	CircleIcon           string
}

// BaseStat returns the base property value for the given propID.
func (a *AvatarMeta) BaseStat(propID int) (int, bool) {
	val, ok := a.BaseProps[propID]
	return val, ok
}

// GrowthStat returns the growth property value for the given propID.
func (a *AvatarMeta) GrowthStat(propID int) (int, bool) {
	val, ok := a.GrowthProps[propID]
	return val, ok
}

// PromotionStat returns the promotion property value for the given promotionLevel and propID.
func (a *AvatarMeta) PromotionStat(promotionLevel, propID int) (int, bool) {
	if promotionLevel > 0 && promotionLevel-1 < len(a.PromotionProps) {
		val, ok := a.PromotionProps[promotionLevel-1][propID]
		return val, ok
	}
	return 0, false
}

// CoreEnhancementStat returns the core enhancement property value.
func (a *AvatarMeta) CoreEnhancementStat(coreEnhancement, propID int) (int, bool) {
	if coreEnhancement > 0 && coreEnhancement < len(a.CoreEnhancementProps) {
		val, ok := a.CoreEnhancementProps[coreEnhancement][propID]
		return val, ok
	}
	return 0, false
}

// PropertyStat represents a stat in WeaponMeta.
// It is a key-value pair extracted directly from the weapon base data.
type PropertyStat struct {
	PropertyID    int `json:"PropertyId"`
	PropertyValue int `json:"PropertyValue"`
}

// WeaponMeta represents immutable weapon metadata.
type WeaponMeta struct {
	ItemName       string
	Rarity         int
	ProfessionType string
	MainStat       PropertyStat
	SecondaryStat  PropertyStat
	ImagePath      string
	// PassiveDescKeys contains localization keys for the weapon's passive ability description.
	PassiveDescKeys []string
}

// EquipmentMeta represents equipment (drive disc) metadata.
type EquipmentMeta struct {
	Rarity int
	SuitID int
}

// EquipmentSuitMeta represents a drive disc set (suit).
type EquipmentSuitMeta struct {
	Icon          string
	Name          string
	SetBonusProps map[int]int
	// Set2DescKey is the localization key for the 2-piece set bonus description.
	Set2DescKey string
	// Set4DescKey is the localization key for the 4-piece set bonus description.
	Set4DescKey string
}

// SkillMeta represents avatar skill metadata.
type SkillMeta struct {
	NameKey string `json:"NameKey"`
	DescKey string `json:"DescKey"`
}

// PropertyMeta represents property metadata.
type PropertyMeta struct {
	Name   string `json:"Name"`
	Format string `json:"Format"`
}

// templateKey is used as a composite map key for caching template lookups.
// It combines item Rarity and Level (or Phase) for O(1) retrieval.
type templateKey struct {
	Rarity int
	Level  int
}

// WeaponLevelTemplate represents the multiplier for weapon's main stat at a given level.
// IMPORTANT: The obfuscated JSON field tags (e.g., "ICPMKHFGPOG") correspond directly
// to the minified keys in the datamined game data files. These keys change almost
// every major patch and must be updated alongside new data dumps.
type WeaponLevelTemplate struct {
	Rarity             int `json:"ICPMKHFGPOG"`
	Level              int `json:"EMLFBEMHINK"`
	MainStat           int `json:"AHMDJCIHNKG"` // Multiplier per 10000
	SubStatDenominator int `json:"IDBKOAPHGLC"` // Divisor for substat
}

// WeaponStarTemplate represents the multipliers for weapon's stats based on refinement phase.
// IMPORTANT: The obfuscated JSON field tags correspond to datamined keys and are patch-dependent.
type WeaponStarTemplate struct {
	Rarity     int `json:"ICPMKHFGPOG"`
	BreakLevel int `json:"BBOCBHBGMML"`
	MainStat   int `json:"NMFHJKEFLOG"` // Multiplier per 10000
	SubStat    int `json:"FCLIIPBDDKP"` // Multiplier per 10000
}

// EquipmentLevelTemplate represents the multiplier for drive disc's main stat.
// IMPORTANT: The obfuscated JSON field tags correspond to datamined keys and are patch-dependent.
type EquipmentLevelTemplate struct {
	Rarity   int `json:"GMKDLJLLBPO"`
	Level    int `json:"FNPIELBFDEJ"`
	MainStat int `json:"JEKGLLBALFE"` // Multiplier per 10000
}

// MedalMeta represents a medal.
type MedalMeta struct {
	Name       string `json:"Name"`
	Icon       string `json:"Icon"`
	TipNum     string `json:"TipNum"`
	PrefixIcon string `json:"PrefixIcon"`
}

// TitleMeta represents a title.
type TitleMeta struct {
	TitleText string `json:"TitleText"`
	ColorA    string `json:"ColorA"`
	ColorB    string `json:"ColorB"`
}

// NamecardMeta represents a namecard.
type NamecardMeta struct {
	Icon string `json:"Icon"`
}

// PfpMeta represents a profile picture.
type PfpMeta struct {
	Icon string `json:"Icon"`
}

// SkinMeta represents a skin.
type SkinMeta struct {
	Id        int    `json:"Id"`
	AvatarId  int    `json:"AvatarId"`
	NameKey   string `json:"NameKey"`
	DescKey   string `json:"DescKey"`
	Icon      string `json:"Icon"`
	IsDefault bool   `json:"IsDefault"`
}

type avatarRecord struct {
	Name                 string           `json:"Name"`
	Rarity               int              `json:"Rarity"`
	ProfessionType       string           `json:"ProfessionType"`
	Camp                 string           `json:"Camp"`
	ElementTypes         []string         `json:"ElementTypes"`
	BaseProps            map[string]int   `json:"BaseProps"`
	GrowthProps          map[string]int   `json:"GrowthProps"`
	PromotionProps       []map[string]int `json:"PromotionProps"`
	CoreEnhancementProps []map[string]int `json:"CoreEnhancementProps"`
	Image                string           `json:"Image"`
	CircleIcon           string           `json:"CircleIcon"`
}

type weaponRecord struct {
	ItemName        string       `json:"ItemName"`
	Rarity          int          `json:"Rarity"`
	ProfessionType  string       `json:"ProfessionType"`
	MainStat        PropertyStat `json:"MainStat"`
	SecondaryStat   PropertyStat `json:"SecondaryStat"`
	ImagePath       string       `json:"ImagePath"`
	PassiveDescKeys []string     `json:"PassiveDescKeys"`
}

type equipmentRecord struct {
	Rarity int `json:"Rarity"`
	SuitId int `json:"SuitId"`
}

type equipmentSuitRecord struct {
	Icon          string         `json:"Icon"`
	Name          string         `json:"Name"`
	SetBonusProps map[string]int `json:"SetBonusProps"`
	Set2DescKey   string         `json:"Set2DescKey"`
	Set4DescKey   string         `json:"Set4DescKey"`
}

type equipmentContainer struct {
	Items map[string]equipmentRecord     `json:"Items"`
	Suits map[string]equipmentSuitRecord `json:"Suits"`
}

type weaponLevelContainer struct {
	List []WeaponLevelTemplate `json:"OOFFGGKCDID"`
}

type weaponStarContainer struct {
	List []WeaponStarTemplate `json:"OOFFGGKCDID"`
}

type equipmentLevelContainer struct {
	List []EquipmentLevelTemplate `json:"MIJCMCEDADM"`
}

type titleContainer struct {
	Titles        map[string]TitleMeta `json:"Titles"`
	TitleVariants map[string]TitleMeta `json:"TitleVariants"`
}

// MindscapeMeta represents immutable metadata for a Mindscape Cinema node.
type MindscapeMeta struct {
	Rank     int    `json:"Rank"`
	TitleKey string `json:"TitleKey"`
	DescKey  string `json:"DescKey"`
}

// PotentialVisionMeta represents immutable metadata for a Potential Vision node.
type PotentialVisionMeta struct {
	ID           int    `json:"ID"`
	Level        int    `json:"Level"`
	LevelNameKey string `json:"LevelNameKey"`
	TitleKey     string `json:"TitleKey"`
	DescKey      string `json:"DescKey"`
}
