package handler

import "net/http"

func HandlePing(w http.ResponseWriter, r *http.Request) {
	jsonResp(w, 200, map[string]string{"message": "pong"})
}
