package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wangyaxings/qqnotify-go/internal/httpbridge"
	"github.com/wangyaxings/qqnotify-go/qqnotify"
)

func main() {
	cfg, err := qqnotify.LoadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := qqnotify.NewClient(cfg, &http.Client{Timeout: 15 * time.Second})
	addr := ":" + port
	log.Printf("qqnotifyd listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, httpbridge.NewHandler(client)); err != nil {
		log.Fatal(err)
	}
}
