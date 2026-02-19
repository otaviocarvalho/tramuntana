package render

import (
	"strings"
	"testing"
)

func TestToMarkdownV2_PlainText(t *testing.T) {
	got := ToMarkdownV2("Hello world")
	if got != "Hello world" {
		t.Errorf("got %q, want 'Hello world'", got)
	}
}

func TestToMarkdownV2_SpecialChars(t *testing.T) {
	got := ToMarkdownV2("2 + 3 = 5")
	if !strings.Contains(got, "\\+") {
		t.Errorf("should escape +: got %q", got)
	}
	if !strings.Contains(got, "\\=") {
		t.Errorf("should escape =: got %q", got)
	}
}

func TestToMarkdownV2_Bold(t *testing.T) {
	got := ToMarkdownV2("**hello**")
	// Should convert **hello** to *hello* in MarkdownV2
	if !strings.Contains(got, "*hello*") {
		t.Errorf("bold not converted: got %q", got)
	}
	// Should NOT contain double asterisks
	if strings.Contains(got, "**") {
		t.Errorf("should not have double asterisks: got %q", got)
	}
}

func TestToMarkdownV2_CodeBlock(t *testing.T) {
	input := "```\nfmt.Println(\"hello\")\n```"
	got := ToMarkdownV2(input)
	if !strings.Contains(got, "```") {
		t.Errorf("should preserve code block markers: got %q", got)
	}
	// Content inside code blocks should not have extra escaping of special chars
	if strings.Contains(got, "\\(") {
		t.Errorf("should not escape ( inside code block: got %q", got)
	}
}

func TestToMarkdownV2_InlineCode(t *testing.T) {
	input := "Use `fmt.Println()` function"
	got := ToMarkdownV2(input)
	if !strings.Contains(got, "`") {
		t.Errorf("should preserve inline code: got %q", got)
	}
}

func TestToMarkdownV2_Link(t *testing.T) {
	input := "[Google](https://google.com)"
	got := ToMarkdownV2(input)
	if !strings.Contains(got, "[") || !strings.Contains(got, "(") {
		t.Errorf("should preserve link format: got %q", got)
	}
}

func TestToMarkdownV2_DotEscape(t *testing.T) {
	got := ToMarkdownV2("file.go")
	if !strings.Contains(got, "\\.") {
		t.Errorf("should escape dot: got %q", got)
	}
}

func TestToMarkdownV2_ExclamationEscape(t *testing.T) {
	got := ToMarkdownV2("Hello!")
	if !strings.Contains(got, "\\!") {
		t.Errorf("should escape exclamation: got %q", got)
	}
}

func TestToMarkdownV2_ExpandableQuote(t *testing.T) {
	input := "Summary\n" + ExpQuoteStart + "line1\nline2" + ExpQuoteEnd
	got := ToMarkdownV2(input)
	if !strings.Contains(got, ">") {
		t.Errorf("should have blockquote marker: got %q", got)
	}
	if !strings.Contains(got, "||") {
		t.Errorf("should have spoiler for expandable: got %q", got)
	}
}

func TestToPlainText_RemovesBold(t *testing.T) {
	got := ToPlainText("**hello** world")
	if got != "hello world" {
		t.Errorf("got %q, want 'hello world'", got)
	}
}

func TestToPlainText_RemovesCode(t *testing.T) {
	got := ToPlainText("Use `fmt.Println()` here")
	if got != "Use fmt.Println() here" {
		t.Errorf("got %q", got)
	}
}

func TestToPlainText_RemovesCodeBlock(t *testing.T) {
	got := ToPlainText("```\ncode here\n```")
	if !strings.Contains(got, "code here") {
		t.Errorf("should preserve code content: got %q", got)
	}
	if strings.Contains(got, "```") {
		t.Errorf("should remove code markers: got %q", got)
	}
}

func TestToPlainText_ConvertLinks(t *testing.T) {
	got := ToPlainText("[Google](https://google.com)")
	if got != "Google (https://google.com)" {
		t.Errorf("got %q, want 'Google (https://google.com)'", got)
	}
}

func TestToPlainText_RemovesExpQuoteMarkers(t *testing.T) {
	input := "Hello " + ExpQuoteStart + "quoted" + ExpQuoteEnd + " world"
	got := ToPlainText(input)
	if strings.Contains(got, ExpQuoteStart) || strings.Contains(got, ExpQuoteEnd) {
		t.Errorf("should strip markers: got %q", got)
	}
	if !strings.Contains(got, "quoted") {
		t.Errorf("should preserve quote content: got %q", got)
	}
}

func TestEscapeMarkdownV2(t *testing.T) {
	got := escapeMarkdownV2("hello_world *bold* [link]")
	if !strings.Contains(got, "\\_") {
		t.Error("should escape underscore")
	}
	if !strings.Contains(got, "\\*") {
		t.Error("should escape asterisk")
	}
	if !strings.Contains(got, "\\[") {
		t.Error("should escape bracket")
	}
}

func TestEscapeCodeContent(t *testing.T) {
	got := escapeCodeContent("hello `world` and \\n")
	if got != "hello \\`world\\` and \\\\n" {
		t.Errorf("got %q", got)
	}
}
