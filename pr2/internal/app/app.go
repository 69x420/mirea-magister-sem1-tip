package app

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/zaldk/mirea-sem1-tid-private/pr2/internal/app/handlers"
	"github.com/zaldk/mirea-sem1-tid-private/pr2/utils"
)

func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(r)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintln(w, "Hello, Go project structure!")
	})

	mux.HandleFunc("GET /ping", handlers.Ping)

	mux.HandleFunc("GET /fail", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(r)
		utils.WriteErr(w, http.StatusBadRequest, "bad_request_example")
	})

	return withRequestID(mux)
}

func Run() {
	addr := ":" + Port()

	utils.LogInfo("Server is starting on " + addr)
	srv := &http.Server{
		Addr:              addr,
		Handler:           Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		utils.LogError("server error: " + err.Error())
	}
}

func Port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = utils.NewID16()
		}
		w.Header().Set("X-Request-Id", id)
		utils.LogRequestID(id, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
