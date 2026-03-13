package handler

import (
	"encoding/json"
	"net/http"

	"destinyServer/store"
)

// POST /api/analysis/free  { "open_id": "..." }
func HandleAnalysisFree(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResp(w, 405, map[string]string{"message": "method not allowed"})
		return
	}

	var req struct {
		OpenID string `json:"open_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OpenID == "" {
		jsonResp(w, 400, map[string]string{"message": "invalid request"})
		return
	}

	remaining, err := store.UseFreeCount(req.OpenID)
	if err != nil {
		jsonResp(w, 500, map[string]string{"message": "db error"})
		return
	}
	if remaining == 0 {
		if store.GetFreeCount(req.OpenID) <= 0 {
			jsonResp(w, 200, map[string]any{
				"success":    false,
				"message":    "没有免费次数了",
				"free_count": 0,
			})
			return
		}
	}

	jsonResp(w, 200, map[string]any{
		"success":    true,
		"free_count": remaining,
	})
}
