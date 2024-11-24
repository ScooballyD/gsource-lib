package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ScooballyD/gsource-lib/scrapers"
)

func (cfg *apiConfig) freeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching games") //test
	games, err := cfg.db.GetGames(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = template.Must(template.ParseFiles("templates/games.html")).Execute(w, games)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// filler for now
func (cfg *apiConfig) discountsHandler(w http.ResponseWriter, r *http.Request) {
	games, err := cfg.db.GetGames(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = template.Must(template.ParseFiles("templates/discounts.html")).Execute(w, games)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	//if cfg.Platform != "dev" {
	//w.WriteHeader(403)
	//return
	//}
	fmt.Println("Resetting") //testing
	//cfg.fileserverHits.Store(0)
	cfg.db.ResetGames(r.Context())
	scrapers.AmzScrape(cfg.db)
	scrapers.EpicScrape(cfg.db)
} //wil change

func (cfg *apiConfig) collecttHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Collecting") //testing
	scrapers.AmzScrape(cfg.db)
	scrapers.EpicScrape(cfg.db)
}
