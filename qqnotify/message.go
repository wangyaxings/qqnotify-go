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
		lines = append(lines, title)
	}
	if body := strings.TrimSpace(n.Body); body != "" {
		lines = append(lines, body)
	}
	if status := strings.TrimSpace(n.Status); status != "" {
		lines = append(lines, "status: "+status)
	}
	if source := strings.TrimSpace(n.Source); source != "" {
		lines = append(lines, "source: "+source)
	}
	if traceID := strings.TrimSpace(n.TraceID); traceID != "" {
		lines = append(lines, "trace_id: "+traceID)
	}
	if !n.Timestamp.IsZero() {
		lines = append(lines, "time: "+n.Timestamp.Format("2006-01-02 15:04:05"))
	}

	return strings.Join(lines, "\n")
}
