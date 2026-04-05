package qqnotify

import (
	"fmt"
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
		lines = append(lines, "Task: "+task)
	}
	if summary := strings.TrimSpace(t.Summary); summary != "" {
		lines = append(lines, "Summary: "+summary)
	}
	if len(t.Files) > 0 {
		lines = append(lines, "Files: "+strings.Join(t.Files, ", "))
	}

	return Notification{
		Title:     "Codex task finished",
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
		lines = append(lines, "Workflow: "+workflow)
	}
	if job := strings.TrimSpace(t.Job); job != "" {
		lines = append(lines, "Job: "+job)
	}
	if summary := strings.TrimSpace(t.Summary); summary != "" {
		lines = append(lines, "Summary: "+summary)
	}
	if runURL := strings.TrimSpace(t.RunURL); runURL != "" {
		lines = append(lines, "Run URL: "+runURL)
	}

	status := normalizeStatus(t.Status)
	return Notification{
		Title:     fmt.Sprintf("CI workflow %s", status),
		Body:      strings.Join(lines, "\n"),
		Status:    status,
		Source:    "github-actions",
		TraceID:   strings.TrimSpace(t.TraceID),
		Timestamp: t.Timestamp,
	}
}

func NewCronNotification(t CronTemplate) Notification {
	var lines []string
	if name := strings.TrimSpace(t.Name); name != "" {
		lines = append(lines, "Job: "+name)
	}
	if scheduled := strings.TrimSpace(t.Scheduled); scheduled != "" {
		lines = append(lines, "Schedule: "+scheduled)
	}
	if summary := strings.TrimSpace(t.Summary); summary != "" {
		lines = append(lines, "Summary: "+summary)
	}

	status := normalizeStatus(t.Status)
	return Notification{
		Title:     fmt.Sprintf("Cron job %s", status),
		Body:      strings.Join(lines, "\n"),
		Status:    status,
		Source:    "cron",
		TraceID:   strings.TrimSpace(t.TraceID),
		Timestamp: t.Timestamp,
	}
}

func normalizeStatus(status string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return "success"
	}
	return status
}
