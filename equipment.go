package fairy

// WEngine represents an enriched W-Engine equipped by an agent.
type WEngine struct {
	ID                 int       `json:"id"`                  // The internal ID of the W-Engine.
	UID                string    `json:"uid"`                 // The unique instance ID of this specific W-Engine.
	Name               string    `json:"name"`                // The localized name of the W-Engine.
	Level              int       `json:"level"`               // The current level of the W-Engine (1-60).
	Phase              int       `json:"phase"`               // The star level from the API (ascension phase).
	Modification       int       `json:"modification"`        // The refinement/upgrade level of the passive (1-5).
	Rarity             Rarity    `json:"rarity"`              // The rarity tier (S, A, or B).
	Specialty          Specialty `json:"specialty"`           // The intended role specialty for this W-Engine.
	SpecialtyName      string    `json:"specialty_name"`      // The localized name of the specialty.
	IconURL            string    `json:"icon_url"`            // The URL to the W-Engine's visual icon.
	MainStat           StatValue `json:"main_stat"`           // The primary stat provided by the W-Engine.
	SecondaryStat      StatValue `json:"secondary_stat"`      // The secondary stat provided by the W-Engine.
	PassiveDescription string    `json:"passive_description"` // The localized description of the passive skill. Note: Changes based on Modification phase.
}

// Set represents a specific Drive Disc equipment set.
// A Set grants bonus effects when an agent equips 2 or 4 pieces of the same set.
type Set struct {
	ID   int    `json:"id"`   // The internal ID of the set.
	Name string `json:"name"` // The localized name of the set (e.g., "Woodpecker Electro").
}

// DriveDisc represents an enriched Drive Disc (artifact/equipment).
type DriveDisc struct {
	ID       int         `json:"id"`        // The internal ID of the specific disc variation.
	UID      string      `json:"uid"`       // The unique instance ID of this specific Drive Disc.
	Set      Set         `json:"set"`       // The equipment set this disc belongs to.
	Slot     int         `json:"slot"`      // The equip slot number (1 to 6). Slots 1-3 have fixed main stats, while 4-6 are randomized.
	Level    int         `json:"level"`     // The upgrade level of the disc (0-15).
	Rarity   Rarity      `json:"rarity"`    // The rarity tier (S, A, B).
	IconPath string      `json:"icon_path"` // The URL to the disc's icon.
	MainStat StatValue   `json:"main_stat"` // The primary stat provided by this disc.
	SubStats []StatValue `json:"sub_stats"` // The randomly rolled sub-stats (up to 4).
}

// CountEffectiveRolls returns the total number of sub-stat rolls on this specific Drive Disc
// that match any of the provided target property IDs.
// Example usage for an Anomaly agent's disc:
//
//	usefulRolls := disc.CountEffectiveRolls(fairy.PropAnomalyProficiency, fairy.PropATKPercent)
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

// DriveDiscSetBonus represents an active set bonus from equipped Drive Discs.
type DriveDiscSetBonus struct {
	Set         Set    `json:"set"`         // The set granting the bonus.
	PieceCount  int    `json:"piece_count"` // The number of pieces equipped from this set (typically 2 or 4).
	Description string `json:"description"` // The localized HTML description of the set bonus from game data.
}
