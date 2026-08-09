<div align="center">
  <img src=".github/assets/fairy.png" alt="Fairy — Zenless Zone Zero profile library" width="100%">
</div>

**Fairy** is a Go library for fetching, enriching, and calculating [Zenless Zone Zero](https://zenless.hoyoverse.com) player game profiles via the [EnkaNetwork API](https://enka.network). Just like the AI assistant from New Eridu, it handles all the heavy lifting — mapping raw game IDs to localized names, building asset URLs, and computing final combat stats from scratch.

[![Go Reference](https://pkg.go.dev/badge/github.com/kirinyoku/fairy.svg)](https://pkg.go.dev/github.com/kirinyoku/fairy)
[![Go Version](https://img.shields.io/github/go-mod/go-version/kirinyoku/fairy)](https://golang.org/doc/devel/release.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Fetching a profile](#fetching-a-profile)
  - [Custom client](#custom-client)
  - [Stat breakdown for UI](#stat-breakdown-for-ui)
  - [Drive Disc analysis](#drive-disc-analysis)
  - [Localize raw data yourself](#localize-raw-data-yourself)
- [Supported Languages](#supported-languages)
- [License](#license)

## Overview

The EnkaNetwork API returns player profiles as raw data — agents, W-Engines, and Drive Discs are represented by internal numeric IDs, stat values have no names, and there are no image URLs. To build anything user-facing, you'd need to maintain your own mapping tables, host localization files, implement stat calculations, and keep up with every game patch.

Fairy eliminates this entire layer. It takes a single UID, fetches the raw profile from Enka, and returns a fully enriched model — human-readable names in 13 languages, ready-to-use asset URLs, computed final stats, and Drive Disc roll analysis. One function call, zero boilerplate.

The comparison below shows what this looks like in practice: a raw API response on the left versus the enriched output Fairy produces on the right.

<table>
<tr>
<th>ENKANETWORK API RESPONSE</th>
<th>FAIRY ENRICHED OUTPUT</th>
</tr>
<tr>
<td>

```json
{
  "Id": 1511,
  "Level": 60,
  "Exp": 0,
  "PromotionLevel": 6,
  "TalentLevel": 0,
  "SkinId": 3115111,
  "UpgradeId": 0,
  "CoreSkillEnhancement": 6,
  "Weapon": {
    "Id": 15388,
    "Level": 60,
    "StarMark": 1,
    "BreakLevel": 6
  },
  "EquippedList": [{
    "Slot": 1,
    "Equipment": {
      "Id": 33041,
      "Level": 15,
      "MainPropertyList": [{
        "PropertyId": 11103,
        "PropertyValue": 550
      }],
      "RandomPropertyList": [
        {"PropertyId": 12103, "PropertyValue": 19},
        {"PropertyId": 31203, "PropertyValue": 9},
        {"PropertyId": 11102, "PropertyValue": 300},
        {"PropertyId": 12102, "PropertyValue": 300}
      ]
    }
  }]
}
```

</td>
<td>

```json
{
  "name": "Nangong Yu",
  "level": 60,
  "rarity": "S",
  "attribute_name": "Ether",
  "specialty_name": "Stun",
  "w_engine": {
    "name": "Neon Fantasies",
    "level": 60,
    "modification": 1,
    "rarity": "S",
    "main_stat": {
      "name": "Base ATK",
      "value": 713
    }
  },
  "drive_discs": [{
    "slot": 1,
    "set_name": "Phaethon's Melody",
    "level": 15,
    "main_stat": {"name": "HP", "value": 2200},
    "sub_stats": [
      {"name": "ATK",         "value": 38,   "rolls": 2},
      {"name": "Anomaly Prof","value": 27,   "rolls": 3},
      {"name": "Percent ATK", "value": 0.09, "rolls": 3}
    ]
  }],
  "stats": {
    "hp": 11188, "atk": 2866,
    "crit_rate": 0.074, "crit_dmg": 0.548,
    "pen_ratio": 0.24, "energy_regen": 1.2
  }
}
```

</td>
</tr>
</table>

> [!NOTE]
> The JSON examples above are **simplified and shortened** to highlight the key differences. The actual API responses and Fairy's models contain significantly more data.

## Features

- 🧮 **Stat Calculation** — Computes final combat stats (HP, ATK, DEF, CRIT, PEN, Energy Regen, and more) by combining agent base values, W-Engine scaling, Drive Disc main/sub stats, and set bonuses. All percentage stats are stored as decimals internally and can be formatted for display with a single call.

- 🎨 **UI-Ready Stat Breakdown** — `FormattedUIStats()` splits every stat into **Base**, **Added**, and **Total** components, pre-formatted as strings — matching exactly what players see in the in-game stat panel. `Stats.Formatted()` gives you a simpler flat view when you don't need the breakdown.

- 🔍 **Drive Disc Analysis** — `SubStatTotals()` aggregates sub-stats across all six discs, grouping by property and summing values and rolls. `CountEffectiveRolls()` counts how many rolls landed on the stats you care about — available on both the agent (all discs) and individual disc level.

- 🌍 **13 Languages** — Every string in the output — agent names, skill descriptions, stat labels, W-Engine passives, set bonus text, titles, and badges — is fully localized. Fetch raw data once, then call `Localize()` to produce the same profile in any supported language without extra network calls.

- 🖼️ **Asset URLs** — Generates ready-to-use image URLs for agent splash arts, skins, W-Engine icons, Drive Disc icons, profile avatars, namecards, and badges. No manual URL construction needed.

- 📦 **Zero-Config Data** — All game metadata (stat scaling tables, localization strings, item definitions) is embedded in the binary via `go:embed`. No external files, no database, no CDN — just `go get` and start building.

- 🧩 **Flexible Client** — Functional options let you configure the default language, swap in a custom `MetadataStore` implementation, and pass through HTTP settings (timeouts, retries, User-Agent, caching) to the underlying [`enkanetwork-go`](https://github.com/kirinyoku/enkanetwork-go) client.

## Installation

Requires **Go 1.22+**

```bash
go get github.com/kirinyoku/fairy
```

## Usage

### Fetching a profile

The easiest way to get started is with the `GetProfile` function. It uses a default client with **English** localization and built-in game data.

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/kirinyoku/fairy"
)

func main() {
    profile, err := fairy.GetProfile(context.Background(), "1504687050")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Player: %s (Inter-Knot Level: %d)\n", profile.Nickname, profile.InterknotLevel)

    for _, agent := range profile.Agents {
        fmt.Printf("  • %s  Lv.%d  %s %s\n", agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName)
    }
}
```

---

### Custom client

Use `NewClient` to set a different default language, configure HTTP settings, retries, caching, or a custom User-Agent header (required by Enka.Network).

```go
client, err := fairy.NewClient(
    // Change default language for all requests
    fairy.WithDefaultLang(fairy.LangJA),
    
    // Configure the underlying `github.com/kirinyoku/enkanetwork-go/client/zzz` client
    fairy.WithEnkaOptions(zzz.Options{
        UserAgent: "MyApp/1.0 (github.com/you/myapp)",
        HTTPClient: &http.Client{Timeout: 10 * time.Second},
        Retry:      &zzz.RetryOptions{MaxAttempts: 2, Delay: 2 * time.Second},
        Cache:      myCacheInstance,
    }),
)
```

You can also override the language on a per-request basis without recreating the client:

```go
// Use the shared global client, but respond in Korean for this call
profile, err := fairy.GetProfileWithLang(ctx, "1504687050", fairy.LangKO)
```

---

### Stat breakdown for UI

`FormattedUIStats()` returns every stat split into **Base**, **Added**, and **Total**, formatted exactly as they appear in the in-game stat panel.

```go
agent := profile.Agents[0]
ui := agent.FormattedUIStats()

fmt.Printf("HP:        %s  (base %s + %s)\n", ui.HP.Total,        ui.HP.Base,        ui.HP.Added)
fmt.Printf("ATK:       %s  (base %s + %s)\n", ui.ATK.Total,       ui.ATK.Base,       ui.ATK.Added)
fmt.Printf("CRIT Rate: %s  (base %s + %s)\n", ui.CritRate.Total,  ui.CritRate.Base,  ui.CritRate.Added)
fmt.Printf("CRIT DMG:  %s  (base %s + %s)\n", ui.CritDMG.Total,   ui.CritDMG.Base,   ui.CritDMG.Added)
fmt.Printf("PEN Ratio: %s  (base %s + %s)\n", ui.PenRatio.Total,  ui.PenRatio.Base,  ui.PenRatio.Added)
```

---

### Drive Disc analysis

Measure how many sub-stat rolls landed on stats that actually matter for your agent.

```go
agent := profile.Agents[0]

// Count effective rolls for an Attack agent:
usefulRolls := agent.CountEffectiveRolls(
    fairy.PropCritRate,
    fairy.PropCritDMG,
    fairy.PropATKPercent,
)
fmt.Printf("Effective rolls: %d\n", usefulRolls)

// Full sub-stat breakdown across all 6 discs, grouped and summed:
for _, stat := range agent.SubStatTotals() {
    fmt.Printf("  %-20s %s  (×%d rolls)\n", stat.Name, stat.DisplayValue(), stat.Rolls)
}
```

---

### Localize raw data yourself

To display the same profile in multiple languages, fetch the raw data once and localize it in memory — no extra network calls needed:

```go
// 1. Fetch the raw data from the API just once
rawProfile, err := client.GetRawProfile(ctx, "1504687050")
if err != nil {
    log.Fatal(err)
}

// 2. Localize the same raw data into different languages without extra network calls
enProfile, _ := client.Localize(rawProfile, fairy.LangEN)
jaProfile, _ := client.Localize(rawProfile, fairy.LangJA)
```

## Supported Languages

| Language | Language |
|----------|----------|
| 🇬🇧 English | 🇰🇷 Korean |
| 🇷🇺 Russian | 🇵🇹 Portuguese |
| 🇩🇪 German | 🇹🇭 Thai |
| 🇪🇸 Spanish | 🇻🇳 Vietnamese |
| 🇫🇷 French | 🇨🇳 Chinese (Simplified) |
| 🇮🇩 Indonesian | 🇹🇼 Chinese (Traditional) |
| 🇯🇵 Japanese | |

## License

Licensed under the [MIT License](LICENSE).
