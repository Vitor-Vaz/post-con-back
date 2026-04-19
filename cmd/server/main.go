package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"post-con-back/extension/database"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	fmt.Println("post-con-back listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
