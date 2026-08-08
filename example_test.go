package fairy_test

import (
	"context"
	"fmt"
	"log"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/fairy"
)

func ExampleGetProfile() {
	// 1. You can fetch a profile directly using the package-level helper method.
	// This will use a shared default client with English localization.
	profile, err := fairy.GetProfile(context.Background(), "100000000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Player: %s (UID: %s)\n", profile.Nickname, profile.UID)
	fmt.Printf("Level: %d | Region: %s\n", profile.InterknotLevel, profile.Region)

	if len(profile.Agents) > 0 {
		agent := profile.Agents[0]
		fmt.Printf("Featured Agent: %s (Level %d, %s)\n", agent.Name, agent.Level, agent.AttributeName)
	}
}

func ExampleLocalize() {
	// 1. Initialize a custom client if you want to set specific options
	// (e.g. setting up caching or changing the default language).
	client, err := fairy.NewClient(
		fairy.WithDefaultLang(fairy.LangRU), // Use Russian by default
		fairy.WithEnkaOptions(zzz.Options{
			UserAgent: "MyAwesomeApp/1.0",
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Fetch the profile. It will be localized in Russian automatically.
	profile, err := client.GetProfile(context.Background(), "100000000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Игрок: %s (Сервер: %s)\n", profile.Nickname, profile.Region)
}
