package qqnotify

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

func TestClientSendTextRequestsAccessTokenAndMessage(t *testing.T) {
	var tokenRequestBody []byte
	var messageRequestBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch path.Clean(r.URL.Path) {
		case "/app/getAppAccessToken":
			tokenRequestBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"test-access-token","expires_in":7200}`))
		case "/v2/users/user-openid/messages":
			if got := r.Header.Get("Authorization"); got != "QQBot test-access-token" {
				t.Fatalf("expected auth header, got %q", got)
			}

			messageRequestBody, _ = io.ReadAll(r.Body)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"message-id","timestamp":1743868800}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(Config{
		AppID:        "1903697734",
		AppSecret:    "secret-value",
		UserOpenID:   "user-openid",
		TokenBaseURL: server.URL,
		APIBaseURL:   server.URL,
	}, server.Client())

	if err := client.SendText(context.Background(), "task finished"); err != nil {
		t.Fatalf("expected send to succeed, got %v", err)
	}

	var tokenPayload map[string]string
	if err := json.Unmarshal(tokenRequestBody, &tokenPayload); err != nil {
		t.Fatalf("expected valid token payload, got %v", err)
	}
	if tokenPayload["appId"] != "1903697734" {
		t.Fatalf("expected appId, got %q", tokenPayload["appId"])
	}
	if tokenPayload["clientSecret"] != "secret-value" {
		t.Fatalf("expected clientSecret, got %q", tokenPayload["clientSecret"])
	}

	var messagePayload map[string]any
	if err := json.Unmarshal(messageRequestBody, &messagePayload); err != nil {
		t.Fatalf("expected valid message payload, got %v", err)
	}
	if messagePayload["content"] != "task finished" {
		t.Fatalf("expected content, got %#v", messagePayload["content"])
	}
}

func TestClientSendTextRetriesOnTransientMessageFailure(t *testing.T) {
	messageAttempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch path.Clean(r.URL.Path) {
		case "/app/getAppAccessToken":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"test-access-token","expires_in":7200}`))
		case "/v2/users/user-openid/messages":
			messageAttempts++
			if messageAttempts == 1 {
				http.Error(w, "temporary upstream error", http.StatusBadGateway)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"message-id","timestamp":1743868800}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(Config{
		AppID:        "1903697734",
		AppSecret:    "secret-value",
		UserOpenID:   "user-openid",
		TokenBaseURL: server.URL,
		APIBaseURL:   server.URL,
	}, server.Client())

	if err := client.SendText(context.Background(), "retry me"); err != nil {
		t.Fatalf("expected send to succeed after retry, got %v", err)
	}
	if messageAttempts != 2 {
		t.Fatalf("expected 2 message attempts, got %d", messageAttempts)
	}
}
