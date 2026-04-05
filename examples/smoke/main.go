package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wangyaxings/qqnotify-go/internal/smokeenv"
	"github.com/wangyaxings/qqnotify-go/qqnotify"
)

func candidateConfigPaths() []string {
	return []string{
		"examples/smoke/.env.local",
		".env.local",
	}
}

func main() {
	usedPath, err := smokeenv.LoadFirst(candidateConfigPaths()...)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := qqnotify.LoadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	client := qqnotify.NewClientWithOptions(cfg, nil, qqnotify.ClientOptions{
		RetryAttempts: 3,
		Timeout:       20 * time.Second,
	})

	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Send(ctx, qqnotify.NewCodexNotification(qqnotify.CodexTemplate{
		Task:      "Run qqnotify-go smoke example",
		Summary:   fmt.Sprintf("Smoke test sent successfully with config file %s", usedPath),
		Status:    "success",
		TraceID:   now.Format("20060102150405"),
		Files:     []string{"examples/smoke/.env.local"},
		Timestamp: now,
	}))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("smoke notification sent successfully using %s", usedPath)
}
