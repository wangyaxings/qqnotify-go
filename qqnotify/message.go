package qqnotify

import (
	"strings"
	"time"
)

type Notification struct {
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    string    `json:"status,omitempty"`
	Source    string    `json:"source,omitempty"`
	TraceID   string    `json:"trace_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func RenderNotification(n Notification) string {
	var lines []string

	if title := strings.TrimSpace(n.Title); title != "" {
		lines = append(lines, "### "+title)
	}
	if body := strings.TrimSpace(n.Body); body != "" {
		lines = append(lines, formatMarkdownBody(body))
	}

	return strings.Join(lines, "\n\n")
}

func formatMarkdownBody(body string) string {
	formatted := make([]string, 0, 4)
	for _, raw := range strings.Split(body, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		label, value, found := strings.Cut(line, "：")
		switch label {
		case "集数", "结果", "原因", "状态":
			if found {
				line = "- **" + label + "**：" + strings.TrimSpace(value)
			}
		case "下一步":
			if found {
				line = "> **下一步**：" + strings.TrimSpace(value)
			}
		}
		formatted = append(formatted, line)
	}
	return strings.Join(formatted, "\n")
}
