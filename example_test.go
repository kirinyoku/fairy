package fairy_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/fairy"
)

// ExampleGetProfile demonstrates fetching a player's showcase profile with context timeout
// and handling domain sentinel errors (e.g. player not found, rate limit).
func ExampleGetProfile() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const uid = "1504687050"

	profile, err := fairy.GetProfile(ctx, uid)
	if err != nil {
		switch {
		case errors.Is(err, fairy.ErrProfileNotFound):
			log.Printf("Player %s not found\n", uid)
		case errors.Is(err, fairy.ErrRateLimit):
			log.Printf("EnkaNetwork API rate limit reached, back off required\n")
		case errors.Is(err, fairy.ErrMaintenance):
			log.Printf("EnkaNetwork API or game servers are under maintenance or temporarily unavailable\n")
		default:
			log.Printf("Failed to fetch %s profile: %v\n", uid, err)
		}
		return
	}

	fmt.Printf("Player: %s (Inter-Knot Lv.%d, Server: %s)\n",
		profile.Nickname, profile.InterknotLevel, profile.Region)

	if profile.Title != nil {
		fmt.Printf("Title: %s\n", profile.Title.Text)
	}

	for _, agent := range profile.Agents {
		weapon := "None"
		if agent.WEngine != nil {
			weapon = fmt.Sprintf("%s (Lv.%d, M%d)", agent.WEngine.Name, agent.WEngine.Level, agent.WEngine.Modification)
		}
		fmt.Printf("  • %-16s Lv.%-2d [%s / %s]  W-Engine: %s\n",
			agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName, weapon)
	}
}

// ExampleGetProfileWithLang demonstrates fetching a profile in a specific language (e.g. Japanese).
func ExampleGetProfileWithLang() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfileWithLang(ctx, "1504687050", fairy.LangJA)
	if err != nil {
		log.Printf("Failed to fetch profile: %v\n", err)
		return
	}

	fmt.Printf("Player: %s\n", profile.Nickname)
	for _, agent := range profile.Agents {
		fmt.Printf("  • %s Lv.%-2d [%s / %s]\n",
			agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName)
	}
}

// ExampleNewClient demonstrates configuring a custom client instance for a production service
// with a custom User-Agent, HTTP timeout, retry policy, caching, and default language.
func ExampleNewClient() {
	// Optional: initialize a Redis or in-memory cache instance
	// redisCache := NewRedisCache(rdb)

	client, err := fairy.NewClient(
		fairy.WithDefaultLang(fairy.LangJA),
		fairy.WithEnkaOptions(zzz.Options{
			UserAgent: "MyApp/v1.0.0 (contact@example.com)",
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Retry: &zzz.RetryOptions{
				MaxAttempts: 3,
				Delay:       1 * time.Second,
			},
			// Cache: redisCache, // Plug in Redis cache
		}),
	)
	if err != nil {
		log.Fatalf("Failed to initialize fairy client: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := client.GetProfile(ctx, "1504687050")
	if err != nil {
		log.Printf("Profile request failed: %v\n", err)
		return
	}

	fmt.Printf("Player: %s\n", profile.Nickname)
	for _, agent := range profile.Agents {
		fmt.Printf("  • %s Lv.%-2d [%s / %s]\n",
			agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName)
	}
}

// ExampleLocalize demonstrates fetching raw API data once and localizing it into multiple
// languages in memory without making extra network requests.
func ExampleLocalize() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Fetch the raw upstream profile once over HTTP
	rawProfile, err := fairy.GetRawProfile(ctx, "1504687050")
	if err != nil {
		log.Printf("Failed to fetch raw profile: %v", err)
		return
	}

	// 2. Localize into German and Japanese in-memory (zero additional network calls)
	enProfile, err := fairy.Localize(rawProfile, fairy.LangDE)
	if err != nil {
		log.Printf("German localization failed: %v", err)
		return
	}

	jaProfile, err := fairy.Localize(rawProfile, fairy.LangJA)
	if err != nil {
		log.Printf("Japanese localization failed: %v", err)
		return
	}

	fmt.Printf("German:\n")
	for _, agent := range enProfile.Agents {
		fmt.Printf("  • %s Lv.%-2d [%s / %s]\n",
			agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName)
	}

	fmt.Printf("Japanese:\n")
	for _, agent := range jaProfile.Agents {
		fmt.Printf("  • %s Lv.%-2d [%s / %s]\n",
			agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName)
	}
}

// ExampleUIStats_List demonstrates rendering an agent's combat stats panel.
func ExampleUIStats_List() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil {
		log.Printf("Failed to fetch profile: %v\n", err)
		return
	}

	if len(profile.Agents) == 0 {
		return
	}

	agent := profile.Agents[0]
	fmt.Printf("=== Agent Stats: %s Lv.%d ===\n", agent.Name, agent.Level)

	for _, stat := range agent.UIStats.List() {
		fmt.Printf("%-24s %10s  (%s + %s)\n",
			stat.Name, stat.Total, stat.Base, stat.Added)
	}
}

// ExampleAgent_SubStatTotals demonstrates analyzing an agent's Drive Discs to summarize
// total sub-stat rolls and aggregate values across all 6 disc slots.
func ExampleAgent_SubStatTotals() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil {
		log.Printf("Failed to fetch profile: %v\n", err)
		return
	}

	if len(profile.Agents) == 0 {
		return
	}

	agent := profile.Agents[0]
	fmt.Printf("=== Drive Disc Substat Summary: %s ===\n", agent.Name)

	for _, stat := range agent.SubStatTotals() {
		fmt.Printf("  • %-24s +%-8s  (%2d rolls)\n",
			stat.Name, stat.DisplayValue(), stat.Rolls)
	}
}

// ExampleAgent_CountEffectiveRolls demonstrates evaluating build quality by counting how many
// sub-stat upgrade rolls landed on target priority stats for a specific combat role.
func ExampleAgent_CountEffectiveRolls() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil {
		log.Printf("Failed to fetch profile: %v", err)
		return
	}

	if len(profile.Agents) == 0 {
		return
	}

	for _, agent := range profile.Agents {
		if agent.Specialty == fairy.SpecialtyAttack {
			fmt.Printf("Agent %s has %d effective substat rolls on equipped Drive Discs\n",
				agent.Name, agent.CountEffectiveRolls(fairy.PropCritRate, fairy.PropCritDMG, fairy.PropATKPercent))
		}
	}
}

// ExampleAgent_GroupedSkills demonstrates organizing an agent's abilities by in-game UI tabs
// (Basic, Special, Dodge, Chain, Assist, Passives) with group progression levels.
func ExampleAgent_GroupedSkills() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil {
		log.Printf("Failed to fetch profile: %v\n", err)
		return
	}

	if len(profile.Agents) == 0 {
		return
	}

	agent := profile.Agents[0]
	fmt.Printf("=== Skills for %s ===\n", agent.Name)

	for _, group := range agent.GroupedSkills() {
		fmt.Printf("\n[%s Tab] (Level %d)\n", group.TypeName, group.Level)
		for _, skill := range group.Skills {
			fmt.Printf("  • %s\n", skill.Name)
		}
	}
}

// ExampleSkill_FormatHTML demonstrates rendering a skill's description with colored text spans
// and embedded icon tags suitable for a web frontend.
func ExampleSkill_FormatHTML() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil || len(profile.Agents) == 0 || len(profile.Agents[0].Skills) == 0 {
		return
	}

	skill := profile.Agents[0].Skills[0]
	html := skill.FormatHTML()

	fmt.Printf("<h3>%s (Lv.%d)</h3>\n<div class=\"skill-desc\">\n%s\n</div>\n",
		skill.Name, skill.Level, html)
}

// ExampleSkill_FormatPlainText demonstrates stripping all markup and evaluating level scaling
// formulas for plain text environments like terminal output or logs.
func ExampleSkill_FormatPlainText() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil || len(profile.Agents) == 0 || len(profile.Agents[0].Skills) == 0 {
		return
	}

	skill := profile.Agents[0].Skills[0]
	plainText := skill.FormatPlainText()

	fmt.Printf("%s (Lv.%d):\n%s\n", skill.Name, skill.Level, plainText)
}

// ExampleSkill_FormatMarkdown demonstrates formatting a skill description with Markdown bold highlights.
func ExampleSkill_FormatMarkdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil || len(profile.Agents) == 0 || len(profile.Agents[0].Skills) == 0 {
		return
	}

	skill := profile.Agents[0].Skills[0]
	markdown := skill.FormatMarkdown()

	fmt.Printf("**%s** (Lv.%d)\n%s\n", skill.Name, skill.Level, markdown)
}
