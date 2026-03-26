package resume

import (
	"regexp"
	"strings"
	"unicode"
)

// CleanText normalizes and sanitizes raw extracted PDF text for LLM consumption.
func CleanText(raw string) string {
	// Replace common PDF artifacts
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")

	// Remove non-printable characters (keep newlines, tabs, spaces)
	raw = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' || r == ' ' {
			return r
		}
		if !unicode.IsPrint(r) {
			return -1
		}
		return r
	}, raw)

	// Collapse multiple spaces into one
	spaceRegex := regexp.MustCompile(`[ \t]+`)
	raw = spaceRegex.ReplaceAllString(raw, " ")

	// Collapse 3+ newlines into 2
	newlineRegex := regexp.MustCompile(`\n{3,}`)
	raw = newlineRegex.ReplaceAllString(raw, "\n\n")

	// Trim each line
	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	raw = strings.Join(lines, "\n")

	// Remove lines that are just repeated special chars (e.g. "-----", "=====")
	junkLineRegex := regexp.MustCompile(`^[\-=_\.\*\#\~]{3,}$`)
	var cleanLines []string
	for _, line := range strings.Split(raw, "\n") {
		if junkLineRegex.MatchString(line) {
			continue
		}
		cleanLines = append(cleanLines, line)
	}

	result := strings.Join(cleanLines, "\n")
	return strings.TrimSpace(result)
}
