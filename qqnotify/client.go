package qqnotify

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
)

type Sender interface {
	SendText(ctx context.Context, text string) error
}

type Client struct {
	cfg        Config
	httpClient *http.Client
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type sendMessageRequest struct {
	Content string `json:"content"`
	MsgType int    `json:"msg_type"`
}

func NewClient(cfg Config, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &Client{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

func (c *Client) Send(ctx context.Context, n Notification) error {
	return c.SendText(ctx, RenderNotification(n))
}

func (c *Client) SendText(ctx context.Context, text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return errors.New("message content cannot be empty")
	}

	accessToken, err := c.fetchAccessToken(ctx)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(sendMessageRequest{
		Content: text,
		MsgType: 0,
	})
	if err != nil {
		return fmt.Errorf("marshal message request: %w", err)
	}

	messageURL := strings.TrimRight(c.cfg.APIBaseURL, "/") + "/v2/users/" + c.cfg.UserOpenID + "/messages"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, messageURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build message request: %w", err)
	}
	req.Header.Set("Authorization", "QQBot "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send qq notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("send qq notification: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}

func (c *Client) fetchAccessToken(ctx context.Context) (string, error) {
	payload, err := json.Marshal(map[string]string{
		"appId":        c.cfg.AppID,
		"clientSecret": c.cfg.AppSecret,
	})
	if err != nil {
		return "", fmt.Errorf("marshal access token request: %w", err)
	}

	tokenURL := strings.TrimRight(c.cfg.TokenBaseURL, "/") + "/app/getAppAccessToken"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("build access token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("fetch access token: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var tokenResp accessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decode access token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return "", errors.New("fetch access token: empty access_token in response")
	}

	return tokenResp.AccessToken, nil
}
