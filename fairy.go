package fairy

import (
	"context"
	"sync"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
)

// Language represents a supported localization language for in-game text.
// It determines the translation strings retrieved from the embedded metadata store
// for Agent names, W-Engine titles, Drive Disc sets, Skill descriptions,
// Mindscape Cinema nodes, and combat stat labels.
type Language string

// Supported in-game localization languages matching the official Zenless Zone Zero game clients.
const (
	// LangEN represents English localization ("en").
	LangEN Language = "en"
	// LangRU represents Russian localization ("ru").
	LangRU Language = "ru"
	// LangDE represents German localization ("de").
	LangDE Language = "de"
	// LangES represents Spanish localization ("es").
	LangES Language = "es"
	// LangFR represents French localization ("fr").
	LangFR Language = "fr"
	// LangID represents Indonesian localization ("id").
	LangID Language = "id"
	// LangJA represents Japanese localization ("ja").
	LangJA Language = "ja"
	// LangKO represents Korean localization ("ko").
	LangKO Language = "ko"
	// LangPT represents Portuguese localization ("pt").
	LangPT Language = "pt"
	// LangTH represents Thai localization ("th").
	LangTH Language = "th"
	// LangVI represents Vietnamese localization ("vi").
	LangVI Language = "vi"
	// LangZHCN represents Simplified Chinese localization ("zh-cn").
	LangZHCN Language = "zh-cn"
	// LangZHTW represents Traditional Chinese localization ("zh-tw").
	LangZHTW Language = "zh-tw"
)

var allLanguages = [...]Language{
	LangEN,
	LangRU,
	LangDE,
	LangES,
	LangFR,
	LangID,
	LangJA,
	LangKO,
	LangPT,
	LangTH,
	LangVI,
	LangZHCN,
	LangZHTW,
}

// AllLanguages returns a newly allocated slice containing all 13 supported [Language] localizations.
// The returned slice is a defensive copy and can be safely mutated by the caller.
func AllLanguages() []Language {
	langs := make([]Language, len(allLanguages))
	copy(langs, allLanguages[:])
	return langs
}

// IsValid reports whether the language is one of the 13 officially supported [Language] constants.
func (l Language) IsValid() bool {
	for _, lang := range allLanguages {
		if l == lang {
			return true
		}
	}
	return false
}

var (
	defaultClient     *Client
	defaultClientOnce sync.Once
	defaultClientErr  error
)

// getDefaultClient returns the shared thread-safe singleton Client instance.
func getDefaultClient() (*Client, error) {
	defaultClientOnce.Do(func() {
		client, err := NewClient()
		if err != nil {
			defaultClientErr = err
			return
		}
		defaultClient = client
	})
	return defaultClient, defaultClientErr
}

// GetProfile fetches a player profile by UID via the EnkaNetwork API
// and enriches it using the shared default client in English ([LangEN]).
//
// The returned [Profile] contains:
//   - Player account details (UID, Nickname, Inter-Knot Level, Region, Title, Avatar, Badges).
//   - Showcase [Agent] list (up to 6 agents), where each agent contains:
//   - Metadata (localized name, Attribute, Specialty, Rarity, Skin, CDN asset URLs).
//   - Progression (categorized [SkillGroup] entries, [MindscapeNode] unlock states, [PotentialVision]).
//   - Equipped [WEngine] (with level-scaled stats and refinement modification effects).
//   - Equipped [DriveDisc] entries (slots 1–6 with roll counts) and active [DriveDiscSetBonus] thresholds.
//   - Pre-calculated combat [Stats] and frontend-ready [UIStats] breakdowns.
//
// The provided [context.Context] controls the HTTP request lifecycle, cancellation, and timeout.
// Returns sentinel errors such as [ErrInvalidUID], [ErrProfileNotFound], [ErrRateLimit], [ErrMaintenance], [ErrNetwork], or [ErrEnrichment].
func GetProfile(ctx context.Context, uid string) (*Profile, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}
	return client.GetProfile(ctx, uid)
}

// GetProfileWithLang fetches a player profile by UID via the EnkaNetwork API
// and enriches it using the specified [Language] localization.
//
// This is identical to [GetProfile], but overrides the default language for the single request
// without modifying the shared default client.
//
// The provided [context.Context] controls the HTTP request lifecycle, cancellation, and timeout.
// Returns sentinel errors such as [ErrInvalidUID], [ErrProfileNotFound], [ErrRateLimit], [ErrMaintenance], [ErrNetwork], or [ErrEnrichment].
func GetProfileWithLang(ctx context.Context, uid string, lang Language) (*Profile, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}
	return client.GetProfileWithLang(ctx, uid, lang)
}

// GetRawProfile fetches the un-enriched [zzz.Profile] directly from the EnkaNetwork API
// using the shared default client.
//
// Use this function when you only need the raw numeric IDs from the upstream API without metadata enrichment,
// or when you want to fetch the upstream payload once and enrich it into multiple languages via [EnrichWithLang].
//
// The provided [context.Context] controls the HTTP request lifecycle, cancellation, and timeout.
// Returns sentinel errors such as [ErrInvalidUID], [ErrProfileNotFound], [ErrRateLimit], [ErrMaintenance], or [ErrNetwork].
func GetRawProfile(ctx context.Context, uid string) (*zzz.Profile, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}
	return client.GetRawProfile(ctx, uid)
}

// Enrich transforms a raw upstream [zzz.Profile] into an enriched [Profile]
// using the default [Language] (English).
//
// This function operates completely in-memory using the embedded metadata store and makes ZERO network requests.
// It resolves all progression data, computes scaled combat stats, parses Unity Rich Text into HTML,
// and assembles the full domain model.
// Returns [ErrEnrichment] if the raw profile payload is nil or corrupt.
func Enrich(raw *zzz.Profile) (*Profile, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}
	return client.Enrich(raw)
}

// EnrichWithLang transforms a raw upstream [zzz.Profile] into an enriched [Profile]
// in the requested [Language].
//
// This function operates completely in-memory using the embedded metadata store and makes ZERO network requests.
// It is ideal for multi-language applications that fetch a player's raw profile once via [GetRawProfile]
// and render it dynamically across different languages.
// Returns [ErrEnrichment] if the raw profile payload is nil or corrupt.
func EnrichWithLang(raw *zzz.Profile, lang Language) (*Profile, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}
	return client.EnrichWithLang(raw, lang)
}
