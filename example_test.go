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

// ExampleGetProfile demonstrates the standard quick-start workflow: fetching a player's
// showcase profile with a context timeout and handling sentinel errors with errors.Is.
func ExampleGetProfile() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const uid = "1504687050"

	profile, err := fairy.GetProfile(ctx, uid)
	if err != nil {
		switch {
		case errors.Is(err, fairy.ErrProfileNotFound):
			log.Printf("Player %s not found on game servers\n", uid)
		case errors.Is(err, fairy.ErrRateLimit):
			log.Printf("EnkaNetwork API rate limit reached, back off required\n")
		case errors.Is(err, fairy.ErrMaintenance):
			log.Printf("EnkaNetwork API or game servers are under maintenance\n")
		case errors.Is(err, fairy.ErrNetwork):
			log.Printf("Network transport or timeout error\n")
		default:
			log.Printf("Unexpected error fetching profile %s: %v\n", uid, err)
		}
		return
	}

	fmt.Printf("Player: %s (Inter-Knot Lv.%d, Server: %s)\n",
		profile.Nickname, profile.InterknotLevel, profile.Region)

	for _, agent := range profile.Agents {
		weapon := "None"
		if agent.WEngine != nil {
			weapon = fmt.Sprintf("%s (Lv.%d, M%d)", agent.WEngine.Name, agent.WEngine.Level, agent.WEngine.Modification)
		}
		fmt.Printf("  • %-16s Lv.%-2d [%s / %s]  W-Engine: %s\n",
			agent.Name, agent.Level, agent.AttributeName, agent.SpecialtyName, weapon)
	}
}

// ExampleProfile demonstrates extracting and rendering player showcase customization
// data, such as avatar icons, namecard backgrounds, achievement badges, and two-color gradient titles.
func ExampleProfile() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil {
		log.Printf("Failed to fetch profile: %v\n", err)
		return
	}

	fmt.Printf("=== Player Showcase: %s (UID: %s) ===\n", profile.Nickname, profile.UID)
	fmt.Printf("Inter-Knot Level: %d | Server: %s\n", profile.InterknotLevel, profile.Region)

	if profile.Avatar != nil && profile.Avatar.URL != "" {
		fmt.Printf("Avatar URL: %s\n", profile.Avatar.URL)
	}

	if profile.Namecard != nil && profile.Namecard.URL != "" {
		fmt.Printf("Namecard URL: %s\n", profile.Namecard.URL)
	}

	if profile.Title != nil {
		// Generate an HTML gradient style using the Title's primary and secondary hex colors
		gradientCSS := fmt.Sprintf("background: linear-gradient(90deg, %s, %s); -webkit-background-clip: text;",
			profile.Title.PrimaryColorHex(), profile.Title.SecondaryColorHex())
		fmt.Printf("Title: <span style=\"%s\">%s</span>\n", gradientCSS, profile.Title.Text)
	}

	if len(profile.Badges) > 0 {
		fmt.Println("Showcased Badges:")
		for _, badge := range profile.Badges {
			fmt.Printf("  • %-24s Value: %-5d (Icon: %s)\n", badge.Title, badge.Value, badge.IconURL)
		}
	}
}

// ExampleNewClient demonstrates configuring a custom Client for production backend services
// with a custom User-Agent, HTTP client timeout, retry policy, caching layer, and default language.
func ExampleNewClient() {
	// Optional: initialize a custom persistent or in-memory cache instance
	// redisCache := NewRedisCache(rdb)

	client, err := fairy.NewClient(
		fairy.WithDefaultLang(fairy.LangJA),
		fairy.WithEnkaOptions(zzz.Options{
			UserAgent: "MyZZZApp/v1.0.0 (contact@example.com)",
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Retry: &zzz.RetryOptions{
				MaxAttempts: 3,
				Delay:       1 * time.Second,
			},
			// Cache: redisCache, // Plug in your cache implementation
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

// ExampleEnrich demonstrates fetching raw API data once with GetRawProfile, caching it,
// and enriching it into multiple languages in memory using Enrich and EnrichWithLang with zero additional network requests.
func ExampleEnrich() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Fetch the raw upstream profile once over HTTP
	rawProfile, err := fairy.GetRawProfile(ctx, "1504687050")
	if err != nil {
		log.Printf("Failed to fetch raw profile: %v", err)
		return
	}

	// 2. Enrich in memory with zero extra network overhead
	// Default English:
	enProfile, err := fairy.Enrich(rawProfile)
	if err != nil {
		log.Printf("English enrichment failed: %v", err)
		return
	}

	// Specific languages on the fly:
	jaProfile, err := fairy.EnrichWithLang(rawProfile, fairy.LangJA)
	if err != nil {
		log.Printf("Japanese enrichment failed: %v", err)
		return
	}

	ruProfile, err := fairy.EnrichWithLang(rawProfile, fairy.LangRU)
	if err != nil {
		log.Printf("Russian enrichment failed: %v", err)
		return
	}

	fmt.Println("English Agent 0:", enProfile.Agents[0].Name, "—", enProfile.Agents[0].SpecialtyName)
	fmt.Println("Japanese Agent 0:", jaProfile.Agents[0].Name, "—", jaProfile.Agents[0].SpecialtyName)
	fmt.Println("Russian Agent 0:", ruProfile.Agents[0].Name, "—", ruProfile.Agents[0].SpecialtyName)
}

// ExampleUIStats_List demonstrates rendering an Agent's combat stats panel with pre-formatted
// Base + Added = Total breakdowns and localized stat names matching the in-game attributes screen.
func ExampleUIStats_List() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil || len(profile.Agents) == 0 {
		return
	}

	agent := profile.Agents[0]
	fmt.Printf("=== Agent Stats: %s Lv.%d ===\n", agent.Name, agent.Level)

	for _, stat := range agent.UIStats.List() {
		fmt.Printf("%-24s %10s  (Base: %s + Added: %s)\n",
			stat.Name, stat.Total, stat.Base, stat.Added)
	}
}

// ExampleAgent_SkillGroups demonstrates organizing an Agent's abilities into the 6 in-game UI tabs
// (Basic, Special, Dodge, Chain, Assist, Passive) with tab upgrade levels and formatted descriptions.
func ExampleAgent_SkillGroups() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil || len(profile.Agents) == 0 {
		return
	}

	agent := profile.Agents[0]
	fmt.Printf("=== Skill Tabs for %s ===\n", agent.Name)

	for _, group := range agent.SkillGroups {
		fmt.Printf("\n[%s Tab] (Level %d)\n", group.TypeName, group.Level)
		for _, skill := range group.Skills {
			// FormatHTML converts Unity Rich Text color tags and button icons into web-ready HTML
			htmlDesc := skill.FormatHTML()
			fmt.Printf("  • %s\n    HTML: %s\n", skill.Name, htmlDesc)

			// Print evaluated scaling parameters (e.g. Damage Ratio: 142.5%)
			for _, param := range skill.Params {
				fmt.Printf("    - %s: %s\n", param.Name, param.Value)
			}
		}
	}
}

// ExampleAgent_Mindscapes demonstrates inspecting an Agent's Mindscape Cinema (Ranks 1–6)
// and Potential Vision passive tree progression to determine unlocked constellation states.
func ExampleAgent_Mindscapes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil || len(profile.Agents) == 0 {
		return
	}

	agent := profile.Agents[0]
	fmt.Printf("=== %s: Mindscape Cinema (Rank M%d) ===\n", agent.Name, agent.MindscapeCinema)

	for _, node := range agent.Mindscapes {
		status := "🔒 Locked"
		if node.Unlocked {
			status = "✨ Unlocked"
		}
		// FormatPlainText strips rich text tags for clean console or log display
		fmt.Printf("  M%d [%s] %s\n     %s\n",
			node.Rank, status, node.Name, node.FormatPlainText())
	}

	if agent.PotentialVision != nil && agent.PotentialVision.IsUnlocked {
		fmt.Println("\nPotential Vision Upgrades:")
		for _, pvNode := range agent.PotentialVision.Nodes {
			if pvNode.IsActive {
				fmt.Printf("  • [Active] %s: %s\n", pvNode.LevelName, pvNode.Title)
			}
		}
	}
}

// ExampleDriveDiscs_SubStatTotals demonstrates comprehensive gear and build quality analysis:
// detecting active 2-pc / 4-pc set bonuses, aggregating substats, and scoring priority rolls.
func ExampleDriveDiscs_SubStatTotals() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, err := fairy.GetProfile(ctx, "1504687050")
	if err != nil || len(profile.Agents) == 0 {
		return
	}

	agent := profile.Agents[0]
	fmt.Printf("=== Drive Disc Analysis: %s ===\n", agent.Name)

	// 1. Detect active Drive Disc Set Bonuses (2-piece and 4-piece thresholds)
	fmt.Println("Active Set Bonuses:")
	for _, bonus := range agent.DriveDiscs.SetBonuses {
		fmt.Printf("  • %s (%d pieces equipped)\n", bonus.Set.Name, bonus.Count)
		for _, effect := range bonus.Effects {
			if effect.IsActive {
				fmt.Printf("    [%d-Piece Active]: %s\n", effect.PieceCount, effect.FormatPlainText())
			}
		}
	}

	// 2. Aggregate substat totals and upgrade roll counts across all 6 equipped discs
	fmt.Println("\nSubstat Totals Across All 6 Discs:")
	for _, stat := range agent.DriveDiscs.SubStatTotals() {
		fmt.Printf("  • %-22s +%-8s (%d rolls)\n",
			stat.Name, stat.DisplayValue(), stat.Rolls)
	}

	// 3. Count effective upgrade rolls for the agent's priority stats (e.g. Crit Rate, Crit DMG, ATK%)
	if agent.Specialty == fairy.SpecialtyAttack {
		usefulRolls := agent.DriveDiscs.CountEffectiveRolls(
			fairy.PropCritRate,
			fairy.PropCritDMG,
			fairy.PropATKPercent,
		)
		fmt.Printf("\nBuild Rating: %d effective substat rolls on priority stats\n", usefulRolls)
	}
}
