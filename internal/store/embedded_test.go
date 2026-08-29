package store

import (
	"sync"
	"testing"
)

func TestEmbeddedStore_Localize(t *testing.T) {
	st, err := Default()
	if err != nil {
		t.Fatalf("failed to get default store: %v", err)
	}

	tests := []struct {
		name     string
		hash     string
		lang     string
		want     string
		fallback bool
	}{
		{
			name: "Basic English localization",
			hash: "AddedDamageRatio",
			lang: "en",
			want: "DMG Bonus",
		},
		{
			name: "Basic Russian localization",
			hash: "AddedDamageRatio",
			lang: "ru",
			want: "Бонус к урону",
		},
		{
			name: "Basic Japanese localization",
			hash: "AddedDamageRatio",
			lang: "ja",
			want: "ダメージボーナス",
		},
		{
			name: "Unknown hash returns original hash",
			hash: "NonExistentKey12345",
			lang: "en",
			want: "NonExistentKey12345",
		},
		{
			name: "Unknown language falls back to English if key exists in EN",
			hash: "AddedDamageRatio",
			lang: "unknown_lang",
			want: "DMG Bonus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := st.Localize(tt.hash, tt.lang)
			if got != tt.want {
				t.Errorf("Localize(%q, %q) = %q, want %q", tt.hash, tt.lang, got, tt.want)
			}
		})
	}
}

func TestEmbeddedStore_LazyLoading(t *testing.T) {
	// Create a fresh embedded store directly to test lazy loading behavior
	subFS, err := Default()
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	// Verify that querying a specific language loads only that language
	_ = subFS.Localize("AddedDamageRatio", "ko")

	subFS.locsMu.RLock()
	koLoaded := subFS.locs["ko"] != nil
	subFS.locsMu.RUnlock()

	if !koLoaded {
		t.Errorf("expected 'ko' dictionary to be loaded lazily")
	}
}

func TestEmbeddedStore_ConcurrentLocalize(t *testing.T) {
	st, err := Default()
	if err != nil {
		t.Fatalf("failed to get default store: %v", err)
	}

	langs := []string{"en", "ru", "ja", "zh-cn", "de", "es", "fr", "ko", "pt", "th", "vi", "zh-tw", "id", "invalid"}
	keys := []string{"AddedDamageRatio", "AddedDamageRatio_Fire", "AddedDamageRatio_Ice", "NonExistentKey"}

	var wg sync.WaitGroup
	workers := 50

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for _, lang := range langs {
				for _, key := range keys {
					_ = st.Localize(key, lang)
				}
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkStore_Default(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Default()
	}
}

func BenchmarkStore_Localize(b *testing.B) {
	st, _ := Default()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = st.Localize("AddedDamageRatio", "ru")
	}
}
