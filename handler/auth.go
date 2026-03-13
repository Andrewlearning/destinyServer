package handler

import (
	"encoding/json"
	"net/http"

	"destinyServer/store"
	"destinyServer/wechat"
)

// POST /api/login  { "code": "...", "referrer": "..." }
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResp(w, 405, map[string]string{"message": "method not allowed"})
		return
	}

	var req struct {
		Code     string `json:"code"`
		Referrer string `json:"referrer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Code == "" {
		jsonResp(w, 400, map[string]string{"message": "invalid request"})
		return
	}

	session, err := wechat.Code2Session(req.Code)
	if err != nil {
		jsonResp(w, 500, map[string]string{"message": "wx login failed"})
		return
	}

	freeCount, err := store.GetOrCreateUser(session.OpenID)
	if err != nil {
		jsonResp(w, 500, map[string]string{"message": "db error"})
		return
	}

	if req.Referrer != "" {
		_ = store.AddReferralBonus(req.Referrer, session.OpenID)
	}

	jsonResp(w, 200, map[string]any{
		"open_id":    session.OpenID,
		"free_count": freeCount,
	})
}

// GET /api/user/free-count?open_id=xxx
func HandleFreeCount(w http.ResponseWriter, r *http.Request) {
	openID := r.URL.Query().Get("open_id")
	if openID == "" {
		jsonResp(w, 400, map[string]string{"message": "missing open_id"})
		return
	}

	freeCount, err := store.GetOrCreateUser(openID)
	if err != nil {
		jsonResp(w, 500, map[string]string{"message": "db error"})
		return
	}
	jsonResp(w, 200, map[string]any{"free_count": freeCount})
}
