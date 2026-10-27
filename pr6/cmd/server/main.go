package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zaldk/mirea-sem1-tid-private/pr6/internal/db"
	"github.com/zaldk/mirea-sem1-tid-private/pr6/internal/httpapi"
)

func main() {
	d := db.Connect()

	if err := db.Migrate(d); err != nil {
		log.Fatal("migrate: ", err)
	}

	r := httpapi.BuildRouter(d)

	srv := &http.Server{
		Addr:              ":" + port(),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("listening on", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}
