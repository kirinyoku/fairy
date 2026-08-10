# Changelog

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
