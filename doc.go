// Package fairy provides a high-performance, modular Go library for fetching,
// parsing, enriching, and calculating combat stats for Zenless Zone Zero (ZZZ)
// player profiles via the EnkaNetwork API.
//
// The raw response from Enka.Network contains raw Agent, W-Engine, and Drive Disc IDs.
// Fairy enriches these raw responses by:
//   - Translating IDs into full localized names for Agents, W-Engines, Drive Discs, Skills, and Mindscapes across 13 languages.
//   - Constructing complete Enka CDN URLs and inline base64 SVG data URIs for splash arts, avatars, namecards, badges, and stat icons.
//   - Calculating exact final combat stats according to the in-game formula (accounting for Agent levels, Promotions, Core Skill Enhancements, W-Engine growth curves, Drive Disc substats, and set bonuses).
//   - Evaluating Unity Rich Text formatting, button icon tags, and dynamic level-scaling formulas.
//   - Providing Drive Disc substat aggregation and build quality evaluation tools.
//
// # Key Features
//
//   - Zero-Allocation Metadata Store: Embedded game data loaded once and shared across all queries.
//   - Full Localization: 13 officially supported languages with on-the-fly in-memory localization ([Localize]).
//   - Accurate Combat Math: Emulates the exact ZZZ stat calculations.
//   - UI-Ready Stats Breakdown: Pre-calculated base, added, and total values formatted for frontends ([UIStats], [FormattedStatBreakdown]).
//   - Rich Text Parsers: Convert game descriptions to clean HTML, Plain Text, or Markdown ([Skill.FormatHTML], [Skill.FormatPlainText], [Skill.FormatMarkdown]).
//   - Production-Grade Client: Supports HTTP timeouts, automatic exponential retries, Redis/in-memory caching, and custom User-Agents via [zzz.Options].
//
// # Architecture & Data Flow
//
// The library separates network fetching, metadata mapping, and combat calculation:
//
//	[EnkaNetwork API]
//	       │
//	       ▼ (HTTP Request via internal API client)
//	[zzz.Profile (Raw Upstream Model)]
//	       │
//	       ▼ (Localize / mapper using embedded MetadataStore)
//	[fairy.Profile (Enriched Domain Model)]
//	       ├── Account Info (UID, Nickname, InterknotLevel, Region, Title, Avatar, Badges)
//	       └── Showcase Agents (max 6)
//	             ├── Agent Meta (Attribute, Specialty, Rarity, Skin, SplashArt)
//	             ├── Skills & Groups (Basic, Dodge, Special, Chain, Assist, Passives)
//	             ├── Mindscape Cinema (Ranks 1–6 with unlocked status)
//	             ├── Potential Vision (Active nodes & descriptions)
//	             ├── Equipped W-Engine (MainStat, SecondaryStat, Modification)
//	             ├── Equipped Drive Discs (Slots 1–6, MainStat, Substats with Roll counts)
//	             ├── Active Drive Disc Set Bonuses (2-piece & 4-piece thresholds)
//	             └── Combat Stats Pipeline
//	                   ├── BaseStats (Agent + W-Engine Base ATK)
//	                   ├── Stats (Final calculated combat stats)
//	                   └── UIStats (Pre-formatted Base + Added = Total breakdowns)
//
// # Core Operations
//
// Fairy provides four core operations for working with player profiles:
//   - [GetProfile]: Fetch and enrich a player profile using the client's default language.
//   - [GetProfileWithLang]: Fetch and enrich a player profile with a specific language override for the request.
//   - [GetRawProfile]: Fetch the raw upstream API response ([zzz.Profile]) without enrichment.
//   - [Localize]: Map a raw [zzz.Profile] into an enriched [Profile] in memory (zero additional network requests).
//
// Each operation is available as a global top-level function (using a shared thread-safe default client)
// and as a method on [Client].
//
// # Quick Start
//
// For standard usage, use the global top-level functions:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	profile, err := fairy.GetProfile(ctx, "1504687050")
//	if err != nil {
//		log.Fatalf("Failed to fetch profile: %v", err)
//	}
//
//	fmt.Printf("Player: %s (Inter-Knot Lv.%d, Server: %s)\n",
//		profile.Nickname, profile.InterknotLevel, profile.Region)
//
//	for _, agent := range profile.Agents {
//		fmt.Printf("• %-16s Lv.%-2d [%s / %s]\n",
//			agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName)
//	}
//
// # Custom Client Configuration
//
// For production services, configure a custom [Client] with timeouts, retry policies, custom headers,
// or persistent caching:
//
//	client, err := fairy.NewClient(
//		fairy.WithDefaultLang(fairy.LangJA),
//		fairy.WithEnkaOptions(zzz.Options{
//			UserAgent:  "MyZZZApp/1.0 (contact@example.com)",
//			HTTPClient: &http.Client{Timeout: 10 * time.Second},
//			Retry: &zzz.RetryOptions{
//				MaxAttempts: 3,
//				Delay:       1 * time.Second,
//			},
//			Cache: myCacheInstance, // Implements zzz.Cache interface
//		}),
//	)
//	if err != nil {
//		log.Fatalf("Failed to initialize client: %v", err)
//	}
//
//	profile, err := client.GetProfile(context.Background(), "1504687050")
//
// # In-Memory Multi-Language Localization
//
// If you need to present the same player profile in multiple languages, fetch the raw profile once
// and localize it in memory with [Localize] without repeating network requests.
//
//	raw, err := fairy.GetRawProfile(ctx, "1504687050")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// In-memory mapping — instant, no extra network overhead
//	deProfile, _ := fairy.Localize(raw, fairy.LangDE)
//	jaProfile, _ := fairy.Localize(raw, fairy.LangJA)
//	ruProfile, _ := fairy.Localize(raw, fairy.LangRU)
//
// # Combat Stat Breakdown & UI Display
//
// Each showcased [Agent] includes a pre-computed [UIStats] panel containing formatted breakdowns
// of base stats vs. added stats (from W-Engines, Drive Discs, and Set Bonuses):
//
//	for _, stat := range agent.UIStats.List() {
//		fmt.Printf("%-22s %8s (Base: %s + Added: %s)\n",
//			stat.Name, stat.Total, stat.Base, stat.Added)
//	}
//
// # Drive Disc & Substat Roll Analysis
//
// Fairy provides helper methods on [DriveDiscs] to inspect and evaluate equipped Drive Discs:
//
//	// 1. Group and sum substats across all 6 disc slots
//	totals := agent.DriveDiscs.SubStatTotals()
//	for _, sub := range totals {
//		fmt.Printf("%-20s +%-6s (%d rolls)\n", sub.Name, sub.DisplayValue(), sub.Rolls)
//	}
//
//	// 2. Count "effective" (useful) substat rolls for a build
//	usefulRolls := agent.DriveDiscs.CountEffectiveRolls(
//		fairy.PropCritRate,
//		fairy.PropCritDMG,
//		fairy.PropATKPercent,
//	)
//	fmt.Printf("Useful Substat Rolls: %d\n", usefulRolls)
//
// # Error Handling
//
// All API errors map to strongly-typed sentinel errors that can be inspected with [errors.Is]:
//
//	profile, err := fairy.GetProfile(ctx, uid)
//	if err != nil {
//		switch {
//		case errors.Is(err, fairy.ErrProfileNotFound):
//			// Profile with specified UID does not exist
//		case errors.Is(err, fairy.ErrRateLimit):
//			// Rate limit reached (HTTP 429) — back off and retry
//		case errors.Is(err, fairy.ErrMaintenance):
//			// Upstream API or game servers are under maintenance
//		case errors.Is(err, fairy.ErrNetwork):
//			// Network timeout or connectivity issue
//		default:
//			// Other unexpected errors
//		}
//	}
//
// # Supported Languages
//
// Supported game languages are defined by the [Language] type constants:
// [LangEN] (English), [LangRU] (Russian), [LangDE] (German), [LangES] (Spanish),
// [LangFR] (French), [LangID] (Indonesian), [LangJA] (Japanese), [LangKO] (Korean),
// [LangPT] (Portuguese), [LangTH] (Thai), [LangVI] (Vietnamese),
// [LangZHCN] (Chinese Simplified), and [LangZHTW] (Chinese Traditional).
// Use [AllLanguages] to retrieve the full list programmatically.
//
// # Server Regions
//
// Player server regions are identified by the [Region] type constants:
// [RegionEU] (Europe), [RegionNA] (America), [RegionAsia] (Asia), and [RegionTWHKMO] (TW/HK/MO).
// Use [AllRegions] to retrieve the full list.
package fairy
