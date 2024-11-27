package server

import (
	"html/template"
	"net/http"
)

func (cfg *apiConfig) freeHandler(w http.ResponseWriter, r *http.Request) {
	games, err := cfg.db.GetGames(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = template.Must(template.ParseFiles("templates/games.html")).Execute(w, games)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) discountsHandler(w http.ResponseWriter, r *http.Request) {
	sortParam := r.URL.Query().Get("sort")

	var discounts any
	var err error

	switch sortParam { //used to sort database results depending on sort query
	case "price":
		discounts, err = cfg.db.GetDiscountsPrice(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "title":
		discounts, err = cfg.db.GetDiscountsTitle(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "discount":
		discounts, err = cfg.db.GetDiscountsDiscount(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	default:
		discounts, err = cfg.db.GetDiscounts(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	err = template.Must(template.ParseFiles("templates/discounts.html")).Execute(w, discounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (cfg *apiConfig) homeHandler(w http.ResponseWriter, r *http.Request) {
	game, err := cfg.db.GetEpicGame(r.Context()) //may put epic game on home page since its the only truly free one
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = template.Must(template.ParseFiles("templates/home.html")).Execute(w, game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {   //may include remote reset capabilities
//if cfg.Platform != "dev" {
//w.WriteHeader(403)
//return
//}
//fmt.Println("Resetting") //testing
//cfg.fileserverHits.Store(0)
//cfg.db.ResetGames(r.Context())
//scrapers.EpicScrape(cfg.db)
//scrapers.AmzScrape(cfg.db)
//cfg.db.ResetDiscounts(r.Context())
//scrapers.SteamDeals(cfg.db)
//} //wil change

//func (cfg *apiConfig) collecttHandler(w http.ResponseWriter, r *http.Request) {
//fmt.Println("Collecting") //testing
//scrapers.AmzScrape(cfg.db)
//scrapers.EpicScrape(cfg.db)
//}
