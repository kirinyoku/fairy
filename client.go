package fairy

import (
	"context"
	"fmt"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/fairy/internal/api"
	"github.com/kirinyoku/fairy/internal/store"
)

// Options holds configuration settings for initializing a [Client].
type Options struct {
	// DefaultLang specifies the default localization language used for profile enrichment.
	// Defaults to [LangEN] if not specified.
	DefaultLang Language

	// EnkaOpts specifies configuration settings for the underlying EnkaNetwork HTTP client,
	// such as custom HTTP transport, User-Agent, retry policies, and cache implementations.
	EnkaOpts zzz.Options
}

// Option defines a functional option for configuring a [Client] in [NewClient].
type Option func(*Options)

// WithDefaultLang sets the default localization [Language] for the [Client].
// If not specified, the client defaults to English ([LangEN]).
func WithDefaultLang(lang Language) Option {
	return func(o *Options) {
		o.DefaultLang = lang
	}
}

// WithEnkaOptions configures the underlying enkanetwork-go client settings,
// including custom [http.Client] timeouts, User-Agent strings, retry policies ([zzz.RetryOptions]),
// and persistent caching implementations ([zzz.Cache]).
func WithEnkaOptions(enkaOpts zzz.Options) Option {
	return func(o *Options) {
		o.EnkaOpts = enkaOpts
	}
}

// Client coordinates fetching player profile data from the EnkaNetwork API
// and enriching it with localized names, CDN asset URLs, and combat stat calculations
// using an embedded metadata store.
//
// A Client is safe for concurrent use by multiple goroutines.
type Client struct {
	apiClient *api.Client         // The enkanetwork-go client.
	store     store.MetadataStore // The data store used for localization mapping.
	lang      Language            // The default language for string localization.
	enkaOpts  zzz.Options         // The original options passed down to the enkanetwork-go client.
}

// NewClient creates a new configured [Client] instance.
//
// If [WithDefaultLang] is omitted, the client defaults to English ([LangEN]).
// Returns an error if the embedded game metadata store fails to load its internal data files.
//
// Example:
//
//	client, err := fairy.NewClient(
//		fairy.WithDefaultLang(fairy.LangJA),
//		fairy.WithEnkaOptions(zzz.Options{
//			UserAgent:  "MyZZZApp/1.0 (contact@example.com)",
//			HTTPClient: &http.Client{Timeout: 10 * time.Second},
//		}),
//	)
func NewClient(opts ...Option) (*Client, error) {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	lang := options.DefaultLang
	if lang == "" {
		lang = LangEN
	}

	st, err := store.Default()
	if err != nil {
		return nil, err
	}

	c := &Client{
		lang:      lang,
		store:     st,
		enkaOpts:  options.EnkaOpts,
		apiClient: api.NewClient(options.EnkaOpts),
	}

	return c, nil
}

// GetProfile fetches a player profile by UID from the EnkaNetwork API
// and enriches it using the client's configured default [Language].
//
// The returned [Profile] contains:
//   - Player account details (UID, Nickname, Inter-Knot Level, Region, Title, Avatar, Badges).
//   - Showcase [Agent] list (up to 6 agents), where each agent contains:
//   - Metadata (localized name, Attribute, Specialty, Rarity, Skin, CDN asset URLs).
//   - Progression (categorized [SkillGroup] entries, [MindscapeNode] unlock states, [PotentialVision]).
//   - Equipped [WEngine] (with level-scaled stats and refinement modification effects).
//   - Equipped [DriveDisc] entries (slots 1–6 with roll counts) and active [DriveDiscSetBonus] thresholds.
//   - Pre-calculated combat [Stats] and frontend-ready [UIStats] breakdowns.
//
// The provided [context.Context] controls the HTTP request lifecycle, cancellation, and timeout.
// Returns sentinel errors such as [ErrInvalidUID], [ErrProfileNotFound], [ErrRateLimit], [ErrMaintenance], [ErrNetwork], or [ErrEnrichment].
func (c *Client) GetProfile(ctx context.Context, uid string) (*Profile, error) {
	return c.GetProfileWithLang(ctx, uid, c.lang)
}

// GetProfileWithLang fetches a player profile by UID from the EnkaNetwork API
// and enriches it using the specified [Language] localization.
//
// This method overrides the client's default language for the single request
// without modifying the client instance, making it safe for concurrent multi-language usage.
//
// The provided [context.Context] controls the HTTP request lifecycle, cancellation, and timeout.
// Returns sentinel errors such as [ErrInvalidUID], [ErrProfileNotFound], [ErrRateLimit], [ErrMaintenance], [ErrNetwork], or [ErrEnrichment].
func (c *Client) GetProfileWithLang(ctx context.Context, uid string, lang Language) (*Profile, error) {
	raw, err := c.GetRawProfile(ctx, uid)
	if err != nil {
		return nil, err
	}

	return c.EnrichWithLang(raw, lang)
}

// GetRawProfile fetches the raw, un-enriched [zzz.Profile] directly from the EnkaNetwork API.
//
// Use this method if you only need the raw numeric IDs provided by the upstream API,
// or when you want to fetch the payload once and enrich it into multiple languages via [Client.EnrichWithLang].
//
// The provided [context.Context] controls the HTTP request lifecycle, cancellation, and timeout.
// Returns sentinel errors such as [ErrInvalidUID], [ErrProfileNotFound], [ErrRateLimit], [ErrMaintenance], or [ErrNetwork].
func (c *Client) GetRawProfile(ctx context.Context, uid string) (*zzz.Profile, error) {
	if err := validateUID(uid); err != nil {
		return nil, err
	}
	return c.apiClient.GetProfile(ctx, uid)
}

// Enrich transforms a raw upstream [zzz.Profile] into an enriched [Profile]
// using the client's configured default [Language].
//
// This method operates entirely in-memory using the client's metadata store and performs ZERO network calls.
// It resolves all progression data, computes scaled combat stats, parses Unity Rich Text into HTML,
// and assembles the full domain model.
// Returns [ErrEnrichment] if the raw profile payload is nil or corrupt.
func (c *Client) Enrich(raw *zzz.Profile) (*Profile, error) {
	return c.EnrichWithLang(raw, c.lang)
}

// EnrichWithLang transforms a raw upstream [zzz.Profile] into an enriched [Profile]
// using the specified [Language] localization.
//
// This method operates entirely in-memory using the client's metadata store and performs ZERO network calls.
// It is ideal for multi-language applications that fetch a player's raw profile once via [Client.GetRawProfile]
// and render it dynamically across different languages.
// Returns [ErrEnrichment] if the raw profile payload is nil or corrupt.
func (c *Client) EnrichWithLang(raw *zzz.Profile, lang Language) (*Profile, error) {
	m := newMapper(c.store, lang)
	p, err := m.ToProfile(raw)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// EnrichAgent transforms a raw upstream [zzz.AvatarData] into an enriched [Agent]
// using the client's configured default [Language].
//
// This method operates entirely in-memory using the client's metadata store and performs ZERO network calls.
// It resolves all agent metadata, progression data, equipped gear ([WEngine], [DriveDisc] entries),
// active set bonuses, computes scaled combat stats, and formats UI stats.
// Returns [ErrEnrichment] if the raw avatar payload is nil or refers to an unknown avatar ID.
func (c *Client) EnrichAgent(raw *zzz.AvatarData) (*Agent, error) {
	return c.EnrichAgentWithLang(raw, c.lang)
}

// EnrichAgentWithLang transforms a raw upstream [zzz.AvatarData] into an enriched [Agent]
// using the specified [Language] localization.
//
// This method operates entirely in-memory using the client's metadata store and performs ZERO network calls.
// It is ideal for multi-language applications that fetch player data once and render individual agents
// dynamically across different languages.
// Returns [ErrEnrichment] if the raw avatar payload is nil or refers to an unknown avatar ID.
func (c *Client) EnrichAgentWithLang(raw *zzz.AvatarData, lang Language) (*Agent, error) {
	if raw == nil {
		return nil, fmt.Errorf("%w: raw avatar data is nil", ErrEnrichment)
	}

	m := newMapper(c.store, lang)
	agent := m.ToAgent(raw)
	if agent == nil {
		return nil, fmt.Errorf("%w: unknown avatar id %d", ErrEnrichment, raw.ID)
	}
	return agent, nil
}

// AgentHighlightProps returns the recommended combat property IDs ([PropertyID]) for an Agent by their numeric ID,
// or nil if the agent is not found.
// The returned slice is a defensive copy and can be safely mutated by the caller.
func (c *Client) AgentHighlightProps(agentID int) []PropertyID {
	if c == nil || c.store == nil {
		return nil
	}
	meta, ok := c.store.AvatarMeta(agentID)
	if !ok || len(meta.HighlightProps) == 0 {
		return nil
	}
	props := make([]PropertyID, len(meta.HighlightProps))
	for i, p := range meta.HighlightProps {
		props[i] = PropertyID(p)
	}
	return props
}

// AllAgentHighlightProps returns a map of all agents' recommended combat property IDs indexed by agent ID.
// This is useful for code generation, static tooling, or caching build recommendation tables.
func (c *Client) AllAgentHighlightProps() map[int][]PropertyID {
	if c == nil || c.store == nil {
		return nil
	}
	allMetas := c.store.AllAvatarMetas()
	res := make(map[int][]PropertyID, len(allMetas))
	for id, meta := range allMetas {
		props := make([]PropertyID, len(meta.HighlightProps))
		for i, p := range meta.HighlightProps {
			props[i] = PropertyID(p)
		}
		res[id] = props
	}
	return res
}

// AgentRecommendedSubStats returns the recommended Drive Disc sub-stat property IDs ([PropertyID]) for an Agent by numeric ID.
// These values correspond directly to the official in-game Drive Disc recommendation system (yellow highlights in equipment UI).
func (c *Client) AgentRecommendedSubStats(agentID int) []PropertyID {
	if curated, exists := agentRecommendedSubStats[agentID]; exists {
		res := make([]PropertyID, len(curated))
		copy(res, curated)
		return res
	}

	hl := c.AgentHighlightProps(agentID)
	if len(hl) == 0 {
		return []PropertyID{}
	}

	var res []PropertyID
	for _, prop := range hl {
		switch prop {
		case PropBaseATK:
			res = append(res, PropATKPercent)
		case PropBaseHP:
			res = append(res, PropHPPercent)
		case PropBaseDEF:
			res = append(res, PropDEFPercent)
		case PropBaseCritRate:
			res = append(res, PropCritRate)
		case PropBaseCritDMG:
			res = append(res, PropCritDMG)
		case PropBaseAnomalyProficiency:
			res = append(res, PropAnomalyProficiency)
		}
	}
	if res == nil {
		return []PropertyID{}
	}
	return res
}

// AllAgentRecommendedSubStats returns a map of all known agents to their recommended Drive Disc sub-stats,
// matching the in-game equipment UI sub-stat highlights.
func (c *Client) AllAgentRecommendedSubStats() map[int][]PropertyID {
	allAvatars := c.AllAgentHighlightProps()
	res := make(map[int][]PropertyID, len(allAvatars))
	for id := range allAvatars {
		res[id] = c.AgentRecommendedSubStats(id)
	}
	return res
}
