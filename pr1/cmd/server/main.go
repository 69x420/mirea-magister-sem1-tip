package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type health struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	srv := &http.Server{
		Addr:              addr,
		Handler:           newRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Starting on %s ...", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func newRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", handleHello)
	mux.HandleFunc("GET /user", handleUser)
	mux.HandleFunc("GET /health", handleHealth)
	return logRequests(mux)
}

func handleHello(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("Hello, world!\n"))
}

func handleUser(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(user{
		ID:   uuid.NewString(),
		Name: "Gopher",
	})
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(health{
		Status: "ok",
		Time:   time.Now().UTC().Format(time.RFC3339),
	})
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start).Round(time.Microsecond))
	})
}
