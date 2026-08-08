<div align="center">
  <img src=".github/assets/fairy.png" alt="Fairy — Zenless Zone Zero profile library" width="100%">
</div>

<br>

**Fairy** is a Go library for fetching, enriching, and calculating [Zenless Zone Zero](https://zenless.hoyoverse.com) player game profiles via the [EnkaNetwork API](https://enka.network). Just like the AI assistant from New Eridu, she handles all the heavy lifting — mapping raw game IDs into localized names, building asset URLs, and computing final combat stats from first principles.

[![Go Reference](https://pkg.go.dev/badge/github.com/kirinyoku/fairy.svg)](https://pkg.go.dev/github.com/kirinyoku/fairy)
[![Go Version](https://img.shields.io/github/go-mod/go-version/kirinyoku/fairy)](https://golang.org/doc/devel/release.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Table of Contents

- [Why Fairy?](#-why-fairy)
- [Features](#-features)
- [Installation](#-installation)
- [Usage](#-usage)
- [Supported Languages](#-supported-languages)
- [License](#-license)

<br>

## 🤔 Why Fairy?

The raw EnkaNetwork API returns internal game IDs and raw scaled integers, instead of readable names and values. Fairy maps all of it automatically.

<br>

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

## ✨ Features

| Feature | Description |
|---|---|
| 🧮 **Stat Calculation** | Computes final stats from base values, W-Engine scalings, drive disc main stats, and substat rolls |
| 🎨 **UI-Ready Formatting** | `FormattedStats` and `UIStats` (Base + Added) match the in-game stat panel exactly |
| 🔍 **Drive Disc Analysis** | `SubStatTotals()` and `CountEffectiveRolls()` to evaluate drive disc quality in one call |
| 🌍 **13 Languages** | Full localization of agent names, skills, drive disc sets, stats, and titles |
| 🖼️ **Asset URLs** | Auto-resolved URLs for splash arts, icons, namecards, avatars, and badges |
| 📦 **Embedded Data** | Ships with all game metadata built-in — no extra setup or database required |
| 🧩 **Modular Client** | Configurable HTTP options, custom User-Agent, built-in API caching, and per-request language overrides |

## 📦 Installation

Requires **Go 1.22+**

```bash
go get github.com/kirinyoku/fairy
```

## 🚀 Usage

### Fetching a profile

The simplest way to get started is with the package-level `GetProfile` function.
It uses a shared default client with **English** localization and embedded game data.

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

Use `NewClient` when you need a specific language, custom HTTP settings, caching, or a custom User-Agent header (required by Enka.Network's API guidelines).

```go
// You can use the built-in memory cache from the enkanetwork-go library
// to prevent rate-limiting and speed up duplicate profile requests.
import "github.com/kirinyoku/enkanetwork-go/cache"

client, err := fairy.NewClient(
    fairy.WithDefaultLang(fairy.LangJA),
    fairy.WithEnkaOptions(zzz.Options{
        UserAgent: "MyApp/1.0 (github.com/you/myapp)",
        Cache:     cache.NewMemoryCache(),
    }),
)
```

You can also override the language per-request without recreating the client:

```go
// Use the shared global client, but respond in Korean for this call
profile, err := fairy.GetProfileWithLang(ctx, "1504687050", fairy.LangKO)
```

---

### Stat breakdown for UI

`FormattedUIStats()` returns every stat split into **Base**, **Added**, and **Total** — formatted as strings exactly like the in-game stat panel.

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

### Disc analysis

Count how many sub-stat rolls landed on stats that actually matter for your agent's build.

```go
agent := profile.Agents[0]

// Count effective rolls for ATK specialized agent:
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

If you need to display the same profile in multiple languages, you can make a single HTTP request to fetch the raw data, and then localize it efficiently in memory:

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

## 🌍 Supported Languages

| Language | Language |
|----------|----------|
| 🇬🇧 English | 🇰🇷 Korean |
| 🇷🇺 Russian | 🇵🇹 Portuguese |
| 🇩🇪 German | 🇹🇭 Thai |
| 🇪🇸 Spanish | 🇻🇳 Vietnamese |
| 🇫🇷 French | 🇨🇳 Chinese (Simplified) |
| 🇮🇩 Indonesian | 🇹🇼 Chinese (Traditional) |
| 🇯🇵 Japanese | |

## ⚖️ License

Licensed under the [MIT License](LICENSE).
