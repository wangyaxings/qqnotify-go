package httpbridge

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wangyaxings/qqnotify-go/qqnotify"
)

type Handler struct {
	sender qqnotify.Sender
	cfg    Config
}

type markdownSender interface {
	SendMarkdown(context.Context, string) error
}

func NewHandler(sender qqnotify.Sender) http.Handler {
	return NewHandlerWithConfig(sender, Config{})
}

func NewHandlerWithConfig(sender qqnotify.Sender, cfg Config) http.Handler {
	return &Handler{
		sender: sender,
		cfg:    cfg,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/healthz" {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok": true,
		})
		return
	}

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"error": "method not allowed",
		})
		return
	}

	if token := strings.TrimSpace(h.cfg.AuthToken); token != "" {
		if r.Header.Get("Authorization") != "Bearer "+token {
			writeJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "unauthorized",
			})
			return
		}
	}

	var req notifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid json payload",
		})
		return
	}

	payload, err := req.BuildNotification()
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}

	if strings.TrimSpace(payload.Title) == "" && strings.TrimSpace(payload.Body) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "title or body is required",
		})
		return
	}

	rendered := qqnotify.RenderNotification(payload)
	var sendErr error
	if sender, ok := h.sender.(markdownSender); ok {
		sendErr = sender.SendMarkdown(r.Context(), rendered)
	} else {
		sendErr = h.sender.SendText(r.Context(), rendered)
	}
	if sendErr != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{
			"error": sendErr.Error(),
		})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"ok": true,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
