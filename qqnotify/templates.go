package qqnotify

import (
	"strings"
	"time"
)

type CodexTemplate struct {
	Task      string
	Summary   string
	Status    string
	TraceID   string
	Files     []string
	Timestamp time.Time
}

type CITemplate struct {
	Workflow  string
	Job       string
	Status    string
	Summary   string
	RunURL    string
	TraceID   string
	Timestamp time.Time
}

type CronTemplate struct {
	Name      string
	Status    string
	Summary   string
	Scheduled string
	TraceID   string
	Timestamp time.Time
}

func NewCodexNotification(t CodexTemplate) Notification {
	var lines []string
	if task := strings.TrimSpace(t.Task); task != "" {
		lines = append(lines, "**任务**："+task)
	}
	if summary := strings.TrimSpace(t.Summary); summary != "" {
		lines = append(lines, summary)
	}

	status := normalizeStatus(t.Status)
	return Notification{
		Title:     statusTitle(status, "代码任务完成", "代码任务失败"),
		Body:      strings.Join(lines, "\n"),
		Status:    normalizeStatus(t.Status),
		Source:    "codex",
		TraceID:   strings.TrimSpace(t.TraceID),
		Timestamp: t.Timestamp,
	}
}

func NewCINotification(t CITemplate) Notification {
	var lines []string
	if workflow := strings.TrimSpace(t.Workflow); workflow != "" {
		lines = append(lines, "**流程**："+workflow)
	}
	if summary := strings.TrimSpace(t.Summary); summary != "" {
		lines = append(lines, summary)
	}
	if runURL := strings.TrimSpace(t.RunURL); runURL != "" {
		lines = append(lines, "[查看运行记录]("+runURL+")")
	}

	status := normalizeStatus(t.Status)
	return Notification{
		Title:     statusTitle(status, "自动化检查完成", "自动化检查失败"),
		Body:      strings.Join(lines, "\n"),
		Status:    status,
		Source:    "github-actions",
		TraceID:   strings.TrimSpace(t.TraceID),
		Timestamp: t.Timestamp,
	}
}

func NewCronNotification(t CronTemplate) Notification {
	var lines []string
	if summary := strings.TrimSpace(t.Summary); summary != "" {
		lines = append(lines, summary)
	}

	status := normalizeStatus(t.Status)
	return Notification{
		Title:     statusTitle(status, "定时任务完成", "定时任务异常"),
		Body:      strings.Join(lines, "\n"),
		Status:    status,
		Source:    "cron",
		TraceID:   strings.TrimSpace(t.TraceID),
		Timestamp: t.Timestamp,
	}
}

func statusTitle(status, success, failure string) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	if strings.Contains(normalized, "fail") || strings.Contains(normalized, "error") {
		return failure
	}
	return success
}

func normalizeStatus(status string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return "success"
	}
	return status
}
