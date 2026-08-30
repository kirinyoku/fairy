<div align="center">
  <img src=".github/assets/fairy.png" alt="Fairy — Zenless Zone Zero profile library" width="100%">

  # Fairy

**A Go library for fetching and enriching Zenless Zone Zero player profiles**

  [![Go Reference](https://pkg.go.dev/badge/github.com/kirinyoku/fairy.svg)](https://pkg.go.dev/github.com/kirinyoku/fairy)
  [![Go Version](https://img.shields.io/github/go-mod/go-version/kirinyoku/fairy)](https://golang.org/doc/devel/release.html)
  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

  *Fetch and enrich Zenless Zone Zero player profiles via the [EnkaNetwork API](https://enka.network) with localized names, calculated agent stats, and ready-to-use assets.*
</div>

---

## At a Glance

```go
profile, err := fairy.GetProfile(ctx, "1504687050")

fmt.Println(profile.Nickname)                 // "LOWLEVEL"
fmt.Println(profile.Agents[0].Name)           // "Nangong Yu"
fmt.Println(profile.Agents[0].Stats.CritRate) // 0.074 (7.4%)
```

---

## Installation

Requires **Go 1.22+**

```bash
go get github.com/kirinyoku/fairy
```

---

## Quick Start

Fetch a player's showcase profile and inspect their agents with zero configuration:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/kirinyoku/fairy"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    profile, err := fairy.GetProfile(ctx, "1504687050")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Player: %s (Inter-Knot Lv.%d • %s)\n",
        profile.Nickname, profile.InterknotLevel, profile.Region)

    for _, agent := range profile.Agents {
        fmt.Printf("  • %s Lv.%d [%s / %s]\n",
            agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName)
    }
}
```

```text
Player: LOWLEVEL (Inter-Knot Lv.60 • Europe)
  • Nangong Yu Lv.60 [Ether / Stun]
  • Yixuan Lv.60 [Auric Ink / Rupture]
  • Ye Shunguang Lv.60 [Honed Edge / Attack]
  • Harumasa Lv.60 [Electric / Attack]
  • Miyabi Lv.60 [Frost / Anomaly]
  • Remielle Lv.60 [Lumiflux / Anomaly]
```

---

## Why Fairy?

The raw EnkaNetwork API returns internal numeric IDs, unlabeled stat keys, and unparsed Unity formula strings. Building apps directly on top of it requires maintaining mapping tables and stat formulas across every game patch.

**Fairy handles all of this automatically in-memory.**

<details>
<summary><b>🔍 Click to compare: Raw EnkaNetwork API vs Fairy Enriched Output</b></summary>
<br>

<table>
<tr>
<th>Raw EnkaNetwork API</th>
<th>Fairy Enriched Output</th>
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
  "drive_discs": {
    "slots": [{
      "slot": 1,
      "set": {"id": 33000, "name": "Phaethon's Melody"},
      "level": 15,
      "main_stat": {"property_id": 11103, "name": "HP", "value": 2200, "is_percent": false},
      "sub_stats": [
        {"property_id": 12103, "name": "ATK",          "value": 38,   "is_percent": false, "rolls": 2},
        {"property_id": 31203, "name": "Anomaly Prof", "value": 27,   "is_percent": false, "rolls": 3},
        {"property_id": 12102, "name": "ATK",          "value": 0.09, "is_percent": true,  "rolls": 3}
      ]
    }],
    "set_bonuses": [{
      "set": {"id": 33000, "name": "Phaethon's Melody"},
      "count": 4
    }]
  },
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

> *Note: JSON snippets are simplified for illustration. Fairy's models provide complete asset URLs, formulas, and breakdown fields.*

</details>

---

## Core API

Fairy provides global functions for quick one-liners, and a [`Client`](https://pkg.go.dev/github.com/kirinyoku/fairy#Client) struct for full control over networking, caching, and default language.

| Method | Description | Network |
| :--- | :--- | :---: |
| [`fairy.GetProfile(ctx, uid)`](https://pkg.go.dev/github.com/kirinyoku/fairy#GetProfile) | Fetch and enrich player profile in default language | 🌐 HTTP |
| [`fairy.GetProfileWithLang(ctx, uid, lang)`](https://pkg.go.dev/github.com/kirinyoku/fairy#GetProfileWithLang) | Fetch and enrich profile in a specific language ([`fairy.LangJA`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangJA), [`fairy.LangRU`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangRU), etc.) | 🌐 HTTP |
| [`fairy.GetRawProfile(ctx, uid)`](https://pkg.go.dev/github.com/kirinyoku/fairy#GetRawProfile) | Fetch raw unparsed profile directly from EnkaNetwork API | 🌐 HTTP |
| [`fairy.Enrich(raw)`](https://pkg.go.dev/github.com/kirinyoku/fairy#Enrich) | Transform raw profile into an enriched model in default language | ⚡ In-memory |
| [`fairy.EnrichWithLang(raw, lang)`](https://pkg.go.dev/github.com/kirinyoku/fairy#EnrichWithLang) | Transform raw profile into an enriched model in a specific language ([`fairy.LangJA`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangJA), [`fairy.LangRU`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangRU), etc.) | ⚡ In-memory |

---

## Feature Highlights

### 1. Agent Stats — 3 Flexible Representations
Fairy calculates the complete combat stat sheet from base attributes, W-Engines, and Drive Discs:

| Mode | Accessor | Type | Best for |
| :--- | :--- | :--- | :--- |
| **Numeric** | `agent.Stats` | [`Stats`](https://pkg.go.dev/github.com/kirinyoku/fairy#Stats) | Math calculations, damage simulators (`CritRate: 0.05`) |
| **Formatted** | `agent.Stats.Formatted()` | [`FormattedStats`](https://pkg.go.dev/github.com/kirinyoku/fairy#FormattedStats) | Text output, logs, summaries (`CritRate: "5.0%"`) |
| **UI Breakdown** | `agent.UIStats` | [`UIStats`](https://pkg.go.dev/github.com/kirinyoku/fairy#UIStats) | In-game style `Base` + `Added` = `Total` with SVG icons & localized names |

📖 See [`ExampleUIStats_List`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-UIStats.List)

---

### 2. Drive Disc Analysis & Build Scoring
Deep breakdown of Drive Disc sets, substat aggregations, and roll quality:

- **Partition Slots (1–6):** Main stats, substats, upgrade rolls, and set identifiers via [`agent.DriveDiscs.Slots`](https://pkg.go.dev/github.com/kirinyoku/fairy#DriveDiscs).
- **Set Bonuses:** Active 2-pc and 4-pc set bonuses via [`agent.DriveDiscs.SetBonuses`](https://pkg.go.dev/github.com/kirinyoku/fairy#DriveDiscs), or boolean queries via [`Has4Piece(setID)`](https://pkg.go.dev/github.com/kirinyoku/fairy#DriveDiscs.Has4Piece) / [`Has2Piece(setID)`](https://pkg.go.dev/github.com/kirinyoku/fairy#DriveDiscs.Has2Piece) (e.g. [`fairy.SetPolarMetal`](https://pkg.go.dev/github.com/kirinyoku/fairy#SetPolarMetal)).
- **Substat Totals:** Aggregates all rolls and stat values across all discs with [`agent.DriveDiscs.SubStatTotals()`](https://pkg.go.dev/github.com/kirinyoku/fairy#DriveDiscs.SubStatTotals).
- **Roll Scoring:** Count effective rolls for specific priority stats with [`agent.DriveDiscs.CountEffectiveRolls(...)`](https://pkg.go.dev/github.com/kirinyoku/fairy#DriveDiscs.CountEffectiveRolls).

📖 See [`ExampleDriveDiscs_SubStatTotals`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-DriveDiscs.SubStatTotals)

---

### 3. Skills, Scaling & Rich Text
- **Flat & Grouped Views:** Access all abilities via [`agent.Skills`](https://pkg.go.dev/github.com/kirinyoku/fairy#Agent.Skills), or categorized into the 6 in-game UI tabs (`basic`, `special`, `dodge`, `chain`, `assist`, `passive`) with upgrade levels (1–12 active, 0–6 core) via [`agent.SkillGroups`](https://pkg.go.dev/github.com/kirinyoku/fairy#Agent.SkillGroups).
- **Dynamic Formulas:** Scaling values and daze ratios are evaluated dynamically per skill level.
- **Unity Rich Text:** Built-in converters for in-game markup: [`FormatHTML()`](https://pkg.go.dev/github.com/kirinyoku/fairy#Skill.FormatHTML), [`FormatMarkdown()`](https://pkg.go.dev/github.com/kirinyoku/fairy#Skill.FormatMarkdown), or [`FormatPlainText()`](https://pkg.go.dev/github.com/kirinyoku/fairy#Skill.FormatPlainText).

📖 See [`ExampleAgent_SkillGroups`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-Agent_SkillGroups)

---

### 4. Mindscape Cinema & Potential Vision
- **Mindscape Cinema:** Constellation rank (M0–M6) via [`agent.MindscapeCinema`](https://pkg.go.dev/github.com/kirinyoku/fairy#Agent).
- **Node Tracking:** Access all 6 Mindscape nodes with unlock status and formatted effect descriptions via [`agent.Mindscapes`](https://pkg.go.dev/github.com/kirinyoku/fairy#Agent.Mindscapes).
- **Potential Vision:** Character upgrade nodes with active unlock states via [`agent.PotentialVision`](https://pkg.go.dev/github.com/kirinyoku/fairy#Agent.PotentialVision).

📖 See [`ExampleAgent_Mindscapes`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-Agent_Mindscapes)

---

### 5. Player Showcase & Visual Assets
- **Profile Info:** Player Nickname, Inter-Knot Level, Server Region, and Server Cache TTL via [`profile.CacheTTL()`](https://pkg.go.dev/github.com/kirinyoku/fairy#Profile.CacheTTL).
- **Gradient Titles:** Two-color gradient titles with hex color helpers ([`PrimaryColorHex()`](https://pkg.go.dev/github.com/kirinyoku/fairy#Title.PrimaryColorHex), [`SecondaryColorHex()`](https://pkg.go.dev/github.com/kirinyoku/fairy#Title.SecondaryColorHex)).
- **Media CDN Assets:** Direct URLs for high-resolution splash art (agents & skins), W-Engines, discs, badges, namecards, plus vector SVG stat & attribute icons.

📖 See [`ExampleProfile`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-Profile)

---

## Advanced Configuration & Errors

<details>
<summary><b>🔧 Custom Client (Timeouts, Retries, Caching, User-Agent)</b></summary>
<br>

Use [`fairy.NewClient`](https://pkg.go.dev/github.com/kirinyoku/fairy#NewClient) with functional options for production deployments:

```go
// import "github.com/kirinyoku/enkanetwork-go/client/zzz"

client, err := fairy.NewClient(
    fairy.WithDefaultLang(fairy.LangJA),
    fairy.WithEnkaOptions(zzz.Options{
        UserAgent:  "MyApp/1.0 (contact@example.com)",
        HTTPClient: &http.Client{Timeout: 10 * time.Second},
        Retry:      &zzz.RetryOptions{MaxAttempts: 2, Delay: 2 * time.Second},
        // Cache: myCacheInstance, // plug in your own zzz.Cache implementation
    }),
)

profile, err := client.GetProfile(ctx, "1504687050")
```

</details>

<details>
<summary><b>🚨 Error Handling & Sentinel Errors</b></summary>
<br>

Fairy returns structured sentinel errors for reliable error handling with `errors.Is`:

| Error | Description |
| :--- | :--- |
| [`fairy.ErrInvalidUID`](https://pkg.go.dev/github.com/kirinyoku/fairy#ErrInvalidUID) | UID format is invalid (must be 10 numeric digits starting with 10, 13, 15, or 17). |
| [`fairy.ErrProfileNotFound`](https://pkg.go.dev/github.com/kirinyoku/fairy#ErrProfileNotFound) | The requested profile does not exist or has showcase hidden. |
| [`fairy.ErrRateLimit`](https://pkg.go.dev/github.com/kirinyoku/fairy#ErrRateLimit) | EnkaNetwork API rate limit exceeded. |
| [`fairy.ErrMaintenance`](https://pkg.go.dev/github.com/kirinyoku/fairy#ErrMaintenance) | API or game servers are under maintenance. |
| [`fairy.ErrNetwork`](https://pkg.go.dev/github.com/kirinyoku/fairy#ErrNetwork) | Network-level failure (DNS, connection timeout, etc.). |
| [`fairy.ErrEnrichment`](https://pkg.go.dev/github.com/kirinyoku/fairy#ErrEnrichment) | In-memory metadata mapping or formula evaluation error. |

</details>

---

## Supported Languages

Fairy includes **zero-config in-memory localization** for all 13 official languages — embedded directly into the binary with no extra network requests.

<details>
<summary><b>🌐 Supported Languages Table</b></summary>
<br>

| Language | Native Name | Constant | Code |
| :--- | :--- | :--- | :--- |
| English | English | [`fairy.LangEN`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangEN) | `"en"` |
| Russian | Русский | [`fairy.LangRU`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangRU) | `"ru"` |
| Japanese | 日本語 | [`fairy.LangJA`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangJA) | `"ja"` |
| Chinese (Simplified) | 简体中文 | [`fairy.LangZHCN`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangZHCN) | `"zh-cn"` |
| Chinese (Traditional) | 繁體中文 | [`fairy.LangZHTW`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangZHTW) | `"zh-tw"` |
| Korean | 한국어 | [`fairy.LangKO`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangKO) | `"ko"` |
| German | Deutsch | [`fairy.LangDE`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangDE) | `"de"` |
| French | Français | [`fairy.LangFR`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangFR) | `"fr"` |
| Spanish | Español | [`fairy.LangES`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangES) | `"es"` |
| Portuguese | Português | [`fairy.LangPT`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangPT) | `"pt"` |
| Indonesian | Bahasa Indonesia | [`fairy.LangID`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangID) | `"id"` |
| Thai | ภาษาไทย | [`fairy.LangTH`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangTH) | `"th"` |
| Vietnamese | Tiếng Việt | [`fairy.LangVI`](https://pkg.go.dev/github.com/kirinyoku/fairy#LangVI) | `"vi"` |

</details>

---

## License

Fairy is released under the [MIT License](LICENSE).
