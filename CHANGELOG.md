# Changelog

## [0.8.1] - 2026-08-13

### Fixed
- Fixed potential panic in mathematical expression evaluator (`evaluateSimpleExpr`) when encountering malformed operator sequences or missing operands.
- Added HTML sanitization in `FormatHTML` to strip `<script>` tags and prevent XSS injection from untrusted text sources.
- Fixed potential nil-pointer dereference in `Title.PrimaryColorHex()` and `Title.SecondaryColorHex()` when called on a `nil` `*Title`.

### Performance
- Optimized `{CAL:...}` formula parsing and `<IconMap:...>` tag replacement in `text.go` by preallocating buffers and eliminating redundant string allocations.

### Refactored
- Decomposed monolithic stat and equipment calculation functions in `calc.go`, `mapper.go`, and `store/embedded.go` into modular, single-responsibility helpers.

### Tests
- Added comprehensive unit test suites for `Agent.SubStatTotals`, `Agent.CountEffectiveRolls`, `DriveDisc.CountEffectiveRolls`, `Agent.FormattedUIStats`, `Title` color helpers, formatting method wrappers (`WEngine`, `MindscapeNode`, `PotentialVisionNode`, `SetEffect`), `Client` options and `Localize`, and upstream API error mapping (`mapEnkaError`).

## [0.8.0] - 2026-08-13

### Added
- Added **Combat Skill Grouping**: `SkillGroup` struct and `Agent.GroupedSkills()` method, organizing skills into 6 in-game UI category cards (`Passives`, `Basic Attack`, `Special Attack`, `Dodge`, `Chain Attack`, `Assist`) across 13 languages.
- Added **Skill Scaling Parameter Tables**: `SkillParam` struct with `Params []SkillParam` in `Skill`, extracting multiplier coefficients from `AvatarSkillTemplateTb` into `skill_templates.json` and dynamically evaluating damage %, daze %, and energy cost at current skill level.
- Added **Out-of-the-Box JSON UI Enrichment**: Pre-populated `grouped_skills` (`SkillGroups`) and `ui_stats` (`UIStats`) fields on `Agent`, and pre-rendered `formatted_html` fields on `Skill`, `MindscapeNode`, `PotentialVisionNode`, `WEngine`, and `SetEffect` for seamless `json.Marshal()` rendering.

### Fixed
- Fixed Chain Attack & Ultimate skill level mapping by reading raw `SkillLevelList` index 6.

## [0.7.0] - 2026-08-13

### Added
- Added **Mindscape Cinema (M1–M6)** support: `MindscapeNode` struct with `FormatHTML()`, `FormatPlainText()`, and `FormatMarkdown()` formatting methods, and `Agent.Mindscapes` field.
- Added **Potential Vision (Special Awakening Buffs)** support: `PotentialVision` and `PotentialVisionNode` structs, localized titles and descriptions across 13 languages, formatting methods, and mapping `IsUpgradeUnlocked` and `UpgradeID` into `Agent.PotentialVision`.
- Added **Combat Stat SVG Icons**: `SVG()` and `IconURL()` methods to `PropertyID` and `StatValue`, and added `IconURL` field to `StatValue` and `FormattedStatBreakdown` for all combat properties.

### Changed
- Refactored Drive Disc set bonus mapping into `SetEffect` struct with `IsActive` status, grouping 2-piece and 4-piece set bonuses cleanly.

## [0.6.0] - 2026-08-12

### Added
- Added `AttributeDMGBonus` calculation matching agent elemental attribute (`Physical`, `Fire`, `Ice`, `Electric`, `Ether`, `Wind`) and reskins (`AuricInk`, `HonedEdge`, `Frost`).
- Added `Name` field to `FormattedStatBreakdown` and localized all 13 combat stats.
- Added `TitleVariants` support and dynamic `TitleInfo` argument formatting in `Profile`.
- Added `FormatHTML()`, `FormatPlainText()`, `FormatMarkdown()`, and `{CAL:...}` formula evaluation to `WEngine` and `DriveDiscSetBonus`.

### Fixed
- Fixed Anomaly Mastery (`31401`–`31403`) and Anomaly Proficiency (`31201`–`31203`) property IDs and calculation formulas.
- Added explicit element type mappings for `Physical` and `Electric` attributes in `mapper.go`.

### Changed
- Refactored `calc.go` to use `sumPropVariants` helper and `propGroup` constants, removing raw property ID calculations.
- Replaced `if-else` chains in formula evaluator with tagged `switch` statements in `text.go`.

## [0.5.0] - 2026-08-11

### Added
- Added `IconURL()` method to `Attribute`, `Specialty`, and `Rarity` types for obtaining visual icon URLs.
- Added `SVG()` method to `Attribute` returning clean, W3C-compliant vector SVG markup for all 10 ZZZ attributes.
- Implemented Unity Rich Text formatting parser (`FormatHTML`, `FormatPlainText`, `FormatMarkdown`) supporting `<color>`, `{CAL:...}`, and `<IconMap:...>` tags.
- Added `FormatHTML()`, `FormatPlainText()`, and `FormatMarkdown()` methods to `Skill` struct.

### Fixed
- Fixed W-Engine passive description property key extraction and missing localizations.
- Fixed missing assist skill localizations for Ramielle.

## [0.4.1] - 2026-08-11

### Fixed
- Fixed profile avatar, namecard background, and badge asset URL resolution.
- Fixed profile avatar mapping by prioritizing `ProfileID` over `AvatarID` to resolve in-game Inter-Knot avatars properly.

## [0.4.0] - 2026-08-11

### Added
- Added support for `SheerForce` stat in `Stats`, `FormattedStats`, and `UIStats`.
- Implemented HP and ATK stat scaling for Rupture specialty agents.

## [0.3.2] - 2026-08-11

### Changed
- Cleaned up and standardized GoDoc comments across the codebase.

## [0.3.1] - 2026-08-10

### Fixed
- Fixed missing localization for the `Rupture` specialty across all supported languages.

## [0.3.0] - 2026-08-10

### Added
- Exported domain sentinel errors (`ErrRateLimit`, `ErrProfileNotFound`, `ErrMaintenance`, `ErrNetwork`) from `internal/api` to the root `fairy` package.

## [0.2.1] - 2026-08-09

### Changed
- Improved package documentation for `pkg.go.dev` by adding individual docstrings to `Attribute`, `Specialty`, `Rarity`, and `PropertyID` constants.
- Clarified the default behavior of the metadata store in `NewClient` documentation.

## [0.2.0] - 2026-08-09

### Added
- Added `UID` field to `WEngine` and `DriveDisc` structs to allow tracking specific instances of equipment across agents.

## [0.1.0] - 2026-08-08

### Added
- Fetch, enrich, and calculate Zenless Zone Zero player game profiles.
- Stat calculation for characters, weapons, and drive discs.
- UI-ready formatting for in-game stat panels.
- Drive disc analysis with effective rolls counting.
- Auto-resolved URLs for various game assets.
- Built-in metadata (no database required).
- Modular client support with `enkanetwork-go` features.
