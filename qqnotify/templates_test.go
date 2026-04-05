package qqnotify

import (
	"strings"
	"testing"
	"time"
)

func TestNewCodexNotificationBuildsStructuredMessage(t *testing.T) {
	n := NewCodexNotification(CodexTemplate{
		Task:      "Refactor notification bridge",
		Summary:   "All tests passed and the patch is ready.",
		Status:    "success",
		TraceID:   "codex-123",
		Files:     []string{"internal/httpbridge/handler.go", "README.md"},
		Timestamp: time.Date(2026, 4, 5, 16, 30, 0, 0, time.FixedZone("CST", 8*3600)),
	})

	if n.Title != "Codex task finished" {
		t.Fatalf("expected codex title, got %q", n.Title)
	}
	for _, want := range []string{
		"Task: Refactor notification bridge",
		"Summary: All tests passed and the patch is ready.",
		"Files: internal/httpbridge/handler.go, README.md",
	} {
		if !strings.Contains(n.Body, want) {
			t.Fatalf("expected body to contain %q, got %q", want, n.Body)
		}
	}
	if n.Source != "codex" {
		t.Fatalf("expected source codex, got %q", n.Source)
	}
}

func TestNewCINotificationBuildsStructuredMessage(t *testing.T) {
	n := NewCINotification(CITemplate{
		Workflow: "release",
		Job:      "build-linux",
		Status:   "failed",
		Summary:  "Unit tests failed in package qqnotify.",
		RunURL:   "https://github.com/example/repo/actions/runs/123",
		TraceID:  "run-123",
	})

	if n.Title != "CI workflow failed" {
		t.Fatalf("expected ci title, got %q", n.Title)
	}
	for _, want := range []string{
		"Workflow: release",
		"Job: build-linux",
		"Summary: Unit tests failed in package qqnotify.",
		"Run URL: https://github.com/example/repo/actions/runs/123",
	} {
		if !strings.Contains(n.Body, want) {
			t.Fatalf("expected body to contain %q, got %q", want, n.Body)
		}
	}
	if n.Source != "github-actions" {
		t.Fatalf("expected source github-actions, got %q", n.Source)
	}
}

func TestNewCronNotificationBuildsStructuredMessage(t *testing.T) {
	n := NewCronNotification(CronTemplate{
		Name:      "daily-report",
		Status:    "success",
		Summary:   "The daily report was generated successfully.",
		Scheduled: "0 9 * * *",
		TraceID:   "cron-001",
	})

	if n.Title != "Cron job success" {
		t.Fatalf("expected cron title, got %q", n.Title)
	}
	for _, want := range []string{
		"Job: daily-report",
		"Schedule: 0 9 * * *",
		"Summary: The daily report was generated successfully.",
	} {
		if !strings.Contains(n.Body, want) {
			t.Fatalf("expected body to contain %q, got %q", want, n.Body)
		}
	}
	if n.Source != "cron" {
		t.Fatalf("expected source cron, got %q", n.Source)
	}
}
