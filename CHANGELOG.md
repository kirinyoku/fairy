# Changelog

## [1.1.0] - 2026-08-29

### Added
- Added `TTL` field and `CacheTTL()` helper method to `Profile`.
- Added `IsValidUID(uid)` predicate and `RegionFromUID(uid)` server resolution helper.
- Added `ErrInvalidUID` and `ErrEnrichment` sentinel errors.

### Changed
- Added pre-network UID format validation to `GetProfile`, `GetProfileWithLang`, and `GetRawProfile`.
- Replaced monolithic `locs.json` with per-language gzip-compressed assets (`internal/assets/data/locs/*.json.gz`).
- Implemented lazy loading for embedded localization dictionaries with `sync.RWMutex`.
- Updated `internal/tools/extractor` to read and write per-language compressed files.

### Performance
- Reduced binary size by ~62%.
- Reduced cold start store initialization time from ~150 ms to <10 ms.
- Reduced resident heap memory consumption for single-language usage by ~90% (~1.5 MB vs ~25 MB).

### Tests
- Added unit and concurrency tests for lazy loading in `internal/store`.
- Added unit tests for UID validation, server region lookup, and `Profile.CacheTTL()`.
- Added performance benchmarks for `store.Default()` and `Localize()`.

## [1.0.0] - 2026-08-20

### Added
- Added `DriveDiscs` struct encapsulating `Slots` (discs 1–6) and `SetBonuses` (active 2-pc and 4-pc set bonuses).
- Added typed `SetID` constants for all 30 official Zenless Zone Zero Drive Disc sets with `AllSetIDs()` and `IsValid()` helpers.
- Added `Has2Piece(setID)`, `Has4Piece(setID)`, `BySlot(slot)`, and `SetCounts()` methods to `DriveDiscs`.
- Added `EnrichWithLang` method and package function for explicit in-memory multi-language enrichment.

### Changed
- Promoted library to stable `v1.0.0` release with finalized public API.
- Renamed in-memory mapping from `Localize` to `Enrich` and `EnrichWithLang`.
- Moved `ActiveSetBonuses` from root `Agent` into `agent.DriveDiscs.SetBonuses`.
- Encapsulated `store` package into `internal/store`.
- Expanded GoDoc documentation across all public packages, types, and functions.
- Updated `example_test.go` with 8 practical recipes for `pkg.go.dev`.

### Tests
- Added automated coverage test `TestAllSetIDs_StoreCoverage` verifying all Drive Disc set constants against embedded game data.
- Added unit test coverage for `DriveDiscs` methods, `SetID` validation, and `Enrich` / `EnrichWithLang`.

## [0.11.2] - 2026-08-18

### Added
- Added runnable `Example*` functions for `pkg.go.dev`.

### Changed
- Expanded `doc.go` with Core Actions summary and Custom Client example.

## [0.11.1] - 2026-08-17

### Fixed
- Restored canonical in-game stat order in `UIStats.List()` and `Stats.List()`.

## [0.11.0] - 2026-08-17

### Added
- Added `List()` method to `UIStats` returning combat stat breakdowns in canonical in-game display order.
- Added `List()` method to `Stats` returning numeric combat stats as a slice of `StatValue`.
- Added `PropBaseRpRecover`, `PropRpRecoverPercent`, and `PropRpRecover` property ID constants for Adrenaline Auto-Accumulation.

### Fixed
- Fixed stat naming and property ID in `UIStats.EnergyRegen` for Rupture specialty agents to correctly display Automatic Adrenaline Accumulation (`RpRecover`) instead of Energy Regen (`SpRecover`).
- Fixed base stat and bonus calculations for Rupture agents to read `PropBaseRpRecover` (`32001`).

### Tests
- Added unit test coverage for `UIStats.List()` and `Stats.List()`.
- Added unit test coverage for Rupture Adrenaline Auto-Accumulation base stat calculations and localization across languages.

## [0.10.0] - 2026-08-15

### Added
- Added slice helper functions (`AllLanguages()`, `AllAttributes()`, `AllSpecialties()`, `AllRarities()`, `AllSkillTypes()`, `AllRegions()`) for enum iteration.
- Added `IsValid()` validation methods for `Language`, `Attribute`, `Specialty`, `Rarity`, `SkillType`, and `Region`.

### Tests
- Added unit test coverage for enum helpers, validation methods, and slice immutability.

## [0.9.0] - 2026-08-14

### Changed
- Removed `panic()` on client initialization; errors are now safely returned.
- Simplified `Agent.FormattedUIStats(lang)` signature to explicitly take a target language.

### Fixed
- Fixed bug where `agent.UIStats` was localized with default store instead of client store and language.
- Fixed fallback when agent skill levels are missing in API response.

### Performance
- Optimized HTML/text formatting using single-pass `strings.NewReplacer`.
- Optimized stat breakdown formatting by using dynamic precision (`%.*f`).

### Refactored
- Centralized constants (CDN URLs, scale factors, skill indices) in `constants.go`.
- Moved SVG icons and maps into a dedicated `svg.go` file.
- Normalized JSON template files and removed obfuscated game patch tags.
- Cleaned up duplicated logic across text formatters and store loaders.

### Tests & CI
- Added GitHub Actions workflow to run `go test -race` on every push/PR.
- Embedded test fixtures with `//go:embed` to prevent path issues.
- Added test coverage for API error responses (404, 429, 500, timeouts) and formula parser edge cases.

## [0.8.2] - 2026-08-13

### Fixed
- Removed unused `iconHTMLReplacer` variable.

## [0.8.1] - 2026-08-13

### Fixed
- Fixed crash in formula evaluator on malformed math expressions.
- Sanitized HTML in `FormatHTML` to strip `<script>` tags.
- Fixed nil-pointer dereference in `Title` color helpers.

### Performance
- Reduced memory allocations in formula and icon tag parsing.

### Refactored
- Split monolithic stat and equipment calculation functions into modular helpers.

### Tests
- Added unit test coverage for stats, formatters, and API error mapping.

## [0.8.0] - 2026-08-13

### Added
- Added skill grouping (`SkillGroup`, `Agent.GroupedSkills()`) into 6 in-game UI categories.
- Added skill scaling parameter tables (`SkillParam`, `skill_templates.json`) with level formulas.
- Added pre-rendered JSON fields: `grouped_skills`, `ui_stats`, and `formatted_html`.

### Fixed
- Fixed Chain Attack & Ultimate skill level reading from raw profile.

## [0.7.0] - 2026-08-13

### Added
- Added Mindscape Cinema (M1–M6) data and formatting support.
- Added Potential Vision (special awakening buffs) data and formatting support.
- Added combat stat SVG icons and URL methods on `PropertyID` and `StatValue`.

### Changed
- Grouped drive disc set bonuses into `SetEffect` with `IsActive` state.

## [0.6.0] - 2026-08-12

### Added
- Added `AttributeDMGBonus` calculation and localized stat names for all combat properties.
- Added `TitleVariants` and dynamic argument formatting in `TitleInfo`.
- Added rich text formatting (`FormatHTML`, `FormatPlainText`, `FormatMarkdown`) to W-Engines and set bonuses.

### Fixed
- Fixed property IDs and calculation formulas for Anomaly Mastery and Proficiency.
- Added missing element mappings for Physical and Electric attributes.

### Refactored
- Simplified `calc.go` using property group helpers.

## [0.5.0] - 2026-08-11

### Added
- Added `IconURL()` and `SVG()` methods for `Attribute`, `Specialty`, and `Rarity`.
- Implemented Unity rich text parser (`<color>`, `{CAL:...}`, `<IconMap:...>`).
- Added text formatting methods to `Skill`.

### Fixed
- Fixed missing W-Engine passive and Ramielle assist localizations.

## [0.4.1] - 2026-08-11

### Fixed
- Fixed asset URL resolution for avatars, namecards, and badges.

## [0.4.0] - 2026-08-11

### Added
- Added support for `SheerForce` stat and scaling for Rupture specialty.

## [0.3.2] - 2026-08-11

### Changed
- Standardized GoDoc comments across the codebase.

## [0.3.1] - 2026-08-10

### Fixed
- Added missing localization for Rupture specialty.

## [0.3.0] - 2026-08-10

### Added
- Exported domain sentinel errors (`ErrRateLimit`, `ErrProfileNotFound`, `ErrMaintenance`, `ErrNetwork`).

## [0.2.1] - 2026-08-09

### Changed
- Improved package documentation on `pkg.go.dev`.

## [0.2.0] - 2026-08-09

### Added
- Added `UID` field to `WEngine` and `DriveDisc` structs.

## [0.1.0] - 2026-08-08

### Added
- Initial release: player profile fetching, stat calculations, drive disc analysis, and built-in metadata.
