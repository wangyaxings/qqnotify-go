package openidcapture

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wangyaxings/qqnotify-go/qqnotify"
)

const defaultGatewayURL = "wss://api.sgroup.qq.com/websocket/"

type gatewayHelloPayload struct {
	Op int `json:"op"`
	D  struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
}

type gatewayDispatchEnvelope struct {
	Op int             `json:"op"`
	T  string          `json:"t"`
	D  json.RawMessage `json:"d"`
}

type IncomingC2CMessage struct {
	UserOpenID string
	Content    string
}

type c2cMessageData struct {
	Author struct {
		UserOpenID string `json:"user_openid"`
	} `json:"author"`
	Content string `json:"content"`
}

func BuildIdentifyPayload(token string, intents int) ([]byte, error) {
	return json.Marshal(map[string]any{
		"op": 2,
		"d": map[string]any{
			"token":   token,
			"intents": intents,
			"shard":   []int{0, 1},
			"properties": map[string]string{
				"$os":      "windows",
				"$browser": "qqnotify-openid-capture",
				"$device":  "qqnotify-openid-capture",
			},
		},
	})
}

func BuildHeartbeatPayload(seq *int) ([]byte, error) {
	var data any
	if seq != nil {
		data = *seq
	}

	return json.Marshal(map[string]any{
		"op": 1,
		"d":  data,
	})
}

func ExtractUserOpenIDFromPayload(payload []byte) (string, string, bool, error) {
	msg, ok, err := ExtractIncomingC2CMessage(payload)
	if err != nil {
		return "", "", false, err
	}
	if !ok {
		return "", "", false, nil
	}
	return msg.UserOpenID, msg.Content, true, nil
}

func ExtractIncomingC2CMessage(payload []byte) (IncomingC2CMessage, bool, error) {
	var envelope gatewayDispatchEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return IncomingC2CMessage{}, false, fmt.Errorf("parse gateway payload: %w", err)
	}

	if envelope.Op != 0 || envelope.T != "C2C_MESSAGE_CREATE" {
		return IncomingC2CMessage{}, false, nil
	}

	var data c2cMessageData
	if err := json.Unmarshal(envelope.D, &data); err != nil {
		return IncomingC2CMessage{}, false, fmt.Errorf("parse c2c payload: %w", err)
	}

	return IncomingC2CMessage{
		UserOpenID: strings.TrimSpace(data.Author.UserOpenID),
		Content:    strings.TrimSpace(data.Content),
	}, true, nil
}

func CaptureNextMatchingMessage(ctx context.Context, input <-chan IncomingC2CMessage, match func(IncomingC2CMessage) bool) (IncomingC2CMessage, error) {
	for {
		select {
		case <-ctx.Done():
			return IncomingC2CMessage{}, ctx.Err()
		case msg, ok := <-input:
			if !ok {
				return IncomingC2CMessage{}, errors.New("message stream closed before a matching message arrived")
			}
			if match(msg) {
				return msg, nil
			}
		}
	}
}

func CaptureUserOpenID(ctx context.Context, cfg qqnotify.Config) (string, string, error) {
	accessToken, err := fetchAccessToken(ctx, &http.Client{Timeout: 10 * time.Second}, cfg)
	if err != nil {
		return "", "", err
	}

	msg, err := CaptureSingleMessage(ctx, accessToken, func(IncomingC2CMessage) bool {
		return true
	})
	if err != nil {
		return "", "", err
	}
	return msg.UserOpenID, msg.Content, nil
}

func CaptureSingleMessage(ctx context.Context, accessToken string, match func(IncomingC2CMessage) bool) (IncomingC2CMessage, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, defaultGatewayURL, nil)
	if err != nil {
		return IncomingC2CMessage{}, fmt.Errorf("connect gateway websocket: %w", err)
	}
	defer conn.Close()

	heartbeatInterval, err := readHello(conn)
	if err != nil {
		return IncomingC2CMessage{}, err
	}

	identifyPayload, err := BuildIdentifyPayload("QQBot "+accessToken, 1<<25)
	if err != nil {
		return IncomingC2CMessage{}, err
	}
	if err := conn.WriteMessage(websocket.TextMessage, identifyPayload); err != nil {
		return IncomingC2CMessage{}, fmt.Errorf("send identify payload: %w", err)
	}

	go startHeartbeat(ctx, conn, heartbeatInterval)
	messages := make(chan IncomingC2CMessage)
	errs := make(chan error, 1)
	go streamIncomingMessages(conn, messages, errs)

	for {
		select {
		case <-ctx.Done():
			return IncomingC2CMessage{}, errors.New("timeout waiting for C2C_MESSAGE_CREATE event; send a fresh message to the bot while listener is running")
		case err := <-errs:
			return IncomingC2CMessage{}, err
		case msg := <-messages:
			if match(msg) {
				return msg, nil
			}
		}
	}
}

func streamIncomingMessages(conn *websocket.Conn, output chan<- IncomingC2CMessage, errs chan<- error) {
	defer close(output)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			errs <- fmt.Errorf("read gateway event: %w", err)
			return
		}

		msg, ok, err := ExtractIncomingC2CMessage(message)
		if err != nil {
			errs <- err
			return
		}
		if ok && msg.UserOpenID != "" {
			output <- msg
		}
	}
}

func readHello(conn *websocket.Conn) (time.Duration, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return 0, fmt.Errorf("read gateway hello: %w", err)
	}

	var hello gatewayHelloPayload
	if err := json.Unmarshal(message, &hello); err != nil {
		return 0, fmt.Errorf("parse gateway hello: %w", err)
	}
	if hello.Op != 10 {
		return 0, fmt.Errorf("unexpected first gateway opcode: %d", hello.Op)
	}
	if hello.D.HeartbeatInterval <= 0 {
		return 0, errors.New("gateway hello missing heartbeat interval")
	}

	return time.Duration(hello.D.HeartbeatInterval) * time.Millisecond, nil
}

func startHeartbeat(ctx context.Context, conn *websocket.Conn, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			payload, err := BuildHeartbeatPayload(nil)
			if err != nil {
				return
			}
			_ = conn.WriteMessage(websocket.TextMessage, payload)
		}
	}
}

func fetchAccessToken(ctx context.Context, client *http.Client, cfg qqnotify.Config) (string, error) {
	payload, err := json.Marshal(map[string]string{
		"appId":        cfg.AppID,
		"clientSecret": cfg.AppSecret,
	})
	if err != nil {
		return "", fmt.Errorf("marshal access token request: %w", err)
	}

	tokenURL := strings.TrimRight(cfg.TokenBaseURL, "/") + "/app/getAppAccessToken"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("build access token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("fetch access token: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decode access token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return "", errors.New("fetch access token: empty access_token in response")
	}

	return tokenResp.AccessToken, nil
}
