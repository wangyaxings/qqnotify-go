package openidcapture

import (
	"context"
	"encoding/json"
	"testing"
)

func TestExtractUserOpenIDFromPayload(t *testing.T) {
	payload := []byte(`{
		"op": 0,
		"t": "C2C_MESSAGE_CREATE",
		"d": {
			"author": {
				"user_openid": "user-openid-123"
			},
			"content": "hello"
		}
	}`)

	openID, content, ok, err := ExtractUserOpenIDFromPayload(payload)
	if err != nil {
		t.Fatalf("expected payload to parse, got %v", err)
	}
	if !ok {
		t.Fatal("expected payload to be treated as c2c message event")
	}
	if openID != "user-openid-123" {
		t.Fatalf("expected user openid, got %q", openID)
	}
	if content != "hello" {
		t.Fatalf("expected content, got %q", content)
	}
}

func TestExtractUserOpenIDFromPayloadIgnoresOtherEvents(t *testing.T) {
	payload := []byte(`{"op":0,"t":"READY","d":{}}`)

	openID, content, ok, err := ExtractUserOpenIDFromPayload(payload)
	if err != nil {
		t.Fatalf("expected payload to parse, got %v", err)
	}
	if ok {
		t.Fatalf("expected non-message event to be ignored, got openid=%q content=%q", openID, content)
	}
}

func TestBuildIdentifyPayload(t *testing.T) {
	raw, err := BuildIdentifyPayload("QQBot access-token", 1<<25)
	if err != nil {
		t.Fatalf("expected identify payload to marshal, got %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}

	if payload["op"] != float64(2) {
		t.Fatalf("expected opcode 2, got %#v", payload["op"])
	}

	data, ok := payload["d"].(map[string]any)
	if !ok {
		t.Fatalf("expected object payload, got %#v", payload["d"])
	}

	if data["token"] != "QQBot access-token" {
		t.Fatalf("expected token to be preserved, got %#v", data["token"])
	}

	if data["intents"] != float64(1<<25) {
		t.Fatalf("expected intents, got %#v", data["intents"])
	}
}

func TestBuildHeartbeatPayload(t *testing.T) {
	raw, err := BuildHeartbeatPayload(nil)
	if err != nil {
		t.Fatalf("expected heartbeat payload to marshal, got %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}

	if payload["op"] != float64(1) {
		t.Fatalf("expected opcode 1, got %#v", payload["op"])
	}

	if value, exists := payload["d"]; !exists || value != nil {
		t.Fatalf("expected null heartbeat seq, got exists=%v value=%#v", exists, value)
	}
}

func TestCaptureNextMatchingMessageUsesFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := make(chan IncomingC2CMessage, 2)
	input <- IncomingC2CMessage{UserOpenID: "other-user", Content: "继续执行"}
	input <- IncomingC2CMessage{UserOpenID: "target-user", Content: "继续执行"}

	got, err := CaptureNextMatchingMessage(ctx, input, func(msg IncomingC2CMessage) bool {
		return msg.UserOpenID == "target-user"
	})
	if err != nil {
		t.Fatalf("expected matching message, got %v", err)
	}

	if got.UserOpenID != "target-user" {
		t.Fatalf("expected target user, got %q", got.UserOpenID)
	}
}
