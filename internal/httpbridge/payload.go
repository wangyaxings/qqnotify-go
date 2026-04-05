package httpbridge

import (
	"fmt"
	"strings"

	"github.com/wangyaxings/qqnotify-go/qqnotify"
)

type notifyRequest struct {
	Type      string   `json:"type,omitempty"`
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	Status    string   `json:"status,omitempty"`
	Source    string   `json:"source,omitempty"`
	TraceID   string   `json:"trace_id,omitempty"`
	Task      string   `json:"task,omitempty"`
	Summary   string   `json:"summary,omitempty"`
	Files     []string `json:"files,omitempty"`
	Workflow  string   `json:"workflow,omitempty"`
	Job       string   `json:"job,omitempty"`
	RunURL    string   `json:"run_url,omitempty"`
	Name      string   `json:"name,omitempty"`
	Scheduled string   `json:"scheduled,omitempty"`
}

func (r notifyRequest) BuildNotification() (qqnotify.Notification, error) {
	switch strings.TrimSpace(r.Type) {
	case "":
		return qqnotify.Notification{
			Title:   strings.TrimSpace(r.Title),
			Body:    strings.TrimSpace(r.Body),
			Status:  strings.TrimSpace(r.Status),
			Source:  strings.TrimSpace(r.Source),
			TraceID: strings.TrimSpace(r.TraceID),
		}, nil
	case "codex":
		if strings.TrimSpace(r.Task) == "" {
			return qqnotify.Notification{}, fmt.Errorf("task is required for codex notifications")
		}
		return qqnotify.NewCodexNotification(qqnotify.CodexTemplate{
			Task:    strings.TrimSpace(r.Task),
			Summary: strings.TrimSpace(r.Summary),
			Status:  strings.TrimSpace(r.Status),
			TraceID: strings.TrimSpace(r.TraceID),
			Files:   r.Files,
		}), nil
	case "ci":
		if strings.TrimSpace(r.Workflow) == "" {
			return qqnotify.Notification{}, fmt.Errorf("workflow is required for ci notifications")
		}
		return qqnotify.NewCINotification(qqnotify.CITemplate{
			Workflow: strings.TrimSpace(r.Workflow),
			Job:      strings.TrimSpace(r.Job),
			Status:   strings.TrimSpace(r.Status),
			Summary:  strings.TrimSpace(r.Summary),
			RunURL:   strings.TrimSpace(r.RunURL),
			TraceID:  strings.TrimSpace(r.TraceID),
		}), nil
	case "cron":
		if strings.TrimSpace(r.Name) == "" {
			return qqnotify.Notification{}, fmt.Errorf("name is required for cron notifications")
		}
		return qqnotify.NewCronNotification(qqnotify.CronTemplate{
			Name:      strings.TrimSpace(r.Name),
			Status:    strings.TrimSpace(r.Status),
			Summary:   strings.TrimSpace(r.Summary),
			Scheduled: strings.TrimSpace(r.Scheduled),
			TraceID:   strings.TrimSpace(r.TraceID),
		}), nil
	default:
		return qqnotify.Notification{}, fmt.Errorf("unsupported notification type: %s", r.Type)
	}
}
