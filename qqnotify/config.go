package qqnotify

import (
	"fmt"
	"os"
	"strings"
)

const (
	DefaultTokenBaseURL = "https://bots.qq.com"
	DefaultAPIBaseURL   = "https://api.sgroup.qq.com"
)

type Config struct {
	AppID        string
	AppSecret    string
	UserOpenID   string
	TokenBaseURL string
	APIBaseURL   string
}

func LoadConfigFromEnv() (Config, error) {
	cfg := Config{
		AppID:        strings.TrimSpace(os.Getenv("QQ_APP_ID")),
		AppSecret:    strings.TrimSpace(os.Getenv("QQ_APP_SECRET")),
		UserOpenID:   strings.TrimSpace(os.Getenv("QQ_USER_OPENID")),
		TokenBaseURL: strings.TrimSpace(os.Getenv("QQ_BOT_TOKEN_BASE_URL")),
		APIBaseURL:   strings.TrimSpace(os.Getenv("QQ_BOT_API_BASE_URL")),
	}

	if cfg.TokenBaseURL == "" {
		cfg.TokenBaseURL = DefaultTokenBaseURL
	}
	if cfg.APIBaseURL == "" {
		cfg.APIBaseURL = DefaultAPIBaseURL
	}

	var missing []string
	if cfg.AppID == "" {
		missing = append(missing, "QQ_APP_ID")
	}
	if cfg.AppSecret == "" {
		missing = append(missing, "QQ_APP_SECRET")
	}
	if cfg.UserOpenID == "" {
		missing = append(missing, "QQ_USER_OPENID")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}
