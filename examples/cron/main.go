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

	if err := client.Send(ctx, qqnotify.NewCronNotification(qqnotify.CronTemplate{
		Name:      "daily-report",
		Status:    "success",
		Summary:   "The scheduled data sync completed successfully.",
		Scheduled: "0 9 * * *",
		TraceID:   "cron-daily-001",
		Timestamp: time.Now(),
	})); err != nil {
		log.Fatal(err)
	}
}
