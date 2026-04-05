package qqnotify

import (
	"strings"
	"testing"
)

func TestLoadConfigFromEnvRequiresFields(t *testing.T) {
	t.Setenv("QQ_APP_ID", "1903697734")
	t.Setenv("QQ_APP_SECRET", "")
	t.Setenv("QQ_USER_OPENID", "")

	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Fatal("expected missing env error")
	}

	if !strings.Contains(err.Error(), "QQ_APP_SECRET") {
		t.Fatalf("expected error to mention QQ_APP_SECRET, got %v", err)
	}

	if !strings.Contains(err.Error(), "QQ_USER_OPENID") {
		t.Fatalf("expected error to mention QQ_USER_OPENID, got %v", err)
	}
}

func TestLoadConfigFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("QQ_APP_ID", "1903697734")
	t.Setenv("QQ_APP_SECRET", "secret-value")
	t.Setenv("QQ_USER_OPENID", "user-openid")
	t.Setenv("QQ_BOT_TOKEN_BASE_URL", "")
	t.Setenv("QQ_BOT_API_BASE_URL", "")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("expected config to load, got %v", err)
	}

	if cfg.TokenBaseURL != DefaultTokenBaseURL {
		t.Fatalf("expected default token base url %q, got %q", DefaultTokenBaseURL, cfg.TokenBaseURL)
	}

	if cfg.APIBaseURL != DefaultAPIBaseURL {
		t.Fatalf("expected default api base url %q, got %q", DefaultAPIBaseURL, cfg.APIBaseURL)
	}
}

func TestLoadCaptureConfigFromEnvAllowsMissingUserOpenID(t *testing.T) {
	t.Setenv("QQ_APP_ID", "1903697734")
	t.Setenv("QQ_APP_SECRET", "secret-value")
	t.Setenv("QQ_USER_OPENID", "")

	cfg, err := LoadCaptureConfigFromEnv()
	if err != nil {
		t.Fatalf("expected capture config to load, got %v", err)
	}

	if cfg.AppID != "1903697734" {
		t.Fatalf("expected app id, got %q", cfg.AppID)
	}

	if cfg.AppSecret != "secret-value" {
		t.Fatalf("expected app secret, got %q", cfg.AppSecret)
	}

	if cfg.UserOpenID != "" {
		t.Fatalf("expected empty user openid during capture, got %q", cfg.UserOpenID)
	}
}
