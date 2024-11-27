package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ScooballyD/gsource-lib/internal/database"
	"github.com/ScooballyD/gsource-lib/scrapers"
)

type apiConfig struct {
	fileserverHits atomic.Int32 //may use to keep track of server trafic
	db             *database.Queries
	Platform       string //will be useful if remote reset is implemented
}

func StartServer(dbQ *database.Queries) {
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQ,
		Platform:       os.Getenv("PLATFORM"),
	}
	go ResetGame(dbQ)
	go ResetDiscount(dbQ)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", cfg.homeHandler)
	mux.HandleFunc("GET /discounts", cfg.discountsHandler)
	//mux.HandleFunc("POST /admin/reset", cfg.resetHandler)     //needed for remote reset
	mux.HandleFunc("GET /free", cfg.freeHandler)

	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	srv.ListenAndServe()
}

func ResetGame(dbQ *database.Queries) {
	ticker := time.NewTicker(45 * time.Minute) //can reduce if you want the database to update more often

	for range ticker.C {
		ResetGamebp(dbQ)
	}
}

func ResetGamebp(dbQ *database.Queries) {
	fmt.Println("-Beginning Free Game Reset-")
	games, err := scrapers.EpicScrape(dbQ)
	if err != nil {
		fmt.Println("unable to get EpicGames: ", err)
		return
	}

	amz, err := scrapers.AmzScrape(dbQ)
	if err != nil {
		fmt.Println("unable to get EpicGames: ", err)
		return
	}
	games = append(games, amz...)

	err = dbQ.ResetGames(context.Background())
	if err != nil {
		fmt.Println("Unable to reset Games: ", err)
	}

	for _, game := range games {
		_, err = dbQ.AddGame(
			context.Background(),
			database.AddGameParams{
				Title:    game.Title,
				Url:      game.Href,
				Image:    game.Image,
				Category: game.Category,
			},
		)
		if err != nil {
			fmt.Printf("err: %v\n", game.Title)
		}
	}
	fmt.Println("-Finished Free Game Reset-")
}

func ResetDiscount(dbQ *database.Queries) {
	ticker := time.NewTicker(2 * time.Hour) //can reduce if needed

	for range ticker.C {
		ResetDiscountbp(dbQ)
	}
}

func ResetDiscountbp(dbQ *database.Queries) {
	fmt.Println("-Beginning Discount Game Reset-")
	games, err := scrapers.GogDeals(dbQ)
	if err != nil {
		fmt.Println("unable to get GOGDeals: ", err)
		return
	}

	steam, err := scrapers.SteamDeals(dbQ)
	if err != nil {
		fmt.Println("unable to get SteamDeals: ", err)
		return
	}
	games = append(games, steam...)

	//humble, err := scrapers.HumbleDeals(dbQ)
	//if err != nil {
	//fmt.Println("unable to get HumbleDeals: ", err) // work in progress
	//return
	//}
	//games = append(games, humble...)

	err = dbQ.ResetDiscounts(context.Background())
	if err != nil {
		fmt.Println("Unable to reset Discounts: ", err)
	}

	for _, game := range games {
		price, err := strconv.ParseFloat(strings.Trim(game.Price, "$"), 64) //needed to change to float so sorting query worked correctly
		if err != nil {
			fmt.Println("Unable to convert price: ", err)
			return
		}
		_, err = dbQ.AddDiscount(
			context.Background(),
			database.AddDiscountParams{
				Title:    game.Title,
				Url:      game.Href,
				Image:    game.Image,
				Category: game.Category,
				Price:    price,
				OgPrice:  game.OGprice,
				Discount: game.Discount,
				Rating:   game.Rating,
			},
		)

		if err != nil {
			fmt.Printf("err: %v\n", game.Title)
		}
	}
	fmt.Println("-Finished Discount Game Reset-")
}
