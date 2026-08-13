package fairy

import (
	"testing"
)

func TestTitle_ColorHex(t *testing.T) {
	tests := []struct {
		name              string
		title             *Title
		expectedPrimary   string
		expectedSecondary string
	}{
		{
			name:              "nil title receiver",
			title:             nil,
			expectedPrimary:   "",
			expectedSecondary: "",
		},
		{
			name: "empty color strings",
			title: &Title{
				PrimaryColor:   "",
				SecondaryColor: "",
			},
			expectedPrimary:   "",
			expectedSecondary: "",
		},
		{
			name: "valid 6-character hex colors",
			title: &Title{
				PrimaryColor:   "FF5500",
				SecondaryColor: "00AAFF",
			},
			expectedPrimary:   "#FF5500",
			expectedSecondary: "#00AAFF",
		},
		{
			name: "only primary color set",
			title: &Title{
				PrimaryColor:   "FFFFFF",
				SecondaryColor: "",
			},
			expectedPrimary:   "#FFFFFF",
			expectedSecondary: "",
		},
		{
			name: "only secondary color set",
			title: &Title{
				PrimaryColor:   "",
				SecondaryColor: "123456",
			},
			expectedPrimary:   "",
			expectedSecondary: "#123456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrimary := tt.title.PrimaryColorHex()
			if gotPrimary != tt.expectedPrimary {
				t.Errorf("PrimaryColorHex() = %q, want %q", gotPrimary, tt.expectedPrimary)
			}

			gotSecondary := tt.title.SecondaryColorHex()
			if gotSecondary != tt.expectedSecondary {
				t.Errorf("SecondaryColorHex() = %q, want %q", gotSecondary, tt.expectedSecondary)
			}
		})
	}
}
