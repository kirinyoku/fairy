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
// NOTE: This interface primarily exists for testability, allowing us to mock
// network calls during unit tests.
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

// Domain-specific error sentinels.
// NOTE: These are wrapped with upstream errors using errors.Join, meaning that
// clients can check for them using errors.Is(err, api.ErrProfileNotFound) while
// still having access to the original underlying error string.
var (
	ErrProfileNotFound = errors.New("profile not found or hidden")
	ErrRateLimit       = errors.New("rate limited by enkanetwork API")
	ErrMaintenance     = errors.New("enkanetwork API is under maintenance or temporarily unavailable")
	ErrNetwork         = errors.New("network error while contacting enkanetwork API")
)

// mapEnkaError translates raw HTTP/Enka errors into domain-specific errors.
// NOTE: We map typed errors exposed by the upstream enkanetwork-go client
// to our own domain errors. Previously, this might have used brittle string matching
// on the error message, but now we use standard errors.Is() checks.
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
