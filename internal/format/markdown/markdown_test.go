package markdown_test

import (
	"strings"
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

func TestToSlack_Heading1(t *testing.T) {
	got := markdown.ToSlack("# Title", noopResolve)
	if got != "*Title*" {
		t.Fatalf("want %q, got %q", "*Title*", got)
	}
}

func TestToSlack_Heading2(t *testing.T) {
	got := markdown.ToSlack("## Section", noopResolve)
	if got != "*Section*" {
		t.Fatalf("want %q, got %q", "*Section*", got)
	}
}

func TestToSlack_UnorderedListDash(t *testing.T) {
	got := markdown.ToSlack("- first item", noopResolve)
	if got != "• first item" {
		t.Fatalf("want %q, got %q", "• first item", got)
	}
}

func TestToSlack_UnorderedListStar(t *testing.T) {
	got := markdown.ToSlack("* second item", noopResolve)
	if got != "• second item" {
		t.Fatalf("want %q, got %q", "• second item", got)
	}
}

func TestToSlack_ListStarWithNestedBold(t *testing.T) {
	// * at line start is a list item; **bold** inside is inline bold.
	got := markdown.ToSlack("* item with **bold**", noopResolve)
	if got != "• item with *bold*" {
		t.Fatalf("want %q, got %q", "• item with *bold*", got)
	}
}

func TestToSlack_OrderedListUnchanged(t *testing.T) {
	got := markdown.ToSlack("1. first\n2. second", noopResolve)
	if got != "1. first\n2. second" {
		t.Fatalf("want %q, got %q", "1. first\n2. second", got)
	}
}

func TestToSlack_Divider(t *testing.T) {
	got := markdown.ToSlack("---", noopResolve)
	if got != "" {
		t.Fatalf("want empty string, got %q", got)
	}
}

func TestToSlack_TableStripped(t *testing.T) {
	input := "| col1 | col2 |\n|------|------|\n| a    | b    |"
	got := markdown.ToSlack(input, noopResolve)
	if strings.Contains(got, "|") {
		t.Fatalf("expected table pipes removed, got %q", got)
	}
}

func TestToSlack_MixedDividerNotStripped(t *testing.T) {
	// "-*-" is not a valid CommonMark divider — must not be stripped.
	got := markdown.ToSlack("-*-", noopResolve)
	if got == "" {
		t.Fatalf("mixed divider should not be stripped, got empty string")
	}
}
