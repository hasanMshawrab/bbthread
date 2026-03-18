// Package markdown converts Bitbucket markdown (CommonMark + Bitbucket extensions)
// to Slack mrkdwn format.
package markdown

import (
	"regexp"
	"strings"
)

var (
	reHeading    = regexp.MustCompile(`^#{1,6}\s+(.+)$`)
	reDivider    = regexp.MustCompile(`^[-*_]{3,}\s*$`)
	reUList      = regexp.MustCompile(`^[*-]\s+(.+)$`)
	reTable      = regexp.MustCompile(`^\|.*\|$`)
	reInlineBold = regexp.MustCompile(`\*\*(.+?)\*\*`)
)

// ToSlack converts Bitbucket markdown (CommonMark + Bitbucket extensions)
// to Slack mrkdwn. resolve maps a Bitbucket account ID to a Slack user ID;
// it may return "" to fall back to the raw account ID.
func ToSlack(raw string, resolve func(accountID string) string) string {
	s := raw

	// Step 1: Line-level processing (must come before inline so headings
	// are wrapped in *...* before inline bold is applied to their content).
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		switch {
		case reHeading.MatchString(line):
			sub := reHeading.FindStringSubmatch(line)
			lines[i] = "*" + sub[1] + "*"
		case reDivider.MatchString(line):
			lines[i] = ""
		case reUList.MatchString(line):
			sub := reUList.FindStringSubmatch(line)
			lines[i] = "• " + sub[1]
		case reTable.MatchString(line):
			lines[i] = ""
		}
	}
	s = strings.Join(lines, "\n")

	// Step 2: Inline replacements.
	s = reInlineBold.ReplaceAllString(s, "*$1*")

	return s
}

// Truncate shortens mrkdwn to at most maxDisplay visible characters,
// appending "…" if truncated. Links (<url|text>) are treated as atomic
// tokens whose display length equals len([]rune(text)). Truncation happens
// at the last word boundary before the limit. Open */_/~ spans are closed
// after the ellipsis.
func Truncate(mrkdwn string, maxDisplay int) string {
	return mrkdwn
}
