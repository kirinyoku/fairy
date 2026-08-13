package fairy

import (
	"testing"
)

func TestDriveDisc_CountEffectiveRolls(t *testing.T) {
	tests := []struct {
		name        string
		subStats    []StatValue
		targetProps []PropertyID
		expected    int
	}{
		{
			name: "single matching target",
			subStats: []StatValue{
				{PropertyID: PropCritRate, Rolls: 2, Value: 4.8},
				{PropertyID: PropCritDMG, Rolls: 1, Value: 4.8},
				{PropertyID: PropATKFlat, Rolls: 1, Value: 19},
				{PropertyID: PropHPFlat, Rolls: 1, Value: 112},
			},
			targetProps: []PropertyID{PropCritRate},
			expected:    2,
		},
		{
			name: "multiple matching targets",
			subStats: []StatValue{
				{PropertyID: PropCritRate, Rolls: 3, Value: 7.2},
				{PropertyID: PropCritDMG, Rolls: 2, Value: 9.6},
				{PropertyID: PropATKPercent, Rolls: 1, Value: 3.0},
				{PropertyID: PropDEFPercent, Rolls: 1, Value: 4.8},
			},
			targetProps: []PropertyID{PropCritRate, PropCritDMG, PropATKPercent},
			expected:    6,
		},
		{
			name: "no matching targets",
			subStats: []StatValue{
				{PropertyID: PropHPFlat, Rolls: 2, Value: 224},
				{PropertyID: PropDEFFlat, Rolls: 3, Value: 45},
			},
			targetProps: []PropertyID{PropCritRate, PropCritDMG},
			expected:    0,
		},
		{
			name: "empty targets",
			subStats: []StatValue{
				{PropertyID: PropCritRate, Rolls: 3, Value: 7.2},
			},
			targetProps: []PropertyID{},
			expected:    0,
		},
		{
			name:        "empty sub-stats",
			subStats:    []StatValue{},
			targetProps: []PropertyID{PropCritRate},
			expected:    0,
		},
		{
			name: "duplicate target props do not double-count",
			subStats: []StatValue{
				{PropertyID: PropCritRate, Rolls: 3, Value: 7.2},
				{PropertyID: PropCritDMG, Rolls: 2, Value: 9.6},
			},
			targetProps: []PropertyID{PropCritRate, PropCritRate, PropCritDMG},
			expected:    5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disc := &DriveDisc{
				SubStats: tt.subStats,
			}
			got := disc.CountEffectiveRolls(tt.targetProps...)
			if got != tt.expected {
				t.Errorf("CountEffectiveRolls() = %d, want %d", got, tt.expected)
			}
		})
	}
}
