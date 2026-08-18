// Package fairy provides a highly modular and extensible Go library for fetching,
// parsing, and enriching Zenless Zone Zero player profiles via the EnkaNetwork API.
//
// The raw response from Enka.Network provides basic IDs for agents, W-Engines, and Drive Discs.
// Fairy takes care of the heavy lifting by replacing raw IDs with full localized names for agents,
// weapons, discs, and skills, building full URLs for splash arts, icons, and avatars, and
// calculating precise final combat stats taking into account base stats, weapon scalings,
// disc substat rolls, and set bonuses across 13 supported languages.
//
// # Core Actions
//
// Fairy provides four actions for working with player profiles:
//   - [GetProfile] — fetch and enrich a profile in the client's default language.
//   - [GetProfileWithLang] — same, but override the language for a single request.
//   - [GetRawProfile] — fetch the raw upstream API response without enrichment.
//   - [Localize] — map a raw profile into an enriched [Profile] in memory (zero network calls).
//
// Each action is available both as a global convenience function (using a shared default client
// with English localization) and as a method on [Client].
//
// # Quick Start
//
// The easiest way to get started is by using the global functions:
//
//	profile, err := fairy.GetProfile(context.Background(), "1504687050")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Player: %s (Level %d)\n", profile.Nickname, profile.InterknotLevel)
//
// # Custom Client
//
// If you are building a multi-language application or want to configure custom HTTP settings,
// you should create a dedicated client:
//
//	client, err := fairy.NewClient(
//		fairy.WithDefaultLang(fairy.LangJA),
//		fairy.WithEnkaOptions(zzz.Options{
//			UserAgent:  "MyApp/1.0 (contact@example.com)",
//			HTTPClient: &http.Client{Timeout: 10 * time.Second},
//			Retry:      &zzz.RetryOptions{MaxAttempts: 2, Delay: 2 * time.Second},
//			Cache:      myCacheInstance,
//		}),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
package fairy
