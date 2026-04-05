package qqnotify

import (
	"strings"
	"testing"
	"time"
)

func TestRenderNotificationIncludesKeyFields(t *testing.T) {
	text := RenderNotification(Notification{
		Title:     "任务完成",
		Body:      "Codex 已完成代码修改并通过测试",
		Status:    "success",
		Source:    "codex",
		TraceID:   "job-123",
		Timestamp: time.Date(2026, 4, 5, 14, 30, 0, 0, time.FixedZone("CST", 8*3600)),
	})

	for _, want := range []string{
		"任务完成",
		"Codex 已完成代码修改并通过测试",
		"success",
		"codex",
		"job-123",
		"2026-04-05 14:30:00",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected rendered text to contain %q, got %q", want, text)
		}
	}
}
