<div align="center">
  <img src=".github/assets/fairy.png" alt="Fairy — Zenless Zone Zero profile library" width="100%">
</div>

**Fairy** is a Go library that brings full-featured [Zenless Zone Zero](https://zenless.hoyoverse.com) profile processing to your apps via the [EnkaNetwork API](https://enka.network). Just like the AI assistant from New Eridu, it handles all the heavy lifting — transforming raw game payloads into ready-to-use models with 13-language localization, resolved CDN assets, parsed Unity rich text, evaluated skill formulas, and detailed combat stat breakdowns.

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
  - [Rich text formatting](#rich-text-formatting)
  - [In-memory localization](#in-memory-localization)
- [Supported Languages](#supported-languages)

## Overview

The EnkaNetwork API returns raw game data: internal numeric IDs, unlabeled stats, and unparsed skill templates with no asset URLs. Building a UI on top of it usually means maintaining your own mapping tables, localization files, and stat formulas across game patches.

**Fairy solves this in a single call.** Provide a player UID, and Fairy returns a fully hydrated model with 13-language localization, CDN asset links, evaluated skill formulas, combat stats, and Drive Disc roll analysis.
See the difference below — raw API response vs. Fairy's enriched output:


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
    "set": {"id": 33000, "name": "Phaethon's Melody"},
    "level": 15,
    "main_stat": {"property_id": 11103, "name": "HP", "value": 2200, "is_percent": false},
    "sub_stats": [
      {"property_id": 12103, "name": "ATK",         "value": 38,   "is_percent": false, "rolls": 2},
      {"property_id": 31203, "name": "Anomaly Prof","value": 27,   "is_percent": false, "rolls": 3},
      {"property_id": 12102, "name": "Percent ATK", "value": 0.09, "is_percent": true,  "rolls": 3}
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

- 🧮 **Stat Engine** — Computes final combat stats from agent base values, W-Engines, Drive Discs, and set bonuses.
- 🎨 **UI-Ready Breakdown** — Detailed stat profiles with Base / Added / Total splits matching the in-game attributes screen.
- ⚔️ **Skill Scaling & Formulas** — Evaluates dynamic damage and daze formulas per skill level, organized into categorized skill tabs.
- 🔍 **Drive Disc Analysis** — Aggregates sub-stat totals across all discs and calculates effective roll efficiency.
- 📝 **Rich Text Engine** — Formats Unity markup into styled HTML, Markdown, or clean plain text.
- ⭐ **Mindscapes & Potential Vision** — Cinema ranks (1–6) and Potential Vision upgrade nodes with unlock state tracking.
- 👤 **Player Showcase** — Player level, server region, customizable titles with hex gradient colors, avatars, namecards, and badges.
- 🖼️ **Ready-to-Use Assets** — Splash arts (agents & skins), W-Engines, discs, badges, plus vector SVG stat & attribute icons.
- 🌍 **13 Languages** — In-memory localization for all game strings without extra network calls.
- 📦 **Zero-Config Data** — All game tables and localization files embedded directly into the binary.
- 🧩 **Flexible Client** — Functional options for custom stores, caching, timeouts, and retry policies.

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
    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

    profile, err := fairy.GetProfile(ctx, "1504687050")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Player: %s (Inter-Knot Lv.%d • %s)\n", profile.Nickname, profile.InterknotLevel, profile.Region)

    for _, agent := range profile.Agents {
        weapon := "None"
        if agent.WEngine != nil {
            weapon = agent.WEngine.Name
        }
        fmt.Printf("  • %-12s Lv.%-2d [%s/%s]  W-Engine: %s\n", 
            agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName, weapon)
    }
}
```

---

### Custom client

Use `NewClient` to set a different default language, configure HTTP timeouts, retries, caching, or a custom User-Agent header (required by Enka.Network API).

```go
client, err := fairy.NewClient(
    // Set default language for all requests
    fairy.WithDefaultLang(fairy.LangJA),
    
    // Configure the underlying enkanetwork-go HTTP client
    fairy.WithEnkaOptions(zzz.Options{
        UserAgent: "MyApp/1.0 (github.com/you/myapp)",
        HTTPClient: &http.Client{Timeout: 10 * time.Second},
        Retry:      &zzz.RetryOptions{MaxAttempts: 2, Delay: 2 * time.Second},
        Cache:      myCacheInstance,
    }),
)
if err != nil {
    log.Fatal(err)
}

profile, err := client.GetProfile(ctx, "1504687050")
```

You can also override the language on a per-request basis without recreating the client:

```go
// Use the shared global client, but respond in Korean for this call
profile, err := fairy.GetProfileWithLang(ctx, "1504687050", fairy.LangKO)
```

---

### Stat breakdown for UI
 
`agent.UIStats` provides every combat stat split into **Base**, **Added**, and **Total**, matching exactly what players see in the in-game stat panel.
 
```go
agent := profile.Agents[0]
ui := agent.UIStats
 
fmt.Printf("HP:        %s  (base %s + %s)\n", ui.HP.Total,       ui.HP.Base,       ui.HP.Added)
fmt.Printf("ATK:       %s  (base %s + %s)\n", ui.ATK.Total,      ui.ATK.Base,      ui.ATK.Added)
fmt.Printf("CRIT Rate: %s  (base %s + %s)\n", ui.CritRate.Total, ui.CritRate.Base, ui.CritRate.Added)
fmt.Printf("CRIT DMG:  %s  (base %s + %s)\n", ui.CritDMG.Total,  ui.CritDMG.Base,  ui.CritDMG.Added)
fmt.Printf("PEN Ratio: %s  (base %s + %s)\n", ui.PenRatio.Total, ui.PenRatio.Base, ui.PenRatio.Added)
```

---

### Drive Disc analysis

Measure how many sub-stat rolls landed on useful properties for your agent build.

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
    fmt.Printf("  %-22s %-8s (×%d rolls)\n", stat.Name, stat.DisplayValue(), stat.Rolls)
}
```

---

### Rich text formatting

Built-in formatters evaluate level scaling formulas and convert Unity markup (`<color>`, `<IconMap>`) into target output formats:

```go
skill := agent.Skills[0]

// For web frontend: styled <span style="color:..."> and <img> button icons
html := skill.FormatHTML()

// Markdown: **bold** highlights and clean formatting
md := skill.FormatMarkdown()

// Clean plain text: all tags stripped, formulas evaluated
plain := skill.FormatPlainText()
```

---

### In-memory localization

To display the same profile in multiple languages, fetch raw data once and localize it in memory with zero extra network calls:

```go
// 1. Fetch raw API data once
rawProfile, err := client.GetRawProfile(ctx, "1504687050")
if err != nil {
    log.Fatal(err)
}

// 2. Localize the same raw data into different languages in memory
enProfile, _ := client.Localize(rawProfile, fairy.LangEN)
jaProfile, _ := client.Localize(rawProfile, fairy.LangJA)
```

## Supported Languages

Fairy supports all 13 official in-game languages. All localization data is embedded directly into the binary and accessed in memory with zero runtime overhead or external requests.

| Language | Native Name | Constant | Code |
| :--- | :--- | :--- | :--- |
| English | English | `fairy.LangEN` | `"en"` |
| Russian | Русский | `fairy.LangRU` | `"ru"` |
| German | Deutsch | `fairy.LangDE` | `"de"` |
| Spanish | Español | `fairy.LangES` | `"es"` |
| French | Français | `fairy.LangFR` | `"fr"` |
| Indonesian | Bahasa Indonesia | `fairy.LangID` | `"id"` |
| Japanese | 日本語 | `fairy.LangJA` | `"ja"` |
| Korean | 한국어 | `fairy.LangKO` | `"ko"` |
| Portuguese | Português | `fairy.LangPT` | `"pt"` |
| Thai | ภาษาไทย | `fairy.LangTH` | `"th"` |
| Vietnamese | Tiếng Việt | `fairy.LangVI` | `"vi"` |
| Chinese (Simplified) | 简体中文 | `fairy.LangZHCN` | `"zh-cn"` |
| Chinese (Traditional) | 繁體中文 | `fairy.LangZHTW` | `"zh-tw"` |