package main

import (
	"log"
	"net/http"

	"destinyServer/config"
	"destinyServer/handler"
	"destinyServer/store"
)

func main() {
	store.InitDB()
	defer store.DB.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/ping", handler.HandlePing)
	mux.HandleFunc("/api/login", handler.HandleLogin)
	mux.HandleFunc("/api/user/free-count", handler.HandleFreeCount)
	mux.HandleFunc("/api/analysis/free", handler.HandleAnalysisFree)
	mux.HandleFunc("/api/pay/create", handler.HandlePayCreate)
	mux.HandleFunc("/api/pay/notify", handler.HandlePayNotify)

	log.Printf("destinyServer starting on %s", config.Cfg.Port)
	if err := http.ListenAndServe(config.Cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
