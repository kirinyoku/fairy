package api

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
)

func TestMapEnkaError(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedErr error
	}{
		{
			name:        "nil error returns nil",
			inputErr:    nil,
			expectedErr: nil,
		},
		{
			name:        "ErrPlayerNotFound maps to ErrProfileNotFound",
			inputErr:    zzz.ErrPlayerNotFound,
			expectedErr: ErrProfileNotFound,
		},
		{
			name:        "wrapped ErrPlayerNotFound maps to ErrProfileNotFound",
			inputErr:    fmt.Errorf("api failure: %w", zzz.ErrPlayerNotFound),
			expectedErr: ErrProfileNotFound,
		},
		{
			name:        "ErrRateLimited maps to ErrRateLimit",
			inputErr:    zzz.ErrRateLimited,
			expectedErr: ErrRateLimit,
		},
		{
			name:        "ErrServerMaintenance maps to ErrMaintenance",
			inputErr:    zzz.ErrServerMaintenance,
			expectedErr: ErrMaintenance,
		},
		{
			name:        "ErrServiceUnavailable maps to ErrMaintenance",
			inputErr:    zzz.ErrServiceUnavailable,
			expectedErr: ErrMaintenance,
		},
		{
			name:        "ErrServerError maps to ErrMaintenance",
			inputErr:    zzz.ErrServerError,
			expectedErr: ErrMaintenance,
		},
		{
			name:        "unknown error maps to ErrNetwork",
			inputErr:    errors.New("connection refused"),
			expectedErr: ErrNetwork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapEnkaError(tt.inputErr)
			if tt.expectedErr == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", got)
				}
				return
			}

			if !errors.Is(got, tt.expectedErr) {
				t.Errorf("expected error matching %v, got %v", tt.expectedErr, got)
			}

			if tt.inputErr != nil && !errors.Is(got, tt.inputErr) {
				t.Errorf("expected mapped error to preserve original error %v, got %v", tt.inputErr, got)
			}
		})
	}
}
