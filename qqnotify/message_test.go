package qqnotify

import (
	"testing"
	"time"
)

func TestRenderNotificationIncludesKeyFields(t *testing.T) {
	text := RenderNotification(Notification{
		Title:     "待确认",
		Body:      "集数：第 1 集\n结果：发现 2 个候选\n下一步：在订阅详情中确认",
		Status:    "success",
		Source:    "codex",
		TraceID:   "job-123",
		Timestamp: time.Date(2026, 4, 5, 14, 30, 0, 0, time.FixedZone("CST", 8*3600)),
	})

	want := "### 待确认\n\n- **集数**：第 1 集\n- **结果**：发现 2 个候选\n> **下一步**：在订阅详情中确认"
	if text != want {
		t.Fatalf("expected compact markdown %q, got %q", want, text)
	}
}
