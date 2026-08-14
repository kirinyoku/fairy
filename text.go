package fairy

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
)

var (
	colorTagRegex        = regexp.MustCompile(`(?i)<color=(#[0-9a-fA-F]{3,8})>(.*?)</color>`)
	bracketColorRegex    = regexp.MustCompile(`(?i)<color=(#[0-9a-fA-F]{3,8})>\[([^\]]+)\]</color>`)
	termBracketRegex     = regexp.MustCompile(`\[([A-Z][A-Za-z0-9_\s:\'\-·]+)\]`)
	iconTagRegex         = regexp.MustCompile(`<IconMap:Icon_([a-zA-Z0-9_]+)>`)
	anyTagRegex          = regexp.MustCompile(`<[^>]*>`)
	calTagRegex          = regexp.MustCompile(`\{CAL:([^,]+),([^,]+),([^\}]+)\}`)
	skillLevelRegex      = regexp.MustCompile(`(?i)AvatarSkillLevel\(\d+\)`)
	layoutSequenceRegex  = regexp.MustCompile(`(\{LAYOUT_[^}]+\})+`)
	layoutSingleTagRegex = regexp.MustCompile(`\{LAYOUT_([^#]+)#([^}]+)\}`)
)

func resolveLayoutTags(text string) string {
	return layoutSequenceRegex.ReplaceAllStringFunc(text, func(match string) string {
		matches := layoutSingleTagRegex.FindAllStringSubmatch(match, -1)
		if len(matches) == 0 {
			return match
		}
		for _, m := range matches {
			if len(m) >= 3 && strings.Contains(m[1], "FALLBACK") {
				return m[2]
			}
		}
		return matches[len(matches)-1][2]
	})
}

func unwrapTermBrackets(text string) string {
	res := bracketColorRegex.ReplaceAllString(text, "<color=$1>$2</color>")
	return termBracketRegex.ReplaceAllStringFunc(res, func(match string) string {
		sub := termBracketRegex.FindStringSubmatch(match)
		if len(sub) > 1 {
			return sub[1]
		}
		return match
	})
}

var iconImageMap = map[string]string{
	"<IconMap:Icon_UltimateReady>":   EnkaAssetBaseURL + "IconRoleSkillKeyUltimateV2.png",
	"<IconMap:Icon_SpecialReady>":    EnkaAssetBaseURL + "IconRoleSkillKeySpecialV2.png",
	"<IconMap:Icon_SpecialReady_Rp>": EnkaAssetBaseURL + "IconRoleSkillKeySpecialV2.png",
	"<IconMap:Icon_Special>":         EnkaAssetBaseURL + "IconRoleSkillKeySpecial.png",
	"<IconMap:Icon_Normal>":          EnkaAssetBaseURL + "IconRoleSkillKeyNormal.png",
	"<IconMap:Icon_Evade>":           EnkaAssetBaseURL + "IconRoleSkillKeyEvade.png",
	"<IconMap:Icon_Switch>":          EnkaAssetBaseURL + "IconRoleSkillKeySwitch.png",
}

var iconLabelMap = map[string]string{
	"<IconMap:Icon_UltimateReady>":   "[Ultimate]",
	"<IconMap:Icon_SpecialReady>":    "[EX Special Attack]",
	"<IconMap:Icon_SpecialReady_Rp>": "[EX Special Attack]",
	"<IconMap:Icon_Special>":         "[Special Attack]",
	"<IconMap:Icon_Normal>":          "[Basic Attack]",
	"<IconMap:Icon_Evade>":           "[Dodge]",
	"<IconMap:Icon_Switch>":          "[Switch]",
}

var (
	iconPlainTextReplacer *strings.Replacer
	iconMarkdownReplacer  *strings.Replacer
)

func init() {
	var plainPairs []string
	var mdPairs []string

	for tag, label := range iconLabelMap {
		plainPairs = append(plainPairs, tag, label)
		mdPairs = append(mdPairs, tag, "**"+label+"**")
	}
	iconPlainTextReplacer = strings.NewReplacer(plainPairs...)
	iconMarkdownReplacer = strings.NewReplacer(mdPairs...)
}

// EvaluateFormulas evaluates and replaces all Unity skill calculation formulas
// in the format {CAL:expr,mult,precision} with calculated values for the given skill level.
func EvaluateFormulas(text string, skillLevel int) string {
	if text == "" || !strings.Contains(text, "{CAL:") {
		return text
	}

	return calTagRegex.ReplaceAllStringFunc(text, func(match string) string {
		sub := calTagRegex.FindStringSubmatch(match)
		if len(sub) < 4 {
			return match
		}

		rawExpr := sub[1]
		multStr := sub[2]

		mult, err := strconv.ParseFloat(multStr, 64)
		if err != nil || mult == 0 {
			mult = 1
		}

		expr := skillLevelRegex.ReplaceAllString(rawExpr, strconv.Itoa(skillLevel))

		val, err := evaluateSimpleExpr(expr)
		if err != nil {
			return match
		}

		finalVal := val * mult

		if finalVal == float64(int64(finalVal)) {
			return fmt.Sprintf("%.0f", finalVal)
		}
		formatted := fmt.Sprintf("%.1f", finalVal)
		formatted = strings.TrimRight(formatted, "0")
		formatted = strings.TrimRight(formatted, ".")
		return formatted
	})
}

// evaluateSimpleExpr parses and evaluates basic linear math expressions (+, -, *, /).
func evaluateSimpleExpr(expr string) (float64, error) {
	if strings.Contains(expr, " ") {
		expr = strings.ReplaceAll(expr, " ", "")
	}
	terms, ops, err := tokenizeSimpleExpr(expr)
	if err != nil {
		return 0, err
	}
	return evaluateTokens(terms, ops)
}

func tokenizeSimpleExpr(expr string) ([]float64, []rune, error) {
	var terms []float64
	var ops []rune

	i := 0
	for i < len(expr) {
		ch := expr[i]
		if (ch >= '0' && ch <= '9') || ch == '.' || (ch == '-' && (i == 0 || expr[i-1] == '+' || expr[i-1] == '-' || expr[i-1] == '*' || expr[i-1] == '/')) {
			start := i
			i++
			for i < len(expr) && ((expr[i] >= '0' && expr[i] <= '9') || expr[i] == '.') {
				i++
			}
			val, err := strconv.ParseFloat(expr[start:i], 64)
			if err != nil {
				return nil, nil, err
			}
			terms = append(terms, val)
		} else if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			ops = append(ops, rune(ch))
			i++
		} else {
			i++
		}
	}

	if len(terms) == 0 {
		return nil, nil, fmt.Errorf("empty expression")
	}
	if len(terms) != len(ops)+1 {
		return nil, nil, fmt.Errorf("invalid expression: mismatched terms and operators")
	}

	return terms, ops, nil
}

func evaluateTokens(terms []float64, ops []rune) (float64, error) {
	// 1. Multiplication and division
	var finalTerms []float64
	var finalOps []rune

	finalTerms = append(finalTerms, terms[0])
	for idx := 0; idx < len(ops); idx++ {
		op := ops[idx]
		if idx+1 >= len(terms) {
			return 0, fmt.Errorf("invalid expression: missing term")
		}
		nextVal := terms[idx+1]

		switch op {
		case '*':
			finalTerms[len(finalTerms)-1] *= nextVal
		case '/':
			if nextVal == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			finalTerms[len(finalTerms)-1] /= nextVal
		default:
			finalTerms = append(finalTerms, nextVal)
			finalOps = append(finalOps, op)
		}
	}

	// 2. Addition and subtraction
	res := finalTerms[0]
	for idx := 0; idx < len(finalOps); idx++ {
		if idx+1 >= len(finalTerms) {
			return 0, fmt.Errorf("invalid expression: missing term")
		}
		switch finalOps[idx] {
		case '+':
			res += finalTerms[idx+1]
		case '-':
			res -= finalTerms[idx+1]
		}
	}

	return res, nil
}

// prepareText handles empty string checks, optional formula evaluation,
// and unwraps color bracket term syntax.
func prepareText(text string, skillLevel ...int) string {
	if text == "" {
		return ""
	}
	if len(skillLevel) > 0 {
		text = EvaluateFormulas(text, skillLevel[0])
	}
	return unwrapTermBrackets(text)
}

// FormatHTML converts Unity Rich Text tags (e.g. <color=#2BAD00>20%</color> and <IconMap:Icon_Special>)
// into web-compatible HTML with inline CSS styling, Enka CDN icon image tags, and break tags (<br>).
// Untrusted HTML content is safely escaped to prevent XSS attacks.
// If an optional skillLevel is passed, any {CAL:...} scaling formulas in text are evaluated automatically.
func FormatHTML(text string, skillLevel ...int) string {
	text = prepareText(text, skillLevel...)
	if text == "" {
		return ""
	}

	// Step 1: Extract Unity color tags and convert them to placeholders with escaped content
	var spanPlaceholders []string
	text = colorTagRegex.ReplaceAllStringFunc(text, func(match string) string {
		sub := colorTagRegex.FindStringSubmatch(match)
		if len(sub) >= 3 {
			color := html.EscapeString(sub[1])
			content := html.EscapeString(sub[2])
			span := fmt.Sprintf(`<span style="color: %s;">%s</span>`, color, content)
			spanPlaceholders = append(spanPlaceholders, span)
			return "\x00SPAN_" + strconv.Itoa(len(spanPlaceholders)-1) + "\x00"
		}
		return match
	})

	// Step 2: Extract IconMap tags to placeholders
	var iconPlaceholders []string
	text = iconTagRegex.ReplaceAllStringFunc(text, func(match string) string {
		sub := iconTagRegex.FindStringSubmatch(match)
		if len(sub) > 1 {
			iconName := sub[1]
			fullTag := fmt.Sprintf("<IconMap:Icon_%s>", iconName)
			var imgTag string
			if imgURL, ok := iconImageMap[fullTag]; ok {
				label := iconLabelMap[fullTag]
				imgTag = fmt.Sprintf(`<img src="%s" class="fairy-icon" alt="%s" style="height: 1.2em; vertical-align: -0.2em; display: inline-block;" />`, imgURL, label)
			} else {
				imgURL := fmt.Sprintf("%s%s.png", EnkaAssetBaseURL, iconName)
				imgTag = fmt.Sprintf(`<img src="%s" class="fairy-icon" alt="%s" style="height: 1.2em; vertical-align: -0.2em; display: inline-block;" />`, imgURL, iconName)
			}
			iconPlaceholders = append(iconPlaceholders, imgTag)
			return "\x00ICON_" + strconv.Itoa(len(iconPlaceholders)-1) + "\x00"
		}
		return match
	})

	// Step 3: HTML-escape the remaining text to sanitize any injected script or HTML tags
	res := html.EscapeString(text)

	// Step 4: Restore placeholders in a single pass
	if len(spanPlaceholders) > 0 || len(iconPlaceholders) > 0 {
		replacements := make([]string, 0, (len(spanPlaceholders)+len(iconPlaceholders))*2)
		for i, span := range spanPlaceholders {
			replacements = append(replacements, "\x00SPAN_"+strconv.Itoa(i)+"\x00", span)
		}
		for i, img := range iconPlaceholders {
			replacements = append(replacements, "\x00ICON_"+strconv.Itoa(i)+"\x00", img)
		}
		res = strings.NewReplacer(replacements...).Replace(res)
	}

	// Resolve layout fallback tags to native localized terms
	res = resolveLayoutTags(res)

	// Replace newlines with <br>
	res = strings.ReplaceAll(res, "\n", "<br>")

	return res
}

// FormatPlainText strips all Unity Rich Text tags, color tags, and icon placeholders from text,
// replacing IconMap tags with clean readable labels (e.g. [Ultimate], [Special Attack]),
// returning clean, human-readable plain text without any markup.
// If an optional skillLevel is passed, any {CAL:...} scaling formulas in text are evaluated automatically.
func FormatPlainText(text string, skillLevel ...int) string {
	text = prepareText(text, skillLevel...)
	if text == "" {
		return ""
	}

	// Unwrap color tags to keep inner text content
	res := colorTagRegex.ReplaceAllString(text, "$2")

	// Resolve layout fallback tags to native localized terms
	res = resolveLayoutTags(res)

	// Replace known IconMap tags with clean labels
	res = iconPlainTextReplacer.Replace(res)

	// Strip any remaining IconMap placeholders
	res = iconTagRegex.ReplaceAllString(res, "")

	// Strip any remaining XML/HTML tags
	res = anyTagRegex.ReplaceAllString(res, "")

	return res
}

// FormatMarkdown converts Unity Rich Text tags into Markdown-formatted text.
// Colored values are wrapped in bold (**text**), IconMap placeholders are replaced with bold labels (**[Ultimate]**),
// and original text layout is preserved for platforms like Discord, Telegram, or Slack.
// If an optional skillLevel is passed, any {CAL:...} scaling formulas in text are evaluated automatically.
func FormatMarkdown(text string, skillLevel ...int) string {
	text = prepareText(text, skillLevel...)
	if text == "" {
		return ""
	}

	// Convert color tags to bold Markdown
	res := colorTagRegex.ReplaceAllString(text, "**$2**")

	// Resolve layout fallback tags to native localized terms
	res = resolveLayoutTags(res)

	// Replace known IconMap tags with bold labels
	res = iconMarkdownReplacer.Replace(res)

	// Strip any remaining IconMap placeholders
	res = iconTagRegex.ReplaceAllString(res, "")

	// Strip any remaining XML/HTML tags
	res = anyTagRegex.ReplaceAllString(res, "")

	return res
}
