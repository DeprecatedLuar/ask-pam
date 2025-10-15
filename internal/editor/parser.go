package editor

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var sqlKeywords = []string{
	"LEFT JOIN", "RIGHT JOIN", "INNER JOIN", "FULL JOIN", "CROSS JOIN",
	"FULL OUTER JOIN", "LEFT OUTER JOIN", "RIGHT OUTER JOIN",
	"INSERT INTO", "DELETE FROM", "GROUP BY", "ORDER BY",
	"UNION ALL", "FETCH FIRST",
	"SELECT", "FROM", "WHERE", "JOIN", "ON", "HAVING",
	"LIMIT", "OFFSET", "UNION", "UPDATE", "VALUES", "SET",
}

var highlightKeywords = []string{
	"SELECT", "FROM", "WHERE", "JOIN", "LEFT", "RIGHT", "INNER", "FULL", "CROSS", "OUTER",
	"ON", "GROUP", "BY", "HAVING", "ORDER", "LIMIT", "OFFSET", "UNION", "ALL",
	"INSERT", "INTO", "UPDATE", "DELETE", "VALUES", "SET", "AND", "OR", "NOT",
	"IN", "EXISTS", "BETWEEN", "LIKE", "IS", "NULL", "DISTINCT", "AS",
	"CASE", "WHEN", "THEN", "ELSE", "END", "FETCH", "FIRST", "ROWS", "ONLY",
}

func FormatSQLWithLineBreaks(sql string) string {
	if sql == "" {
		return ""
	}

	formatted := sql

	// Process keywords (longer phrases first to avoid breaking them up)
	for _, keyword := range sqlKeywords {
		// Create regex pattern for word boundaries
		// This ensures we match the keyword as a whole phrase
		pattern := regexp.MustCompile(`(?i)\s+` + regexp.QuoteMeta(keyword) + `\s+`)
		
		formatted = pattern.ReplaceAllStringFunc(formatted, func(match string) string {
			// Preserve the trailing space, add newline before keyword
			return "\n" + strings.TrimSpace(match) + " "
		})
	}

	// Clean up extra newlines and trim
	lines := strings.Split(formatted, "\n")
	var cleanedLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanedLines = append(cleanedLines, trimmed)
		}
	}

	return strings.Join(cleanedLines, "\n")
}

// highlightSQL applies syntax highlighting to SQL keywords
func HighlightSQL(sql string) string {
	keywordStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	stringStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("220"))

	highlighted := sql

	// Highlight SQL keywords using word boundaries to avoid partial matches
	for _, keyword := range highlightKeywords {
		// Use regex with word boundaries to match whole words only
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)
		
		highlighted = pattern.ReplaceAllStringFunc(highlighted, func(match string) string {
			return keywordStyle.Render(match)
		})
	}

	// Highlight strings (simple implementation for single quotes)
	var result strings.Builder
	inString := false
	for _, char := range highlighted {
		if char == '\'' {
			if inString {
				result.WriteString(stringStyle.Render("'"))
				inString = false
			} else {
				result.WriteString(stringStyle.Render("'"))
				inString = true
			}
		} else if inString {
			result.WriteString(stringStyle.Render(string(char)))
		} else {
			// Check if this is part of an already styled keyword
			// (this is a simplification - in practice the styled text already has escape codes)
			result.WriteRune(char)
		}
	}

	return result.String()
}

func countLines(s string) int {
	if s == "" {
		return 1
	}
	return strings.Count(s, "\n") + 1
}
