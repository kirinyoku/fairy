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
//   - Providing Drive Disc substat aggregation, set bonus evaluation, and build quality scoring tools.
//   - Validating player UID syntax and resolving server regions locally without network calls.
//
// # Key Features
//
//   - Zero-Allocation Metadata Store: Embedded game data loaded once and shared across all queries with lazy-loaded localization.
//   - Full Localization: 13 officially supported languages with on-the-fly in-memory enrichment ([Enrich], [EnrichWithLang], [EnrichAgent], [EnrichAgentWithLang]).
//   - Accurate Combat Math: Emulates the exact ZZZ stat calculations including base attributes, weapon scaling, and set bonuses.
//   - UI-Ready Stats Breakdown: Pre-calculated base, added, and total values formatted for frontends ([UIStats], [FormattedStatBreakdown]).
//   - Rich Text Parsers: Convert game descriptions to clean HTML, Plain Text, or Markdown ([Skill.FormatHTML], [Skill.FormatPlainText], [Skill.FormatMarkdown], [SetEffect.FormatHTML]).
//   - UID Validation & Region Discovery: Instant offline syntax verification ([IsValidUID]) and server region lookup ([RegionFromUID]).
//   - Cache Governance: Built-in upstream API cache TTL inspection via [Profile.CacheTTL].
//   - Production-Grade Client: Supports HTTP timeouts, automatic exponential retries, Redis/in-memory caching, and custom User-Agents via [zzz.Options].
//
// # Architecture & Data Flow
//
// The library separates network fetching, metadata mapping, and combat calculation:
//
//	Full Profile Flow:
//	  [EnkaNetwork API]
//	         │
//	         ▼ (HTTP Request via internal API client)
//	  [zzz.Profile (Raw Upstream Model)]
//	         │
//	         ▼ (Enrich / mapper using embedded MetadataStore)
//	  [fairy.Profile (Enriched Domain Model)]
//	         ├── Account Info (UID, Nickname, InterknotLevel, Region, Title, Avatar, Badges)
//	         └── Showcase Agents (max 6)
//	               ├── Agent Meta (Attribute, Specialty, Rarity, Skin, SplashArt)
//	               ├── Skills & Groups (Basic, Dodge, Special, Chain, Assist, Passives)
//	               ├── Mindscape Cinema (Ranks 1–6 with unlocked status)
//	               ├── Potential Vision (Active nodes & descriptions)
//	               ├── Equipped W-Engine (MainStat, SecondaryStat, Modification)
//	               ├── Drive Discs (Slots 1–6 with roll counts, active 2-piece & 4-piece Set Bonuses)
//	               └── Combat Stats Pipeline
//	                     ├── BaseStats (Agent + W-Engine Base ATK)
//	                     ├── Stats (Final calculated combat stats)
//	                     └── UIStats (Pre-formatted Base + Added = Total breakdowns)
//
//	Standalone Agent Flow:
//	  [zzz.AvatarData (Raw Upstream Agent)]
//	         │
//	         ▼ (EnrichAgent / mapper using embedded MetadataStore)
//	  [fairy.Agent (Enriched Domain Model)]
//
// # Core Operations
//
// Fairy provides core operations for working with player profiles and individual agents:
//   - [GetProfile]: Fetch and enrich a player profile using the client's default language.
//   - [GetProfileWithLang]: Fetch and enrich a player profile with a specific language override for the request.
//   - [GetRawProfile]: Fetch the raw upstream API response ([zzz.Profile]) without enrichment.
//   - [Enrich]: Transform a raw [zzz.Profile] into an enriched [Profile] using the default language (zero additional network requests).
//   - [EnrichWithLang]: Transform a raw [zzz.Profile] into an enriched [Profile] in the specified language (zero additional network requests).
//   - [EnrichAgent]: Transform a raw [zzz.AvatarData] into an enriched [Agent] using the default language (zero additional network requests).
//   - [EnrichAgentWithLang]: Transform a raw [zzz.AvatarData] into an enriched [Agent] in the specified language (zero additional network requests).
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
//	fmt.Printf("Player: %s (Inter-Knot Lv.%d, Server: %s, Cache TTL: %v)\n",
//		profile.Nickname, profile.InterknotLevel, profile.Region, profile.CacheTTL())
//
//	for _, agent := range profile.Agents {
//		fmt.Printf("• %-16s Lv.%-2d [%s / %s]\n",
//			agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName)
//	}
//
// # UID Validation & Server Region Detection
//
// Validate player UIDs and determine their game server region before making network calls:
//
//	const uid = "1504687050"
//	if !fairy.IsValidUID(uid) {
//		log.Fatalf("Invalid player UID format: %s", uid)
//	}
//
//	if region, ok := fairy.RegionFromUID(uid); ok {
//		fmt.Printf("Server region: %s\n", region) // "Europe"
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
// # In-Memory Multi-Language Enrichment
//
// If you need to present the same player profile in multiple languages, fetch the raw profile once
// and enrich it in memory with [Enrich] or [EnrichWithLang] without repeating network requests.
//
//	raw, err := fairy.GetRawProfile(ctx, "1504687050")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// In-memory mapping — instant, no extra network overhead
//	enProfile, _ := fairy.Enrich(raw) // Default English
//	jaProfile, _ := fairy.EnrichWithLang(raw, fairy.LangJA)
//	ruProfile, _ := fairy.EnrichWithLang(raw, fairy.LangRU)
//
// Individual agents can also be enriched directly without a full profile:
//
//	agentEN, _ := fairy.EnrichAgent(rawAvatar)
//	agentJA, _ := fairy.EnrichAgentWithLang(rawAvatar, fairy.LangJA)
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
// # Drive Disc & Set Bonus Analysis
//
// Fairy provides helper methods on [DriveDiscs] to inspect and evaluate equipped Drive Discs:
//
//	// 1. Check active Drive Disc set bonuses (2-piece and 4-piece thresholds)
//	if agent.DriveDiscs.Has4Piece(fairy.SetPolarMetal) {
//		fmt.Println("Polar Metal 4-piece set bonus active!")
//	}
//
//	// 2. Group and sum substats across all 6 disc slots
//	totals := agent.DriveDiscs.SubStatTotals()
//	for _, sub := range totals {
//		fmt.Printf("%-20s +%-6s (%d rolls)\n", sub.Name, sub.DisplayValue(), sub.Rolls)
//	}
//
//	// 3. Count "effective" (useful) substat rolls for a build
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
//		case errors.Is(err, fairy.ErrInvalidUID):
//			// Player UID format is invalid (must be 10 digits starting with 10, 13, 15, or 17)
//		case errors.Is(err, fairy.ErrProfileNotFound):
//			// Profile with specified UID does not exist on game servers
//		case errors.Is(err, fairy.ErrRateLimit):
//			// Rate limit reached (HTTP 429) — back off and retry
//		case errors.Is(err, fairy.ErrMaintenance):
//			// Upstream API or game servers are under maintenance
//		case errors.Is(err, fairy.ErrNetwork):
//			// Network timeout or connectivity issue
//		case errors.Is(err, fairy.ErrEnrichment):
//			// In-memory transformation, metadata mapping, or stat calculation failed
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
// Use [AllRegions] to retrieve the full list, or [RegionFromUID] to resolve a region from a player's UID.
//
// # Combat Attributes, Specialties & Drive Disc Sets
//
// Combat attributes, character classes, and equipment sets are defined by strongly-typed constants:
//   - [Attribute]: [AttributePhysical], [AttributeFire], [AttributeIce], [AttributeElectric], [AttributeEther], and variants ([AttributeHonedEdge], [AttributeFrost], [AttributeAuricInk], [AttributeWind], [AttributeLumiflux]). Use [AllAttributes] to retrieve all.
//   - [Specialty]: [SpecialtyAttack], [SpecialtyStun], [SpecialtyAnomaly], [SpecialtySupport], [SpecialtyDefense], [SpecialtyRupture], and [SpecialtyArmorer]. Use [AllSpecialties] to retrieve all.
//   - [Rarity]: [RarityS], [RarityA], and [RarityB]. Use [AllRarities] to retrieve all.
//   - [SetID]: Drive Disc set identifiers (e.g. [SetPolarMetal], [SetWoodpeckerElectro], [SetBranchBladeSong], etc.). Use [AllSetIDs] to retrieve all.
package fairy
