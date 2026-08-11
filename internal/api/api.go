// Package api provides a thin adapter layer over the upstream enkanetwork-go client.
// It exists to isolate the rest of the application from changes in the upstream
// library and to provide domain-specific error mapping.
package api

import (
	"context"
	"errors"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
)

// Fetcher defines the interface for fetching raw data from the EnkaNetwork API.
// This interface primarily exists for testability, allowing mocks to be used
// during unit tests.
type Fetcher interface {
	GetProfile(ctx context.Context, uid string) (*zzz.Profile, error)
}

// Client wraps the enkanetwork-go client to fetch raw profile data.
type Client struct {
	enka *zzz.Client
}

// NewClient creates a new API client using the provided enkanetwork options.
func NewClient(opts zzz.Options) *Client {
	return &Client{
		enka: zzz.New(opts),
	}
}

// GetProfile fetches the raw player profile by UID.
func (c *Client) GetProfile(ctx context.Context, uid string) (*zzz.Profile, error) {
	rawProfile, err := c.enka.GetProfile(ctx, uid)
	if err != nil {
		return nil, mapEnkaError(err)
	}
	return rawProfile, nil
}

var (
	// ErrProfileNotFound is returned when the player profile is not found.
	ErrProfileNotFound = errors.New("profile not found")
	// ErrRateLimit is returned when the player is rate limited by the EnkaNetwork API.
	ErrRateLimit = errors.New("rate limited by enkanetwork API")
	// ErrMaintenance is returned when the EnkaNetwork API is under maintenance or temporarily unavailable.
	ErrMaintenance = errors.New("enkanetwork API is under maintenance or temporarily unavailable")
	// ErrNetwork is returned when there is a network error while contacting the EnkaNetwork API.
	ErrNetwork = errors.New("network error while contacting enkanetwork API")
)

func mapEnkaError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, zzz.ErrPlayerNotFound) {
		return errors.Join(ErrProfileNotFound, err)
	}
	if errors.Is(err, zzz.ErrRateLimited) {
		return errors.Join(ErrRateLimit, err)
	}
	if errors.Is(err, zzz.ErrServerMaintenance) || errors.Is(err, zzz.ErrServiceUnavailable) || errors.Is(err, zzz.ErrServerError) {
		return errors.Join(ErrMaintenance, err)
	}

	return errors.Join(ErrNetwork, err)
}
