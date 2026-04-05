package smokeenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFileLoadsEnvPairs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.local")
	content := "# smoke config\nQQ_APP_ID=123\nQQ_APP_SECRET=abc\nQQ_USER_OPENID=\"openid-1\"\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	if err := LoadFile(path); err != nil {
		t.Fatalf("expected env file to load, got %v", err)
	}

	if got := os.Getenv("QQ_APP_ID"); got != "123" {
		t.Fatalf("expected QQ_APP_ID, got %q", got)
	}
	if got := os.Getenv("QQ_APP_SECRET"); got != "abc" {
		t.Fatalf("expected QQ_APP_SECRET, got %q", got)
	}
	if got := os.Getenv("QQ_USER_OPENID"); got != "openid-1" {
		t.Fatalf("expected QQ_USER_OPENID, got %q", got)
	}
}

func TestLoadFirstUsesFirstExistingFile(t *testing.T) {
	dir := t.TempDir()
	second := filepath.Join(dir, ".env.local")
	if err := os.WriteFile(second, []byte("QQ_APP_ID=123\n"), 0o644); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	used, err := LoadFirst(filepath.Join(dir, "missing.env"), second)
	if err != nil {
		t.Fatalf("expected one file to load, got %v", err)
	}
	if used != second {
		t.Fatalf("expected used file %q, got %q", second, used)
	}
}
