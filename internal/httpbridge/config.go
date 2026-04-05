package httpbridge

import (
	"os"
	"strings"
)

type Config struct {
	ListenAddr string
	AuthToken  string
}

func LoadConfigFromEnv() Config {
	listenAddr := strings.TrimSpace(os.Getenv("QQNOTIFY_LISTEN_ADDR"))
	if listenAddr == "" {
		if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
			listenAddr = ":" + port
		} else {
			listenAddr = ":8080"
		}
	}

	return Config{
		ListenAddr: listenAddr,
		AuthToken:  strings.TrimSpace(os.Getenv("QQNOTIFY_AUTH_TOKEN")),
	}
}
