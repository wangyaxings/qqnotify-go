package httpbridge

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeSender struct {
	lastText string
	err      error
}

func (f *fakeSender) SendText(_ context.Context, text string) error {
	f.lastText = text
	return f.err
}

func TestHandlerAcceptsJSONAndSendsNotification(t *testing.T) {
	sender := &fakeSender{}
	handler := NewHandler(sender)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{"title":"任务完成","body":"构建成功","status":"success"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202 accepted, got %d", rec.Code)
	}
	if !strings.Contains(sender.lastText, "任务完成") {
		t.Fatalf("expected sender to receive rendered text, got %q", sender.lastText)
	}
}

func TestHandlerRejectsInvalidPayload(t *testing.T) {
	sender := &fakeSender{}
	handler := NewHandler(sender)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{"title":"","body":""}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 bad request, got %d", rec.Code)
	}
}

func TestHandlerRespondsToHealthCheck(t *testing.T) {
	sender := &fakeSender{}
	handler := NewHandler(sender)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for health check, got %d", rec.Code)
	}
}

func TestHandlerRequiresBearerTokenWhenConfigured(t *testing.T) {
	sender := &fakeSender{}
	handler := NewHandlerWithConfig(sender, Config{
		AuthToken: "secret-token",
	})

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{"title":"任务完成","body":"构建成功"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", rec.Code)
	}
}

func TestHandlerAcceptsBearerTokenWhenConfigured(t *testing.T) {
	sender := &fakeSender{}
	handler := NewHandlerWithConfig(sender, Config{
		AuthToken: "secret-token",
	})

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{"title":"任务完成","body":"构建成功"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202 with valid token, got %d", rec.Code)
	}
}

func TestHandlerBuildsCodexTemplatePayload(t *testing.T) {
	sender := &fakeSender{}
	handler := NewHandler(sender)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{
		"type":"codex",
		"task":"Refactor bridge auth",
		"summary":"All tests passed.",
		"status":"success",
		"trace_id":"codex-456",
		"files":["internal/httpbridge/handler.go","README.md"]
	}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202 for codex template payload, got %d", rec.Code)
	}
	for _, want := range []string{
		"Codex task finished",
		"Task: Refactor bridge auth",
		"Summary: All tests passed.",
		"Files: internal/httpbridge/handler.go, README.md",
	} {
		if !strings.Contains(sender.lastText, want) {
			t.Fatalf("expected rendered text to contain %q, got %q", want, sender.lastText)
		}
	}
}

func TestHandlerRejectsUnknownTemplateType(t *testing.T) {
	sender := &fakeSender{}
	handler := NewHandler(sender)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{
		"type":"unknown",
		"summary":"no-op"
	}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for unknown template type, got %d", rec.Code)
	}
}

func TestHandlerBuildsCITemplatePayload(t *testing.T) {
	sender := &fakeSender{}
	handler := NewHandler(sender)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{
		"type":"ci",
		"workflow":"release",
		"job":"build-linux",
		"summary":"Unit tests failed.",
		"status":"failed",
		"run_url":"https://github.com/example/repo/actions/runs/123"
	}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202 for ci template payload, got %d", rec.Code)
	}
	for _, want := range []string{
		"CI workflow failed",
		"Workflow: release",
		"Job: build-linux",
		"Run URL: https://github.com/example/repo/actions/runs/123",
	} {
		if !strings.Contains(sender.lastText, want) {
			t.Fatalf("expected rendered text to contain %q, got %q", want, sender.lastText)
		}
	}
}

func TestHandlerBuildsCronTemplatePayload(t *testing.T) {
	sender := &fakeSender{}
	handler := NewHandler(sender)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{
		"type":"cron",
		"name":"daily-report",
		"summary":"The report was generated successfully.",
		"status":"success",
		"scheduled":"0 9 * * *"
	}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202 for cron template payload, got %d", rec.Code)
	}
	for _, want := range []string{
		"Cron job success",
		"Job: daily-report",
		"Schedule: 0 9 * * *",
	} {
		if !strings.Contains(sender.lastText, want) {
			t.Fatalf("expected rendered text to contain %q, got %q", want, sender.lastText)
		}
	}
}
