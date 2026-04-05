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

	if err := client.Send(ctx, qqnotify.NewCINotification(qqnotify.CITemplate{
		Workflow:  "release",
		Job:       "build-linux",
		Status:    "failed",
		Summary:   "Unit tests failed in package qqnotify.",
		RunURL:    "https://github.com/wangyaxings/qqnotify-go/actions/runs/123",
		TraceID:   "run-123",
		Timestamp: time.Now(),
	})); err != nil {
		log.Fatal(err)
	}
}
