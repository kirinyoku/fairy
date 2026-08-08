package fairy

import (
	"context"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/fairy/internal/api"
	"github.com/kirinyoku/fairy/store"
)

// Options holds the configuration for the fairy Client.
type Options struct {
	DefaultLang Language            // The default language for string localization.
	Store       store.MetadataStore // The store providing game metadata.
	EnkaOpts    zzz.Options         // Configuration for the underlying enkanetwork-go client.
}

// Option defines a functional option for the fairy Client.
type Option func(*Options)

// WithDefaultLang sets the default language for the Client.
func WithDefaultLang(lang Language) Option {
	return func(o *Options) {
		o.DefaultLang = lang
	}
}

// WithStore sets the metadata store for the Client.
func WithStore(s store.MetadataStore) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// WithEnkaOptions sets the underlying enkanetwork-go client options.
func WithEnkaOptions(enkaOpts zzz.Options) Option {
	return func(o *Options) {
		o.EnkaOpts = enkaOpts
	}
}

// Client coordinates fetching data from the API and enriching it using a metadata store.
// It serves as the main entry point for the fairy library.
type Client struct {
	apiClient *api.Client         // The enkanetwork-go client.
	store     store.MetadataStore // The data store used for localization mapping.
	lang      Language            // The default language for string localization.
	enkaOpts  zzz.Options         // The original options passed down to the enkanetwork-go client.
}

// NewClient creates a new instance of Client.
// If opts.Store is nil, it will automatically load the default EmbeddedStore.
// If opts.DefaultLang is empty, it defaults to LangEN.
// Returns an error if the fallback EmbeddedStore fails to load its internal files.
func NewClient(opts ...Option) (*Client, error) {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	lang := options.DefaultLang
	if lang == "" {
		lang = LangEN
	}

	st := options.Store
	if st == nil {
		emb, err := store.Default()
		if err != nil {
			return nil, err
		}
		st = emb
	}

	c := &Client{
		lang:      lang,
		store:     st,
		enkaOpts:  options.EnkaOpts,
		apiClient: api.NewClient(options.EnkaOpts),
	}

	return c, nil
}

// GetProfile fetches the game profile from the EnkaNetwork API and localizes it using the client's default language.
// The provided context controls the HTTP request timeout and cancellation.
func (c *Client) GetProfile(ctx context.Context, uid string) (*Profile, error) {
	return c.GetProfileWithLang(ctx, uid, c.lang)
}

// GetProfileWithLang fetches the game profile from the EnkaNetwork API and localizes it with a specific language.
// The provided context controls the HTTP request timeout and cancellation.
func (c *Client) GetProfileWithLang(ctx context.Context, uid string, lang Language) (*Profile, error) {
	raw, err := c.GetRawProfile(ctx, uid)
	if err != nil {
		return nil, err
	}

	return c.Localize(raw, lang)
}

// GetRawProfile fetches the raw zzz.Profile without localization.
// Use this if you only need the raw data structure provided by the EnkaNetwork API.
// The provided context controls the HTTP request timeout and cancellation.
func (c *Client) GetRawProfile(ctx context.Context, uid string) (*zzz.Profile, error) {
	return c.apiClient.GetProfile(ctx, uid)
}

// Localize maps a raw zzz.Profile into an enriched Profile using the specified language.
// This is highly useful when you want to fetch the raw profile once, but display it in
// multiple different languages.
func (c *Client) Localize(raw *zzz.Profile, lang Language) (*Profile, error) {
	m := newMapper(c.store, lang)
	p, err := m.ToProfile(raw)
	if err != nil {
		return nil, err
	}
	return p, nil
}
