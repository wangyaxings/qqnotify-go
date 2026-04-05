package main

import (
	"context"
	"log"
	"time"

	"github.com/wangyaxings/qqnotify-go/qqnotify"
)

func main() {
	cfg, err := qqnotify.LoadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	client := qqnotify.NewClient(cfg, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Send(ctx, qqnotify.Notification{
		Title:   "Codex task finished",
		Body:    "Patch generated and verification completed.",
		Status:  "success",
		Source:  "codex",
		TraceID: "codex-run-001",
	}); err != nil {
		log.Fatal(err)
	}
}
