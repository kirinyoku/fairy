<div align="center">
  <img src=".github/assets/fairy.png" alt="Fairy — Zenless Zone Zero profile library" width="100%">
</div>

**Fairy** is a Go library that brings full-featured [Zenless Zone Zero](https://zenless.hoyoverse.com) profile processing to your apps via the [EnkaNetwork API](https://enka.network). Just like the AI assistant from New Eridu, it handles all the heavy lifting — transforming raw game payloads into ready-to-use models with 13-language localization, CDN asset URLs, parsed Unity rich text, evaluated skill formulas, and detailed combat stat breakdowns.

[![Go Reference](https://pkg.go.dev/badge/github.com/kirinyoku/fairy.svg)](https://pkg.go.dev/github.com/kirinyoku/fairy)
[![Go Version](https://img.shields.io/github/go-mod-go-version/kirinyoku/fairy)](https://golang.org/doc/devel/release.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Core actions](#core-actions)
  - [Global functions vs Client](#global-functions-vs-client)
  - [Quick start](#quick-start)
  - [Custom client](#custom-client)
  - [Error handling](#error-handling)
  - [Feature highlights](#feature-highlights)
    - [Agent stats](#agent-stats)
    - [Agent skills](#agent-skills)
    - [Drive Disc analysis](#drive-disc-analysis)
    - [Rich text formatting](#rich-text-formatting)
- [Supported Languages](#supported-languages)
- [License](#license)

## Overview

The EnkaNetwork API returns raw game data: internal numeric IDs, unlabeled stats, and unparsed skill templates with no asset URLs. Building a UI on top of it usually means maintaining your own mapping tables, localization files, and stat formulas across game patches.

**Fairy solves this in a single call.** Provide a player UID, and Fairy returns a fully hydrated model with 13-language localization, CDN asset links, evaluated skill formulas, combat stats, and Drive Disc roll analysis. The comparison below shows the difference — raw API response vs. Fairy's enriched output:


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
- 🎨 **UI-Ready Breakdown** — Base / Added / Total stat splits with localized names and SVG icons, matching the in-game attributes screen.
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

Full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/kirinyoku/fairy). However, auto-generated docs list every type and function alphabetically, which makes it hard to see the big picture or know where to start. The guide below walks you through the library's key concepts in a logical order — from basic fetching to advanced features — so you can get productive quickly.

### Core actions

`fairy` exposes four actions. The first three hit the [EnkaNetwork API](https://enka.network) over HTTP; the fourth is a pure in-memory transformation.

#### `GetProfile`

Fetches a player's profile from the API and returns a fully enriched [`*Profile`](https://pkg.go.dev/github.com/kirinyoku/fairy#Profile) in the client's default language (English for the global function).

```go
profile, err := fairy.GetProfile(ctx, "1504687050")
```

📖 See [`ExampleGetProfile`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-GetProfile)

#### `GetProfileWithLang`

Same as `GetProfile`, but overrides the localization language for this single request.

```go
profile, err := fairy.GetProfileWithLang(ctx, "1504687050", fairy.LangJA)
```

📖 See [`ExampleGetProfileWithLang`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-GetProfileWithLang)

#### `GetRawProfile`

Fetches the raw [`*zzz.Profile`](https://pkg.go.dev/github.com/kirinyoku/enkanetwork-go/client/zzz#Profile) from the API without applying any enrichment or localization.

```go
raw, err := fairy.GetRawProfile(ctx, "1504687050")
```

#### `Localize`

Maps a raw [`*zzz.Profile`](https://pkg.go.dev/github.com/kirinyoku/enkanetwork-go/client/zzz#Profile) into an enriched [`*Profile`](https://pkg.go.dev/github.com/kirinyoku/fairy#Profile) for the specified language. **Makes zero network calls.**

```go
enProfile, err := fairy.Localize(raw, fairy.LangEN)
jaProfile, err := fairy.Localize(raw, fairy.LangJA)
```

📖 See [ExampleLocalize](https://pkg.go.dev/github.com/kirinyoku/fairy#example-Localize).

---

### Global functions vs Client

Every action is available in two forms:

**Global functions** — `fairy.GetProfile(ctx, uid)`, `fairy.Localize(raw, lang)`, etc. — use a lazily initialized shared client with default settings.

**Client methods** — created via [`fairy.NewClient`](https://pkg.go.dev/github.com/kirinyoku/fairy#NewClient) — give you full control over HTTP client configuration, retries, and caching.

```go
// Global function
profile, err := fairy.GetProfile(ctx, "1504687050")

// Client method
client, err := fairy.NewClient()
profile, err := client.GetProfile(ctx, "1504687050")
```

---

### Quick start

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
        fmt.Printf("  %s Lv.%d [%s/%s]\n",
            agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName)
    }
}
```

Example Output:

```
Player: LOWLEVEL (Inter-Knot Lv.60 • Europe)
  Nangong Yu Lv.60 [Ether/Stun]
  Yixuan Lv.60 [Auric Ink/Rupture]
  Ye Shunguang Lv.60 [Honed Edge/Attack]
  Harumasa Lv.60 [Electric/Attack]
  Miyabi Lv.60 [Frost/Anomaly]
  Remielle Lv.60 [Lumiflux/Anomaly]
```

---

### Custom client

Use `NewClient` to set a different default language, configure HTTP timeouts, retries, caching, or a custom User-Agent header.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/kirinyoku/enkanetwork-go/client/zzz"
    "github.com/kirinyoku/fairy"
)

func main() {
    client, err := fairy.NewClient(
        fairy.WithDefaultLang(fairy.LangJA),
        fairy.WithEnkaOptions(zzz.Options{
            UserAgent:  "MyApp/1.0 (contact@example.com)",
            HTTPClient: &http.Client{Timeout: 10 * time.Second},
            Retry:      &zzz.RetryOptions{MaxAttempts: 2, Delay: 2 * time.Second},
            // Cache: myCacheInstance, // plug in your own zzz.Cache implementation
        }),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    profile, err := client.GetProfile(ctx, "1504687050")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Player: %s\n", profile.Nickname)
    for _, agent := range profile.Agents {
        fmt.Printf("  • %s Lv.%-2d [%s / %s]\n",
            agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName)
    }
}
```

---

### Error handling

Fairy exposes sentinel errors for common API failure scenarios.

| Error | Description |
| :--- | :--- |
| `fairy.ErrProfileNotFound` | The requested profile with the UID does not exist. |
| `fairy.ErrRateLimit` | EnkaNetwork API rate limit exceeded. |
| `fairy.ErrMaintenance` | API or game servers are under maintenance or temporarily unavailable. |
| `fairy.ErrNetwork` | Network-level error (timeout, DNS, etc.). |

---

### Feature highlights

#### Agent stats

Fairy provides three distinct ways to work with agent stats:

| Representation | Type | Format | Use case |
| :--- | :--- | :--- | :--- |
| **`agent.Stats`** | [`Stats`](https://pkg.go.dev/github.com/kirinyoku/fairy#Stats) | Raw `float64` values (`CritRate: 0.05`, `HP: 11188`) | Calculations, damage formulas, stat comparisons |
| **`agent.Stats.Formatted()`** | [`FormattedStats`](https://pkg.go.dev/github.com/kirinyoku/fairy#FormattedStats) | Formatted strings (`CritRate: "5.0%"`, `HP: "11188"`) | Simple text views, logging, summaries |
| **`agent.UIStats`** | [`UIStats`](https://pkg.go.dev/github.com/kirinyoku/fairy#UIStats) | `Base` + `Added` = `Total` breakdown with localized names & SVG icons | Rich character sheet UIs matching the in-game stats panel |

📖 See [`ExampleUIStats_List`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-UIStats_List)

#### Agent Skills

Agents have two views of their abilities:

- [`agent.Skills`](https://pkg.go.dev/github.com/kirinyoku/fairy#Agent.Skills) — A flat list of every individual ability, core enhancement, and passive on the agent.
- [`agent.GroupedSkills()`](https://pkg.go.dev/github.com/kirinyoku/fairy#Agent.GroupedSkills) — Categorized into the 6 in-game UI tabs (`basic`, `special`, `dodge`, `chain`, `assist`, `passive`), tracking each tab's upgrade level (`Level`, 1–12 for active skills, 0–6 for core passives) with its nested skills.

📖 See [`ExampleAgent_GroupedSkills`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-Agent_GroupedSkills)

#### Drive Disc analysis

- [`agent.SubStatTotals()`](https://pkg.go.dev/github.com/kirinyoku/fairy#Agent.SubStatTotals) aggregates sub-stat values and roll counts across all 6 Drive Discs. 
- [`agent.CountEffectiveRolls(props ...PropertyID)`](https://pkg.go.dev/github.com/kirinyoku/fairy#Agent.CountEffectiveRolls) counts how many rolls landed on the properties you care about — useful for evaluating build quality.

📖 See [`ExampleAgent_SubStatTotals`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-Agent_SubStatTotals), [`ExampleAgent_CountEffectiveRolls`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-Agent_CountEffectiveRolls)

#### Rich text formatting

Skill descriptions, Mindscape effects, Potential Vision effects, W-Engine passives, and Drive Disc set bonuses contain Unity Rich Text markup. Built-in formatters convert the markup and evaluate level-scaling formulas, outputting your target format:

| Method | Output |
| :--- | :--- |
| `FormatHTML()` | `<span style="color:...">` + `<img>` icons + `<br>` |
| `FormatMarkdown()` | `**bold**` highlights + text icon labels |
| `FormatPlainText()` | Clean text, all tags stripped |

These methods are available on [`Skill`](https://pkg.go.dev/github.com/kirinyoku/fairy#Skill), [`MindscapeNode`](https://pkg.go.dev/github.com/kirinyoku/fairy#MindscapeNode), [`PotentialVisionNode`](https://pkg.go.dev/github.com/kirinyoku/fairy#PotentialVisionNode), [`WEngine`](https://pkg.go.dev/github.com/kirinyoku/fairy#WEngine), and [`SetEffect`](https://pkg.go.dev/github.com/kirinyoku/fairy#SetEffect).

📖 See [`ExampleSkill_FormatHTML`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-Skill_FormatHTML), [`ExampleSkill_FormatMarkdown`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-Skill_FormatMarkdown), [`ExampleSkill_FormatPlainText`](https://pkg.go.dev/github.com/kirinyoku/fairy#example-Skill_FormatPlainText).

## Supported Languages

Fairy supports all 13 official in-game languages. All localization data is embedded directly into the binary — no external files or network requests are needed at runtime.

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

## License

Fairy is released under the [MIT License](LICENSE).