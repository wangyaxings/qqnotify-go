package main

import (
	"log"
	"net/http"
	"time"

	"github.com/wangyaxings/qqnotify-go/internal/httpbridge"
	"github.com/wangyaxings/qqnotify-go/qqnotify"
)

func main() {
	cfg, err := qqnotify.LoadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	bridgeCfg := httpbridge.LoadConfigFromEnv()
	client := qqnotify.NewClient(cfg, &http.Client{Timeout: 15 * time.Second})
	log.Printf("qqnotifyd listening on http://localhost%s", bridgeCfg.ListenAddr)
	if bridgeCfg.AuthToken != "" {
		log.Print("qqnotifyd auth is enabled for /notify")
	}
	if err := http.ListenAndServe(bridgeCfg.ListenAddr, httpbridge.NewHandlerWithConfig(client, bridgeCfg)); err != nil {
		log.Fatal(err)
	}
}
