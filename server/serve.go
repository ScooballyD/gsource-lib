package server

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/ScooballyD/gsource-lib/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	Platform       string
}

func StartServer(dbQ *database.Queries) {
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQ,
		Platform:       os.Getenv("PLATFORM"),
	}
	go Reset()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", cfg.freeHandler)
	mux.HandleFunc("GET /discounts", cfg.discountsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)     //will change from GET
	mux.HandleFunc("GET /admin/collect", cfg.collecttHandler) //will change from GET

	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	srv.ListenAndServe()
}

func Reset() {
	fmt.Println("starting reset")
	ticker := time.NewTicker(1 * time.Minute)
	client := &http.Client{}

	for range ticker.C {
		//test fmt.Println("initiating")
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/admin/reset", nil)
		if err != nil {
			fmt.Println("unable to make request: ", err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("unable to send request: ", err)
			continue
		}
		resp.Body.Close()
	}
}
