// Package markdown converts Bitbucket markdown (CommonMark + Bitbucket extensions)
// to Slack mrkdwn format.
package markdown

// ToSlack converts Bitbucket markdown to Slack mrkdwn.
// resolve maps a Bitbucket account ID to a Slack user ID;
// it may return "" to fall back to the raw account ID.
func ToSlack(raw string, resolve func(accountID string) string) string {
	return raw
}

// Truncate shortens mrkdwn to at most maxDisplay visible characters,
// appending "…" if truncated. Links (<url|text>) are treated as atomic
// tokens whose display length equals len([]rune(text)). Truncation happens
// at the last word boundary before the limit. Open */_/~ spans are closed
// after the ellipsis.
func Truncate(mrkdwn string, maxDisplay int) string {
	return mrkdwn
}
