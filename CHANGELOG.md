# Changelog

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
