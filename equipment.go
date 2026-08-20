package fairy

// WEngine represents an enriched W-Engine (weapon) equipped by an [Agent].
// It provides base combat stats, scaling secondary stats, and a unique passive ability.
type WEngine struct {
	// ID is the internal numeric identifier of the W-Engine.
	ID int `json:"id"`

	// UID is the unique instance identifier of this specific W-Engine.
	UID string `json:"uid"`

	// Name is the localized display name of the W-Engine (e.g. "Deep Sea Visitor", "The Brimstone").
	Name string `json:"name"`

	// Level is the current progression level of the W-Engine (1–60).
	Level int `json:"level"`

	// Phase is the ascension / star phase of the W-Engine (0–5).
	Phase int `json:"phase"`

	// Modification is the refinement level of the W-Engine's passive skill (1–5 / M1–M5).
	Modification int `json:"modification"`

	// Rarity is the rarity tier of the W-Engine ([RarityS], [RarityA], or [RarityB]).
	Rarity Rarity `json:"rarity"`

	// Specialty is the recommended combat role for this W-Engine (e.g. [SpecialtyAttack], [SpecialtyStun]).
	Specialty Specialty `json:"specialty"`

	// SpecialtyName is the localized display name of the intended specialty.
	SpecialtyName string `json:"specialty_name"`

	// IconURL is the absolute HTTPS URL pointing to the W-Engine's visual icon on the EnkaNetwork CDN.
	IconURL string `json:"icon_url"`

	// MainStat is the primary stat provided by the W-Engine ([PropBaseATK]), scaled by level and phase.
	MainStat StatValue `json:"main_stat"`

	// SecondaryStat is the secondary stat provided by the W-Engine (e.g. CRIT Rate, ATK%, PEN Ratio), scaled by level.
	SecondaryStat StatValue `json:"secondary_stat"`

	// PassiveDescription is the localized description text of the W-Engine's passive ability.
	PassiveDescription string `json:"passive_description"`

	// FormattedHTML is the web-ready HTML description with inline CSS colors and evaluated modification parameters.
	FormattedHTML string `json:"formatted_html,omitempty"`
}

// FormatHTML returns the W-Engine passive description formatted as HTML with inline CSS styling.
// Returns an empty string if w is nil.
func (w *WEngine) FormatHTML() string {
	if w == nil {
		return ""
	}
	return formatHTML(w.PassiveDescription)
}

// FormatPlainText returns the W-Engine passive description as clean plain text with all Rich Text tags stripped.
// Returns an empty string if w is nil.
func (w *WEngine) FormatPlainText() string {
	if w == nil {
		return ""
	}
	return formatPlainText(w.PassiveDescription)
}

// FormatMarkdown returns the W-Engine passive description formatted with Markdown syntax.
// Returns an empty string if w is nil.
func (w *WEngine) FormatMarkdown() string {
	if w == nil {
		return ""
	}
	return formatMarkdown(w.PassiveDescription)
}

// SetID represents the unique numeric identifier of a Drive Disc set.
type SetID int

// Strongly-typed Drive Disc [SetID] constants for all equipment sets in Zenless Zone Zero.
const (
	// SetWoodpeckerElectro represents the Woodpecker Electro set.
	SetWoodpeckerElectro SetID = 31000
	// SetPufferElectro represents the Puffer Electro set.
	SetPufferElectro SetID = 31100
	// SetShockstarDisco represents the Shockstar Disco set.
	SetShockstarDisco SetID = 31200
	// SetFreedomBlues represents the Freedom Blues set.
	SetFreedomBlues SetID = 31300
	// SetHormonePunk represents the Hormone Punk set.
	SetHormonePunk SetID = 31400
	// SetSoulRock represents the Soul Rock set.
	SetSoulRock SetID = 31500
	// SetSwingJazz represents the Swing Jazz set.
	SetSwingJazz SetID = 31600
	// SetChaosJazz represents the Chaos Jazz set.
	SetChaosJazz SetID = 31800
	// SetProtoPunk represents the Proto Punk set.
	SetProtoPunk SetID = 31900
	// SetInfernoMetal represents the Inferno Metal set.
	SetInfernoMetal SetID = 32200
	// SetChaoticMetal represents the Chaotic Metal set.
	SetChaoticMetal SetID = 32300
	// SetThunderMetal represents the Thunder Metal set.
	SetThunderMetal SetID = 32400
	// SetPolarMetal represents the Polar Metal set.
	SetPolarMetal SetID = 32500
	// SetFangedMetal represents the Fanged Metal set.
	SetFangedMetal SetID = 32600
	// SetBranchBladeSong represents the Branch & Blade Song set.
	SetBranchBladeSong SetID = 32700
	// SetAstralVoice represents the Astral Voice set.
	SetAstralVoice SetID = 32800
	// SetShadowHarmony represents the Shadow Harmony set.
	SetShadowHarmony SetID = 32900
	// SetPhaethonsMelody represents the Phaethon's Melody set.
	SetPhaethonsMelody SetID = 33000
	// SetYunkuiTales represents the Yunkui Tales set.
	SetYunkuiTales SetID = 33100
	// SetKingOfTheSummit represents the King of the Summit set.
	SetKingOfTheSummit SetID = 33200
	// SetDawnsBloom represents the Dawn's Bloom set.
	SetDawnsBloom SetID = 33300
	// SetMoonlightLullaby represents the Moonlight Lullaby set.
	SetMoonlightLullaby SetID = 33400
	// SetWhiteWaterBallad represents the White Water Ballad set.
	SetWhiteWaterBallad SetID = 33500
	// SetShiningAria represents the Shining Aria set.
	SetShiningAria SetID = 33600
	// SetBunnyInWonderland represents the Bunny in Wonderland set.
	SetBunnyInWonderland SetID = 33700
	// SetNotesFromTheChained represents the Notes From the Chained set.
	SetNotesFromTheChained SetID = 33800
	// SetWutheringSalon represents the Wuthering Salon set.
	SetWutheringSalon SetID = 33900
	// SetTheSkyAblaze represents the The Sky Ablaze set.
	SetTheSkyAblaze SetID = 34000
	// SetFeatheredFate represents the Feathered Fate set.
	SetFeatheredFate SetID = 34100
	// SetThornedRose represents the Thorned Rose set.
	SetThornedRose SetID = 34200
)

var allSetIDs = [...]SetID{
	SetWoodpeckerElectro,
	SetPufferElectro,
	SetShockstarDisco,
	SetFreedomBlues,
	SetHormonePunk,
	SetSoulRock,
	SetSwingJazz,
	SetChaosJazz,
	SetProtoPunk,
	SetInfernoMetal,
	SetChaoticMetal,
	SetThunderMetal,
	SetPolarMetal,
	SetFangedMetal,
	SetBranchBladeSong,
	SetAstralVoice,
	SetShadowHarmony,
	SetPhaethonsMelody,
	SetYunkuiTales,
	SetKingOfTheSummit,
	SetDawnsBloom,
	SetMoonlightLullaby,
	SetWhiteWaterBallad,
	SetShiningAria,
	SetBunnyInWonderland,
	SetNotesFromTheChained,
	SetWutheringSalon,
	SetTheSkyAblaze,
	SetFeatheredFate,
	SetThornedRose,
}

// AllSetIDs returns a newly allocated slice containing all known [SetID] constants.
// The returned slice is a defensive copy and can be safely mutated by the caller.
func AllSetIDs() []SetID {
	sets := make([]SetID, len(allSetIDs))
	copy(sets, allSetIDs[:])
	return sets
}

// IsValid reports whether s is a recognized [SetID].
func (s SetID) IsValid() bool {
	for _, id := range allSetIDs {
		if s == id {
			return true
		}
	}
	return false
}

// Set represents a Drive Disc equipment set.
// A Set grants bonus combat effects when an [Agent] equips 2 or 4 pieces belonging to the same set.
type Set struct {
	// ID is the unique numeric identifier of the Drive Disc set (see [SetID] constants, e.g. [SetWoodpeckerElectro]).
	ID SetID `json:"id"`

	// Name is the localized display name of the set (e.g. "Woodpecker Electro", "Polar Metal", "Fanged Metal").
	Name string `json:"name"`
}

// DriveDisc represents an equipped Drive Disc on an [Agent].
//
// An [Agent] can equip up to 6 Drive Discs across partition slots 1 through 6:
//   - Slots 1–3 have fixed main stats (Slot 1: Flat HP, Slot 2: Flat ATK, Slot 3: Flat DEF).
//   - Slots 4–6 have randomized main stats (e.g. CRIT Rate, CRIT DMG, Attribute DMG Bonus, Energy Regen, Anomaly Mastery).
type DriveDisc struct {
	// ID is the internal numeric identifier of the specific disc variation.
	ID int `json:"id"`

	// UID is the unique instance identifier of this specific Drive Disc piece.
	UID string `json:"uid"`

	// Set is the Drive Disc [Set] metadata this disc belongs to.
	Set Set `json:"set"`

	// Slot is the equip partition slot position (1 to 6).
	Slot int `json:"slot"`

	// Level is the current upgrade level of the Drive Disc (0–15).
	Level int `json:"level"`

	// Rarity is the rarity rank of the disc ([RarityS], [RarityA], or [RarityB]).
	Rarity Rarity `json:"rarity"`

	// IconPath is the absolute HTTPS URL pointing to the disc's partition icon on the EnkaNetwork CDN.
	IconPath string `json:"icon_path"`

	// MainStat is the primary stat provided by this disc, scaled by disc level.
	MainStat StatValue `json:"main_stat"`

	// SubStats is the list of randomly rolled sub-stats (up to 4), including upgrade roll counts.
	SubStats []StatValue `json:"sub_stats"`
}

// DriveDiscs represents a collection of [DriveDisc] entries equipped on an [Agent] or stored in inventory.
type DriveDiscs []DriveDisc

// SubStatTotals aggregates and sums sub-stat values and roll counts across all discs in the collection.
// It groups them by [PropertyID] and preserves the deterministic appearance order of the sub-stats.
//
// Example:
//
//	totals := agent.DriveDiscs.SubStatTotals()
//	for _, stat := range totals {
//		fmt.Printf("%-20s +%-6s (%d rolls)\n", stat.Name, stat.DisplayValue(), stat.Rolls)
//	}
func (discs DriveDiscs) SubStatTotals() []StatValue {
	totals := make(map[PropertyID]StatValue)
	var order []PropertyID

	for _, disc := range discs {
		for _, sub := range disc.SubStats {
			if curr, exists := totals[sub.PropertyID]; exists {
				curr.Value += sub.Value
				curr.Rolls += sub.Rolls
				totals[sub.PropertyID] = curr
			} else {
				totals[sub.PropertyID] = sub
				order = append(order, sub.PropertyID)
			}
		}
	}

	result := make([]StatValue, 0, len(order))
	for _, id := range order {
		result = append(result, totals[id])
	}
	return result
}

// CountEffectiveRolls returns the total number of sub-stat upgrade rolls across all [DriveDisc] entries in the collection
// that match any of the provided target property IDs (also known as "effective" or "useful" rolls for build evaluation).
//
// Example:
//
//	// Evaluate build quality on an Attack Agent across all 6 equipped discs:
//	rolls := agent.DriveDiscs.CountEffectiveRolls(fairy.PropCritRate, fairy.PropCritDMG, fairy.PropATKPercent)
//	fmt.Printf("Effective rolls: %d\n", rolls)
func (discs DriveDiscs) CountEffectiveRolls(targetProps ...PropertyID) int {
	total := 0
	targetMap := make(map[PropertyID]bool)
	for _, p := range targetProps {
		targetMap[p] = true
	}

	for _, disc := range discs {
		for _, sub := range disc.SubStats {
			if targetMap[sub.PropertyID] {
				total += sub.Rolls
			}
		}
	}
	return total
}

// BySlot returns a pointer to the [DriveDisc] equipped in the specified partition slot (1–6),
// or nil if no disc is equipped in that slot.
func (discs DriveDiscs) BySlot(slot int) *DriveDisc {
	for i := range discs {
		if discs[i].Slot == slot {
			return &discs[i]
		}
	}
	return nil
}

// SetCounts groups equipped discs by their [Set] and returns the count of pieces equipped for each set.
func (discs DriveDiscs) SetCounts() map[Set]int {
	counts := make(map[Set]int)
	for _, disc := range discs {
		if disc.Set.ID > 0 {
			counts[disc.Set]++
		}
	}
	return counts
}

// Has2Piece reports whether at least 2 pieces of the specified [SetID] are equipped.
//
// Example:
//
//	if agent.DriveDiscs.Has2Piece(fairy.SetSwingJazz) {
//		fmt.Println("Swing Jazz 2-pc is active (+20% Energy Regen)")
//	}
func (discs DriveDiscs) Has2Piece(setID SetID) bool {
	count := 0
	for _, disc := range discs {
		if disc.Set.ID == setID {
			count++
			if count >= 2 {
				return true
			}
		}
	}
	return false
}

// Has4Piece reports whether at least 4 pieces of the specified [SetID] are equipped.
//
// Example:
//
//	if agent.DriveDiscs.Has4Piece(fairy.SetWoodpeckerElectro) {
//		fmt.Println("Woodpecker Electro 4-pc is active")
//	}
func (discs DriveDiscs) Has4Piece(setID SetID) bool {
	count := 0
	for _, disc := range discs {
		if disc.Set.ID == setID {
			count++
			if count >= 4 {
				return true
			}
		}
	}
	return false
}

// CountEffectiveRolls returns the total number of sub-stat upgrade rolls on this specific [DriveDisc]
// that match any of the provided target property IDs.
//
// Example:
//
//	// Count useful rolls for an Attack Agent on a single disc:
//	rolls := disc.CountEffectiveRolls(fairy.PropCritRate, fairy.PropCritDMG, fairy.PropATKPercent)
func (d *DriveDisc) CountEffectiveRolls(targetProps ...PropertyID) int {
	total := 0
	targetMap := make(map[PropertyID]bool)
	for _, p := range targetProps {
		targetMap[p] = true
	}

	for _, sub := range d.SubStats {
		if targetMap[sub.PropertyID] {
			total += sub.Rolls
		}
	}
	return total
}

// DriveDiscSetBonus represents an aggregated Drive Disc set equipped on an [Agent],
// detailing the equipped piece count and active/inactive set threshold bonuses.
type DriveDiscSetBonus struct {
	// Set is the Drive Disc [Set] metadata (ID and localized name).
	Set Set `json:"set"`

	// Count is the total number of pieces from this set currently equipped on the Agent (e.g. 2, 4, 5).
	Count int `json:"count"`

	// Effects contains the 2-piece and 4-piece [SetEffect] thresholds with their activation status.
	Effects []SetEffect `json:"effects"`
}

// SetEffect represents a specific Drive Disc set effect threshold (2-piece or 4-piece bonus).
type SetEffect struct {
	// PieceCount is the required piece count threshold (2 or 4).
	PieceCount int `json:"piece_count"`

	// Description is the localized description text of the set bonus effect.
	Description string `json:"description"`

	// FormattedHTML is the web-ready HTML description with inline CSS colors.
	FormattedHTML string `json:"formatted_html,omitempty"`

	// IsActive indicates whether this set effect is currently active on the Agent (Count >= PieceCount).
	IsActive bool `json:"is_active"`
}

// FormatHTML returns the set effect description formatted as HTML with inline CSS styling.
func (e SetEffect) FormatHTML() string {
	return formatHTML(e.Description)
}

// FormatPlainText returns the set effect description as clean plain text with all Rich Text tags stripped.
func (e SetEffect) FormatPlainText() string {
	return formatPlainText(e.Description)
}

// FormatMarkdown returns the set effect description formatted with Markdown syntax.
func (e SetEffect) FormatMarkdown() string {
	return formatMarkdown(e.Description)
}
