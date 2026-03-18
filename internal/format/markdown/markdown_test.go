package markdown_test

import (
	"testing"

	"github.com/hasanMshawrab/bitslack/internal/format/markdown"
)

func noopResolve(accountID string) string { return "" }

func TestToSlack_PlainText(t *testing.T) {
	got := markdown.ToSlack("hello world", noopResolve)
	if got != "hello world" {
		t.Fatalf("want %q, got %q", "hello world", got)
	}
}
