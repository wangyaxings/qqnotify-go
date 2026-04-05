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
