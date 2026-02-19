package render

import (
	"fmt"
	"strings"
)

// Sentinel markers for expandable quotes. These are replaced during MarkdownV2 conversion.
const (
	ExpQuoteStart = "\x02EXPQUOTE_START\x02"
	ExpQuoteEnd   = "\x02EXPQUOTE_END\x02"
)

// FormatToolUse formats a tool_use block into a summary line.
func FormatToolUse(name, input string) string {
	if input != "" {
		return fmt.Sprintf("**%s**(%s)", name, input)
	}
	return fmt.Sprintf("**%s**()", name)
}

// FormatToolResult formats a tool_result for Telegram display.
func FormatToolResult(toolName, content string, isError bool) string {
	if isError {
		return formatError(content)
	}

	lines := strings.Split(content, "\n")
	lineCount := len(lines)
	// Don't count trailing empty line
	if lineCount > 0 && lines[lineCount-1] == "" {
		lineCount--
	}

	switch toolName {
	case "Read":
		return fmt.Sprintf("Read %d lines", lineCount)
	case "Write":
		return fmt.Sprintf("Wrote %d lines", lineCount)
	case "Bash":
		summary := fmt.Sprintf("Output %d lines", lineCount)
		if content != "" {
			summary += "\n" + formatExpandableQuote(truncateContent(content, 3000))
		}
		return summary
	case "Grep":
		matchCount := countNonEmpty(lines)
		summary := fmt.Sprintf("Found %d matches", matchCount)
		if content != "" {
			summary += "\n" + formatExpandableQuote(truncateContent(content, 3000))
		}
		return summary
	case "Glob":
		fileCount := countNonEmpty(lines)
		summary := fmt.Sprintf("Found %d files", fileCount)
		if content != "" {
			summary += "\n" + formatExpandableQuote(truncateContent(content, 3000))
		}
		return summary
	case "Edit":
		added, removed := countEditChanges(content)
		summary := fmt.Sprintf("Added %d, removed %d", added, removed)
		if content != "" {
			summary += "\n" + formatExpandableQuote(truncateContent(content, 3000))
		}
		return summary
	case "Task":
		summary := fmt.Sprintf("Agent output %d lines", lineCount)
		if content != "" {
			summary += "\n" + formatExpandableQuote(truncateContent(content, 3000))
		}
		return summary
	case "WebFetch":
		charCount := len(content)
		summary := fmt.Sprintf("Fetched %d characters", charCount)
		if content != "" {
			summary += "\n" + formatExpandableQuote(truncateContent(content, 3000))
		}
		return summary
	case "WebSearch":
		resultCount := countSearchResults(content)
		summary := fmt.Sprintf("%d search results", resultCount)
		if content != "" {
			summary += "\n" + formatExpandableQuote(truncateContent(content, 3000))
		}
		return summary
	default:
		if content != "" {
			return fmt.Sprintf("%d lines", lineCount) + "\n" + formatExpandableQuote(truncateContent(content, 3000))
		}
		return fmt.Sprintf("%d lines", lineCount)
	}
}

// FormatThinking formats a thinking block: truncate to 500 chars and wrap in expandable quote.
func FormatThinking(text string) string {
	truncated := text
	if len(truncated) > 500 {
		truncated = truncated[:500] + "..."
	}
	return formatExpandableQuote(truncated)
}

// FormatText strips system tags and returns clean text.
func FormatText(text string) string {
	return text
}

// formatError formats an error tool result.
func formatError(content string) string {
	lines := strings.SplitN(content, "\n", 2)
	firstLine := lines[0]
	if len(firstLine) > 100 {
		firstLine = firstLine[:100] + "..."
	}

	result := "Error: " + firstLine
	if len(lines) > 1 {
		result += "\n" + formatExpandableQuote(truncateContent(content, 3000))
	}
	return result
}

// formatExpandableQuote wraps text in expandable quote markers.
func formatExpandableQuote(text string) string {
	return ExpQuoteStart + text + ExpQuoteEnd
}

// truncateContent truncates content to maxLen characters.
func truncateContent(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "\n... (truncated)"
}

// countNonEmpty counts non-empty lines.
func countNonEmpty(lines []string) int {
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}

// countEditChanges counts added and removed lines from a diff-like content.
func countEditChanges(content string) (added, removed int) {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			added++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			removed++
		}
	}
	return
}

// countSearchResults counts the number of search results (lines starting with a number or bullet).
func countSearchResults(content string) int {
	count := 0
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && (len(trimmed) > 0 && (trimmed[0] >= '1' && trimmed[0] <= '9' || trimmed[0] == '-' || trimmed[0] == '*')) {
			count++
		}
	}
	if count == 0 {
		count = 1 // at least 1 if there's any content
	}
	return count
}
