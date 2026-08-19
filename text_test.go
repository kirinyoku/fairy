package fairy

import (
	"testing"
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
		{
			name:       "malformed CAL tag with trailing operator",
			input:      "Static value {CAL:10+,1,2}%",
			skillLevel: 1,
			expected:   "Static value {CAL:10+,1,2}%",
		},
		{
			name:       "division by zero returns unparsed match",
			input:      "Value {CAL:10/0,1,2}%",
			skillLevel: 1,
			expected:   "Value {CAL:10/0,1,2}%",
		},
		{
			name:       "unbalanced operators return unparsed match",
			input:      "Value {CAL:10++5,1,2}%",
			skillLevel: 1,
			expected:   "Value {CAL:10++5,1,2}%",
		},
		{
			name:       "leading negative number",
			input:      "Shift {CAL:-5+10,1,2} deg",
			skillLevel: 1,
			expected:   "Shift 5 deg",
		},
		{
			name:       "subtracting negative number",
			input:      "Delta {CAL:10--5,1,2}",
			skillLevel: 1,
			expected:   "Delta 15",
		},
		{
			name:       "complex operator precedence",
			input:      "Value {CAL:10-6/2+1*4,1,2}",
			skillLevel: 1,
			expected:   "Value 11",
		},
		{
			name:       "floating point precision trimmed",
			input:      "Value {CAL:10/3,1,2}",
			skillLevel: 1,
			expected:   "Value 3.3",
		},
		{
			name:       "zero multiplier falls back to 1",
			input:      "Value {CAL:50,0,2}",
			skillLevel: 1,
			expected:   "Value 50",
		},
		{
			name:       "invalid character in formula returns unparsed match",
			input:      "Value {CAL:10@2,1,2}",
			skillLevel: 1,
			expected:   "Value {CAL:10@2,1,2}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateFormulas(tt.input, tt.skillLevel)
			if result != tt.expected {
				t.Errorf("evaluateFormulas() = %q, want %q", result, tt.expected)
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
		{
			name:     "html sanitization of script tags",
			input:    "Deals <color=#98EFF0>20%</color> DMG <script>alert(1)</script>",
			expected: `Deals <span style="color: #98EFF0;">20%</span> DMG &lt;script&gt;alert(1)&lt;/script&gt;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatHTML(tt.input, tt.skillLevel...)
			if result != tt.expected {
				t.Errorf("formatHTML() = %q, want %q", result, tt.expected)
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
			result := formatPlainText(tt.input, tt.skillLevel...)
			if result != tt.expected {
				t.Errorf("formatPlainText() = %q, want %q", result, tt.expected)
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
			result := formatMarkdown(tt.input, tt.skillLevel...)
			if result != tt.expected {
				t.Errorf("formatMarkdown() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSkillMethods(t *testing.T) {
	skill := Skill{
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

func TestFormattingWrappers(t *testing.T) {
	const sampleRichText = "Increases <color=#FE437E>ATK</color> by 15%."
	const expectedHTML = `Increases <span style="color: #FE437E;">ATK</span> by 15%.`
	const expectedPlain = "Increases ATK by 15%."
	const expectedMD = "Increases **ATK** by 15%."

	t.Run("WEngine formatting and nil receiver", func(t *testing.T) {
		var nilWEngine *WEngine
		if nilWEngine.FormatHTML() != "" {
			t.Errorf("nil WEngine FormatHTML() = %q, want empty", nilWEngine.FormatHTML())
		}
		if nilWEngine.FormatPlainText() != "" {
			t.Errorf("nil WEngine FormatPlainText() = %q, want empty", nilWEngine.FormatPlainText())
		}
		if nilWEngine.FormatMarkdown() != "" {
			t.Errorf("nil WEngine FormatMarkdown() = %q, want empty", nilWEngine.FormatMarkdown())
		}

		w := &WEngine{PassiveDescription: sampleRichText}
		if got := w.FormatHTML(); got != expectedHTML {
			t.Errorf("WEngine.FormatHTML() = %q, want %q", got, expectedHTML)
		}
		if got := w.FormatPlainText(); got != expectedPlain {
			t.Errorf("WEngine.FormatPlainText() = %q, want %q", got, expectedPlain)
		}
		if got := w.FormatMarkdown(); got != expectedMD {
			t.Errorf("WEngine.FormatMarkdown() = %q, want %q", got, expectedMD)
		}
	})

	t.Run("MindscapeNode formatting", func(t *testing.T) {
		node := MindscapeNode{Description: sampleRichText}
		if got := node.FormatHTML(); got != expectedHTML {
			t.Errorf("MindscapeNode.FormatHTML() = %q, want %q", got, expectedHTML)
		}
		if got := node.FormatPlainText(); got != expectedPlain {
			t.Errorf("MindscapeNode.FormatPlainText() = %q, want %q", got, expectedPlain)
		}
		if got := node.FormatMarkdown(); got != expectedMD {
			t.Errorf("MindscapeNode.FormatMarkdown() = %q, want %q", got, expectedMD)
		}
	})

	t.Run("PotentialVisionNode formatting", func(t *testing.T) {
		pvn := PotentialVisionNode{Description: sampleRichText}
		if got := pvn.FormatHTML(); got != expectedHTML {
			t.Errorf("PotentialVisionNode.FormatHTML() = %q, want %q", got, expectedHTML)
		}
		if got := pvn.FormatPlainText(); got != expectedPlain {
			t.Errorf("PotentialVisionNode.FormatPlainText() = %q, want %q", got, expectedPlain)
		}
		if got := pvn.FormatMarkdown(); got != expectedMD {
			t.Errorf("PotentialVisionNode.FormatMarkdown() = %q, want %q", got, expectedMD)
		}
	})

	t.Run("SetEffect formatting", func(t *testing.T) {
		eff := SetEffect{Description: sampleRichText}
		if got := eff.FormatHTML(); got != expectedHTML {
			t.Errorf("SetEffect.FormatHTML() = %q, want %q", got, expectedHTML)
		}
		if got := eff.FormatPlainText(); got != expectedPlain {
			t.Errorf("SetEffect.FormatPlainText() = %q, want %q", got, expectedPlain)
		}
		if got := eff.FormatMarkdown(); got != expectedMD {
			t.Errorf("SetEffect.FormatMarkdown() = %q, want %q", got, expectedMD)
		}
	})
}
