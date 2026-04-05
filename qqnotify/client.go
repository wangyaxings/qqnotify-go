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
	options    ClientOptions
}

const (
	defaultMessageRetryAttempts = 2
	defaultHTTPTimeout          = 10 * time.Second
)

type ClientOptions struct {
	RetryAttempts int
	Timeout       time.Duration
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type sendMessageRequest struct {
	Content string `json:"content"`
	MsgType int    `json:"msg_type"`
}

type tokenErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(cfg Config, httpClient *http.Client) *Client {
	return NewClientWithOptions(cfg, httpClient, ClientOptions{})
}

func NewClientWithOptions(cfg Config, httpClient *http.Client, options ClientOptions) *Client {
	if options.Timeout <= 0 {
		options.Timeout = defaultHTTPTimeout
	}
	if options.RetryAttempts <= 0 {
		options.RetryAttempts = defaultMessageRetryAttempts
	}

	if httpClient == nil {
		httpClient = &http.Client{Timeout: options.Timeout}
	}

	return &Client{
		cfg:        cfg,
		httpClient: httpClient,
		options:    options,
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

	var lastErr error
	for attempt := 0; attempt < c.options.RetryAttempts; attempt++ {
		resp, err := c.httpClient.Do(req.Clone(ctx))
		if err != nil {
			lastErr = fmt.Errorf("send qq notification: %w", err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		lastErr = fmt.Errorf("send qq notification: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
		if resp.StatusCode < 500 {
			return lastErr
		}
	}

	return lastErr
}

func (c *Client) fetchAccessToken(ctx context.Context) (string, error) {
	return FetchAccessToken(ctx, c.httpClient, c.cfg)
}

func FetchAccessToken(ctx context.Context, client *http.Client, cfg Config) (string, error) {
	if client == nil {
		client = &http.Client{Timeout: defaultHTTPTimeout}
	}

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

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("fetch access token: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var tokenResp accessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("decode access token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return "", formatTokenResponseError(body)
	}

	return tokenResp.AccessToken, nil
}

func formatTokenResponseError(body []byte) error {
	var apiErr tokenErrorResponse
	if err := json.Unmarshal(body, &apiErr); err == nil {
		message := strings.TrimSpace(strings.ToLower(apiErr.Message))
		if apiErr.Code == 100016 || message == "invalid appid or secret" {
			return errors.New("fetch access token: invalid QQ_APP_ID or QQ_APP_SECRET; verify the copied values or regenerate AppSecret if needed")
		}
	}

	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return errors.New("fetch access token: empty access_token in response")
	}

	return fmt.Errorf("fetch access token: empty access_token in response: %s", trimmed)
}
