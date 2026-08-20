package fairy

// Region represents the game server region hosting a player's account.
type Region string

// Supported server region constants.
const (
	// RegionEU represents the European game server region ("Europe").
	RegionEU Region = "Europe"
	// RegionNA represents the North American game server region ("America").
	RegionNA Region = "America"
	// RegionAsia represents the Asian game server region ("Asia").
	RegionAsia Region = "Asia"
	// RegionTWHKMO represents the Taiwan / Hong Kong / Macau game server region ("TW/HK/MO").
	RegionTWHKMO Region = "TW/HK/MO"
)

var allRegions = [...]Region{
	RegionEU,
	RegionNA,
	RegionAsia,
	RegionTWHKMO,
}

// AllRegions returns a newly allocated slice containing all supported [Region] server constants.
// The returned slice is a defensive copy and can be safely mutated by the caller.
func AllRegions() []Region {
	regions := make([]Region, len(allRegions))
	copy(regions, allRegions[:])
	return regions
}

// IsValid reports whether the region is one of the recognized [Region] server constants.
func (r Region) IsValid() bool {
	for _, reg := range allRegions {
		if r == reg {
			return true
		}
	}
	return false
}

// Profile represents an enriched Zenless Zone Zero player profile.
// It aggregates account-level metadata and the player's showcased [Agent] lineup.
type Profile struct {
	// UID is the unique in-game identifier of the player (typically a 9- or 10-digit string).
	UID string `json:"uid"`

	// Nickname is the display name chosen by the player. Can be empty if omitted by upstream API.
	Nickname string `json:"nickname"`

	// InterknotLevel is the player's overall account progression level (Inter-Knot Level).
	InterknotLevel int `json:"interknot_level"`

	// Region is the game server region hosting the player's account.
	Region Region `json:"region"`

	// Title is the active achievement title equipped on the profile. May be nil if no title is equipped.
	Title *Title `json:"title"`

	// Avatar is the active profile picture equipped by the player. May be nil if none is set.
	Avatar *Avatar `json:"avatar"`

	// Namecard is the active profile background image (calling card). May be nil if none is set.
	Namecard *Namecard `json:"namecard"`

	// Badges is the list of collectible showcase medals displayed on the profile. May be empty.
	Badges []Badge `json:"badges"`

	// Agents is the list of up to 6 showcased [Agent] entries configured by the player in-game.
	// May be empty if the player's in-game showcase is empty or has details hidden.
	Agents []Agent `json:"agents"`
}

// Avatar represents a player's equipped profile avatar (profile picture).
type Avatar struct {
	// ID is the internal numeric identifier of the avatar asset.
	ID int `json:"id"`

	// URL is the absolute HTTPS URL pointing to the avatar's image asset on the EnkaNetwork CDN.
	URL string `json:"url"`
}

// Namecard represents a player's equipped background calling card image.
type Namecard struct {
	// ID is the internal numeric identifier of the namecard asset.
	ID int `json:"id"`

	// URL is the absolute HTTPS URL pointing to the namecard's background asset on the EnkaNetwork CDN.
	URL string `json:"url"`
}

// Title represents an achievement or status title displayed on a player's profile.
//
// Many titles in Zenless Zone Zero feature a two-color linear gradient.
// Use [Title.PrimaryColorHex] and [Title.SecondaryColorHex] to retrieve CSS-ready hex color strings.
type Title struct {
	// ID is the internal numeric identifier of the title.
	ID int `json:"id"`

	// Text is the fully localized title text.
	Text string `json:"text"`

	// PrimaryColor is the 6-character hex color code for the starting gradient color (without #).
	PrimaryColor string `json:"primary_color"`

	// SecondaryColor is the 6-character hex color code for the ending gradient color (without #).
	SecondaryColor string `json:"secondary_color"`
}

// PrimaryColorHex returns the primary gradient color formatted as a standard CSS hex string (e.g. "#F7BA3F").
// Returns an empty string if t is nil or if PrimaryColor is not set.
func (t *Title) PrimaryColorHex() string {
	if t == nil || t.PrimaryColor == "" {
		return ""
	}
	return "#" + t.PrimaryColor
}

// SecondaryColorHex returns the secondary gradient color formatted as a standard CSS hex string (e.g. "#E74C3C").
// Returns an empty string if t is nil or if SecondaryColor is not set.
func (t *Title) SecondaryColorHex() string {
	if t == nil || t.SecondaryColor == "" {
		return ""
	}
	return "#" + t.SecondaryColor
}

// Badge represents a collectible achievement medal displayed in a player's profile showcase.
type Badge struct {
	// ID is the internal numeric identifier of the badge icon.
	ID int `json:"id"`

	// Title is the localized title/name of the badge achievement.
	Title string `json:"title"`

	// Value is the progression value or score associated with the badge.
	Value int `json:"value"`

	// IconURL is the absolute HTTPS URL pointing to the badge's visual icon on the EnkaNetwork CDN.
	IconURL string `json:"icon_url"`
}
