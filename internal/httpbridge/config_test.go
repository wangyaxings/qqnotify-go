package httpbridge

import "testing"

func TestLoadConfigFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("QQNOTIFY_LISTEN_ADDR", "")
	t.Setenv("QQNOTIFY_AUTH_TOKEN", "")

	cfg := LoadConfigFromEnv()

	if cfg.ListenAddr != ":8080" {
		t.Fatalf("expected default listen addr :8080, got %q", cfg.ListenAddr)
	}
	if cfg.AuthToken != "" {
		t.Fatalf("expected empty auth token, got %q", cfg.AuthToken)
	}
}

func TestLoadConfigFromEnvUsesExplicitValues(t *testing.T) {
	t.Setenv("QQNOTIFY_LISTEN_ADDR", "127.0.0.1:9090")
	t.Setenv("QQNOTIFY_AUTH_TOKEN", "secret-token")

	cfg := LoadConfigFromEnv()

	if cfg.ListenAddr != "127.0.0.1:9090" {
		t.Fatalf("expected explicit listen addr, got %q", cfg.ListenAddr)
	}
	if cfg.AuthToken != "secret-token" {
		t.Fatalf("expected auth token, got %q", cfg.AuthToken)
	}
}

func TestLoadConfigFromEnvFallsBackToPort(t *testing.T) {
	t.Setenv("QQNOTIFY_LISTEN_ADDR", "")
	t.Setenv("PORT", "9090")

	cfg := LoadConfigFromEnv()

	if cfg.ListenAddr != ":9090" {
		t.Fatalf("expected listen addr from PORT, got %q", cfg.ListenAddr)
	}
}
