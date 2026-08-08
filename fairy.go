package fairy

import (
	"context"
	"sync"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
)

// Language represents the localization language for game strings.
type Language string

// Supported language constants matching the in-game localizations.
// These determine which translation strings are pulled from the metadata store.
const (
	LangEN   Language = "en"    // English
	LangRU   Language = "ru"    // Russian
	LangDE   Language = "de"    // German
	LangES   Language = "es"    // Spanish
	LangFR   Language = "fr"    // French
	LangID   Language = "id"    // Indonesian
	LangJA   Language = "ja"    // Japanese
	LangKO   Language = "ko"    // Korean
	LangPT   Language = "pt"    // Portuguese
	LangTH   Language = "th"    // Thai
	LangVI   Language = "vi"    // Vietnamese
	LangZHCN Language = "zh-cn" // Chinese (Simplified)
	LangZHTW Language = "zh-tw" // Chinese (Traditional)
)

var (
	defaultClient     *Client
	defaultClientOnce sync.Once
)

// getDefaultClient lazily initializes and returns the shared global Client.
func getDefaultClient() *Client {
	defaultClientOnce.Do(func() {
		// NewClient handles the default Language (EN) and MetadataStore (Embedded)
		client, err := NewClient()
		if err != nil {
			panic("fairy: failed to initialize default client: " + err.Error())
		}
		defaultClient = client
	})
	return defaultClient
}

// GetProfile fetches the game profile by UID and localizes it into the default language
// using a shared, global client.
func GetProfile(ctx context.Context, uid string) (*Profile, error) {
	return getDefaultClient().GetProfile(ctx, uid)
}

// GetProfileWithLang fetches the game profile by UID and localizes it into the specified language
// using a shared, global client.
func GetProfileWithLang(ctx context.Context, uid string, lang Language) (*Profile, error) {
	return getDefaultClient().GetProfileWithLang(ctx, uid, lang)
}

// GetRawProfile fetches the raw zzz.Profile by UID without applying any localization
// using a shared, global client.
func GetRawProfile(ctx context.Context, uid string) (*zzz.Profile, error) {
	return getDefaultClient().GetRawProfile(ctx, uid)
}

// Localize maps a raw zzz.Profile into an enriched Profile using the specified language.
// It uses the global default client's metadata store for the conversion.
func Localize(raw *zzz.Profile, lang Language) (*Profile, error) {
	return getDefaultClient().Localize(raw, lang)
}
