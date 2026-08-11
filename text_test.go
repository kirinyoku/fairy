package fairy_test

import (
	"testing"

	"github.com/kirinyoku/fairy"
)

func TestEvaluateFormulas(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		skillLevel int
		expected   string
	}{
		{
			name:       "simple CAL formula at Lv 12",
			input:      "Deals {CAL:0.07+AvatarSkillLevel(1)*0.015,100,2}% Ice DMG.",
			skillLevel: 12,
			expected:   "Deals 25% Ice DMG.",
		},
		{
			name:       "scaling formula with base ATK addition",
			input:      "DMG increases by {CAL:100+AvatarSkillLevel(3)*60,1,2} points.",
			skillLevel: 12,
			expected:   "DMG increases by 820 points.",
		},
		{
			name:       "CAL tag without level",
			input:      "Static value {CAL:2*20*0.8,1,2}%",
			skillLevel: 1,
			expected:   "Static value 32%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fairy.EvaluateFormulas(tt.input, tt.skillLevel)
			if result != tt.expected {
				t.Errorf("EvaluateFormulas() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatHTML(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		skillLevel []int
		expected   string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "plain text without tags",
			input:    "Increases ATK by 20%.",
			expected: "Increases ATK by 20%.",
		},
		{
			name:     "unity color tag",
			input:    "Increases ATK by <color=#2BAD00>20%</color>.",
			expected: `Increases ATK by <span style="color: #2BAD00;">20%</span>.`,
		},
		{
			name:       "unity color tag with CAL formula evaluation",
			input:      "Deals <color=#98EFF0>{CAL:0.07+AvatarSkillLevel(1)*0.015,100,2}%</color> Ice DMG.",
			skillLevel: []int{12},
			expected:   `Deals <span style="color: #98EFF0;">25%</span> Ice DMG.`,
		},
		{
			name:     "iconmap tag",
			input:    "Hold <IconMap:Icon_SpecialReady> to activate",
			expected: `Hold <img src="https://enka.network/ui/zzz/IconRoleSkillKeySpecialV2.png" class="fairy-icon" alt="[EX Special Attack]" style="height: 1.2em; vertical-align: -0.2em; display: inline-block;" /> to activate`,
		},
		{
			name:     "layout tag resolution russian",
			input:    "Для активации потяните {LAYOUT_XBOXCONTROLLER#мини-джойстик}",
			expected: "Для активации потяните мини-джойстик",
		},
		{
			name:     "layout tag resolution russian xbox fallback",
			input:    "Для активации потяните {LAYOUT_XBOXCONTROLLER#мини-джойстик}{LAYOUT_FALLBACK#джойстик}",
			expected: "Для активации потяните джойстик",
		},
		{
			name:     "layout tag resolution english fallback",
			input:    "While tilting the {LAYOUT_CONSOLECONTROLLER#stick}{LAYOUT_FALLBACK#joystick}, hold <IconMap:Icon_Evade> to activate",
			expected: `While tilting the joystick, hold <img src="https://enka.network/ui/zzz/IconRoleSkillKeyEvade.png" class="fairy-icon" alt="[Dodge]" style="height: 1.2em; vertical-align: -0.2em; display: inline-block;" /> to activate`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fairy.FormatHTML(tt.input, tt.skillLevel...)
			if result != tt.expected {
				t.Errorf("FormatHTML() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatPlainText(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		skillLevel []int
		expected   string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:       "unity color tag and CAL formula stripped to plain text",
			input:      "Increases ATK by <color=#2BAD00>{CAL:0.07+AvatarSkillLevel(1)*0.015,100,2}%</color>.",
			skillLevel: []int{12},
			expected:   "Increases ATK by 25%.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fairy.FormatPlainText(tt.input, tt.skillLevel...)
			if result != tt.expected {
				t.Errorf("FormatPlainText() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatMarkdown(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		skillLevel []int
		expected   string
	}{
		{
			name:       "unity color tag to markdown bold with CAL formula",
			input:      "Increases ATK by <color=#2BAD00>{CAL:0.07+AvatarSkillLevel(1)*0.015,100,2}%</color>.",
			skillLevel: []int{12},
			expected:   "Increases ATK by **25%**.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fairy.FormatMarkdown(tt.input, tt.skillLevel...)
			if result != tt.expected {
				t.Errorf("FormatMarkdown() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSkillMethods(t *testing.T) {
	skill := fairy.Skill{
		Level:       12,
		Name:        "Special Attack",
		Description: "Deals <color=#98EFF0>{CAL:0.07+AvatarSkillLevel(1)*0.015,100,2}%</color> Ice DMG.",
	}

	eval := skill.EvaluatedDescription()
	if eval != "Deals <color=#98EFF0>25%</color> Ice DMG." {
		t.Errorf("Skill.EvaluatedDescription() = %q", eval)
	}

	html := skill.FormatHTML()
	if html != `Deals <span style="color: #98EFF0;">25%</span> Ice DMG.` {
		t.Errorf("Skill.FormatHTML() = %q", html)
	}

	plain := skill.FormatPlainText()
	if plain != "Deals 25% Ice DMG." {
		t.Errorf("Skill.FormatPlainText() = %q", plain)
	}

	md := skill.FormatMarkdown()
	if md != "Deals **25%** Ice DMG." {
		t.Errorf("Skill.FormatMarkdown() = %q", md)
	}
}
