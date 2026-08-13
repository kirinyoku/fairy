package fairy

// Region represents the game server region (e.g., Europe, America).
type Region string

const (
	RegionEU     Region = "Europe"   // European server region.
	RegionNA     Region = "America"  // North American server region.
	RegionAsia   Region = "Asia"     // Asian server region.
	RegionTWHKMO Region = "TW/HK/MO" // Taiwan/Hong Kong/Macau server region.
)

// Profile represents the enriched user profile data.
// It contains player-level metadata and the showcased agents.
type Profile struct {
	UID            string    `json:"uid"`             // The unique identifier of the player (typically a 9-digit string).
	Nickname       string    `json:"nickname"`        // The player's chosen nickname. Can be empty if the API returned no name.
	InterknotLevel int       `json:"interknot_level"` // The player's overall account level.
	Region         Region    `json:"region"`          // The server region the player belongs to.
	Title          *Title    `json:"title"`           // The active title displayed on the profile. May be nil if none is equipped.
	Avatar         *Avatar   `json:"avatar"`          // The active avatar (profile picture) displayed. May be nil.
	Namecard       *Namecard `json:"namecard"`        // The active background namecard. May be nil.
	Badges         []Badge   `json:"badges"`          // The showcase badges selected by the player. Can be empty.
	Agents         []Agent   `json:"agents"`          // The list of agents showcased on the profile (max 6). Can be empty.
}

// Avatar represents a selectable proxy avatar (profile picture).
type Avatar struct {
	ID  int    `json:"id"`  // The internal ID of the avatar.
	URL string `json:"url"` // The URL to the avatar's image asset.
}

// Namecard represents a profile background image.
type Namecard struct {
	ID  int    `json:"id"`  // The internal ID of the namecard.
	URL string `json:"url"` // The URL to the namecard's background asset.
}

// Title represents an achievement or status title chosen by the player.
// Titles often have gradients represented by two hex colors. To properly display
// this gradient, use the PrimaryColorHex() and SecondaryColorHex() helpers.
type Title struct {
	ID             int    `json:"id"`              // The internal ID of the title.
	Text           string `json:"text"`            // The localized text of the title.
	PrimaryColor   string `json:"primary_color"`   // The primary color hex (without #).
	SecondaryColor string `json:"secondary_color"` // The secondary color hex (without #).
}

// PrimaryColorHex returns the primary gradient color formatted as a standard hex string (#RRGGBB).
func (t *Title) PrimaryColorHex() string {
	if t == nil || t.PrimaryColor == "" {
		return ""
	}
	return "#" + t.PrimaryColor
}

// SecondaryColorHex returns the secondary gradient color formatted as a standard hex string (#RRGGBB).
func (t *Title) SecondaryColorHex() string {
	if t == nil || t.SecondaryColor == "" {
		return ""
	}
	return "#" + t.SecondaryColor
}

// Badge represents a collectible medal or badge displayed on the profile.
type Badge struct {
	ID      int    `json:"id"`       // The internal ID of the badge.
	Title   string `json:"title"`    // The localized name/title of the badge.
	Value   int    `json:"value"`    // The progression value associated with the badge.
	IconURL string `json:"icon_url"` // The URL to the badge's visual icon.
}
