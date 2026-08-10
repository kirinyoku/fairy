package fairy

import (
	"github.com/kirinyoku/fairy/internal/api"
)

var (
	// ErrProfileNotFound is returned when the requested UID does not exist.
	ErrProfileNotFound = api.ErrProfileNotFound
	// ErrRateLimit is returned when the client exceeds the API rate limit.
	ErrRateLimit = api.ErrRateLimit
	// ErrMaintenance is returned when the API is undergoing maintenance or temporarily unavailable.
	ErrMaintenance = api.ErrMaintenance
	// ErrNetwork is returned when there is a network error.
	ErrNetwork = api.ErrNetwork
)
