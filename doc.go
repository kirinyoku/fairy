// Package fairy provides a highly modular and extensible Go library for fetching,
// parsing, and enriching Zenless Zone Zero player profiles via the EnkaNetwork API.
//
// The raw response from Enka.Network provides basic IDs for agents, W-Engines, and Drive Discs.
// Fairy takes care of the heavy lifting by replacing raw IDs with full localized names for agents,
// weapons, discs, and skills, building full URLs for splash arts, icons, and avatars, and
// calculating precise final combat stats taking into account base stats, weapon scalings,
// disc substat rolls, and set bonuses.
//
// # Quick Start
//
// The easiest way to get started is by using the global functions. By default, this uses
// the embedded metadata store and English localization:
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
//		fairy.WithDefaultLang(fairy.LangJA), // Default to Japanese
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// See https://github.com/kirinyoku/fairy for more advanced features.
package fairy
