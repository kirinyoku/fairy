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

func TestIsValidUID(t *testing.T) {
	tests := []struct {
		name    string
		uid     string
		wantErr bool
	}{
		{name: "Valid 10-digit NA UID", uid: "1004687050", wantErr: false},
		{name: "Valid 10-digit Asia UID", uid: "1304687050", wantErr: false},
		{name: "Valid 10-digit EU UID", uid: "1504687050", wantErr: false},
		{name: "Valid 10-digit TW/HK/MO UID", uid: "1704687050", wantErr: false},
		{name: "Invalid 10-digit with unknown prefix 19", uid: "1904687050", wantErr: true},
		{name: "Invalid 10-digit with prefix 12", uid: "1204687050", wantErr: true},
		{name: "9-digit UID", uid: "100000001", wantErr: true},
		{name: "8-digit China UID", uid: "12345678", wantErr: true},
		{name: "Empty string", uid: "", wantErr: true},
		{name: "Too short (7 digits)", uid: "1234567", wantErr: true},
		{name: "Too long (11 digits)", uid: "12345678901", wantErr: true},
		{name: "Non-numeric characters", uid: "150468705a", wantErr: true},
		{name: "Special characters", uid: "150468-705", wantErr: true},
		{name: "Spaces included", uid: " 150468705", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUID(tt.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUID(%q) error = %v, wantErr %v", tt.uid, err, tt.wantErr)
			}
			if tt.wantErr && err != ErrInvalidUID {
				t.Errorf("validateUID(%q) returned %v, want %v", tt.uid, err, ErrInvalidUID)
			}
			if got := IsValidUID(tt.uid); got != !tt.wantErr {
				t.Errorf("IsValidUID(%q) = %v, want %v", tt.uid, got, !tt.wantErr)
			}
		})
	}
}

func TestRegionFromUID(t *testing.T) {
	tests := []struct {
		name       string
		uid        string
		wantRegion Region
		wantOk     bool
	}{
		{name: "America prefix 10", uid: "1004687050", wantRegion: RegionNA, wantOk: true},
		{name: "Asia prefix 13", uid: "1304687050", wantRegion: RegionAsia, wantOk: true},
		{name: "Europe prefix 15", uid: "1504687050", wantRegion: RegionEU, wantOk: true},
		{name: "TW/HK/MO prefix 17", uid: "1704687050", wantRegion: RegionTWHKMO, wantOk: true},
		{name: "Unknown prefix 19", uid: "1904687050", wantRegion: "", wantOk: false},
		{name: "8-digit China UID", uid: "12345678", wantRegion: "", wantOk: false},
		{name: "Short string", uid: "15", wantRegion: "", wantOk: false},
		{name: "Empty string", uid: "", wantRegion: "", wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg, ok := RegionFromUID(tt.uid)
			if ok != tt.wantOk {
				t.Errorf("RegionFromUID(%q) ok = %v, want %v", tt.uid, ok, tt.wantOk)
			}
			if reg != tt.wantRegion {
				t.Errorf("RegionFromUID(%q) region = %q, want %q", tt.uid, reg, tt.wantRegion)
			}
		})
	}
}
