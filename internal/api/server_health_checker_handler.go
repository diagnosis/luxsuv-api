package api

import (
	"encoding/json"
	"net/http"
)

type ServerHealthCheckerHandler struct {
}

func NewServerHealthCheckerHandler() *ServerHealthCheckerHandler {
	return &ServerHealthCheckerHandler{}
}

func (sh *ServerHealthCheckerHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "luxsuv server seems healthy..."})
}
