package fairy

import (
	"github.com/kirinyoku/fairy/internal/api"
)

// Sentinel errors returned by [GetProfile], [GetProfileWithLang], [GetRawProfile], and [*Client] methods.
// Callers should inspect errors using standard [errors.Is] checks.
//
// Example:
//
//	profile, err := fairy.GetProfile(ctx, uid)
//	if err != nil {
//		switch {
//		case errors.Is(err, fairy.ErrInvalidUID):
//			// Player UID format is invalid (must be 10 digits starting with 10, 13, 15, or 17)
//		case errors.Is(err, fairy.ErrProfileNotFound):
//			// Player profile does not exist on game servers
//		case errors.Is(err, fairy.ErrRateLimit):
//			// Rate limit exceeded (HTTP 429) — wait before retrying
//		case errors.Is(err, fairy.ErrMaintenance):
//			// Upstream API or game servers are under maintenance
//		case errors.Is(err, fairy.ErrNetwork):
//			// Network connection reset or request timeout
//		default:
//			// Other unexpected error
//		}
//	}
var (
	// ErrInvalidUID is returned when a provided player UID has an invalid format
	// (e.g. empty, non-numeric characters, length other than 10 digits, or unrecognized server prefix).
	ErrInvalidUID = api.ErrInvalidUID

	// ErrProfileNotFound is returned when the requested player profile with specified UID does not exist
	// or cannot be found by the upstream EnkaNetwork API (HTTP 404).
	ErrProfileNotFound = api.ErrProfileNotFound

	// ErrRateLimit is returned when requests exceed the EnkaNetwork API rate limit (HTTP 429 Too Many Requests).
	// Callers should apply exponential backoff or use a caching layer before retrying.
	ErrRateLimit = api.ErrRateLimit

	// ErrMaintenance is returned when the EnkaNetwork API or upstream game servers are undergoing maintenance,
	// experiencing an outage, or temporarily unavailable (HTTP 500, 502, 503, or 504).
	ErrMaintenance = api.ErrMaintenance

	// ErrNetwork is returned when a transport-level network error occurs while communicating with the API
	// (e.g. DNS resolution failure, connection refused, TLS handshake failure, or context timeout).
	ErrNetwork = api.ErrNetwork
)
