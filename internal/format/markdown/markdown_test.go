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

func TestToSlack_Strikethrough(t *testing.T) {
	got := markdown.ToSlack("~~strike~~", noopResolve)
	if got != "~strike~" {
		t.Fatalf("want %q, got %q", "~strike~", got)
	}
}

func TestToSlack_NamedLink(t *testing.T) {
	got := markdown.ToSlack("[click here](https://example.com)", noopResolve)
	if got != "<https://example.com|click here>" {
		t.Fatalf("want %q, got %q", "<https://example.com|click here>", got)
	}
}

func TestToSlack_ImageWithAlt(t *testing.T) {
	got := markdown.ToSlack("![screenshot](https://example.com/img.png)", noopResolve)
	if got != "<https://example.com/img.png|📎 screenshot>" {
		t.Fatalf("want %q, got %q", "<https://example.com/img.png|📎 screenshot>", got)
	}
}

func TestToSlack_ImageNoAlt(t *testing.T) {
	got := markdown.ToSlack("![](https://example.com/img.png)", noopResolve)
	if got != "<https://example.com/img.png|📎 Image>" {
		t.Fatalf("want %q, got %q", "<https://example.com/img.png|📎 Image>", got)
	}
}

func TestToSlack_KramdownAttrStripped(t *testing.T) {
	got := markdown.ToSlack("![](https://example.com/img.png){: .attr}", noopResolve)
	if got != "<https://example.com/img.png|📎 Image>" {
		t.Fatalf("want %q, got %q", "<https://example.com/img.png|📎 Image>", got)
	}
}

func TestToSlack_KramdownAfterNamedLink(t *testing.T) {
	got := markdown.ToSlack("[text](https://example.com){: .class}", noopResolve)
	if got != "<https://example.com|text>" {
		t.Fatalf("want %q, got %q", "<https://example.com|text>", got)
	}
}

func TestToSlack_Mention_Resolved(t *testing.T) {
	resolve := func(id string) string {
		if id == "557058:abc123" {
			return "USLACKID"
		}
		return ""
	}
	got := markdown.ToSlack("@{557058:abc123}", resolve)
	if got != "<@USLACKID>" {
		t.Fatalf("want %q, got %q", "<@USLACKID>", got)
	}
}

func TestToSlack_Mention_Unresolved(t *testing.T) {
	got := markdown.ToSlack("@{557058:abc123}", noopResolve)
	if got != "@557058:abc123" {
		t.Fatalf("want %q, got %q", "@557058:abc123", got)
	}
}

func TestToSlack_MultilineBlock(t *testing.T) {
	input := "## Title\n\nSome **bold** text.\n\n- item one\n- item two"
	got := markdown.ToSlack(input, noopResolve)
	if !strings.Contains(got, "*Title*") {
		t.Fatalf("heading not converted; got %q", got)
	}
	if !strings.Contains(got, "*bold*") {
		t.Fatalf("bold not converted; got %q", got)
	}
	if !strings.Contains(got, "• item one") {
		t.Fatalf("list not converted; got %q", got)
	}
}
