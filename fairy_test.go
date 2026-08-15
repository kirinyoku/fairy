package fairy

import (
	"testing"
)

func TestAllLanguages(t *testing.T) {
	langs := AllLanguages()
	if len(langs) != 13 {
		t.Fatalf("AllLanguages() returned %d languages, want 13", len(langs))
	}

	// Verify mutability protection (copy is returned)
	langs[0] = Language("corrupted")
	freshLangs := AllLanguages()
	if freshLangs[0] == Language("corrupted") {
		t.Errorf("AllLanguages() is vulnerable to external slice mutation")
	}
}

func TestLanguage_IsValid(t *testing.T) {
	tests := []struct {
		lang  Language
		valid bool
	}{
		{LangEN, true},
		{LangRU, true},
		{LangDE, true},
		{LangES, true},
		{LangFR, true},
		{LangID, true},
		{LangJA, true},
		{LangKO, true},
		{LangPT, true},
		{LangTH, true},
		{LangVI, true},
		{LangZHCN, true},
		{LangZHTW, true},
		{Language("invalid"), false},
		{Language(""), false},
		{Language("EN"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.lang), func(t *testing.T) {
			if got := tt.lang.IsValid(); got != tt.valid {
				t.Errorf("Language(%q).IsValid() = %v, want %v", tt.lang, got, tt.valid)
			}
		})
	}
}
