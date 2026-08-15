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

func TestAllRegions(t *testing.T) {
	regions := AllRegions()
	if len(regions) != 4 {
		t.Fatalf("AllRegions() returned %d regions, want 4", len(regions))
	}

	// Verify mutability protection (copy is returned)
	regions[0] = Region("corrupted")
	freshRegions := AllRegions()
	if freshRegions[0] == Region("corrupted") {
		t.Errorf("AllRegions() is vulnerable to external slice mutation")
	}
}

func TestRegion_IsValid(t *testing.T) {
	tests := []struct {
		region Region
		valid  bool
	}{
		{RegionEU, true},
		{RegionNA, true},
		{RegionAsia, true},
		{RegionTWHKMO, true},
		{Region("Mars"), false},
		{Region(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.region), func(t *testing.T) {
			if got := tt.region.IsValid(); got != tt.valid {
				t.Errorf("Region(%q).IsValid() = %v, want %v", tt.region, got, tt.valid)
			}
		})
	}
}
