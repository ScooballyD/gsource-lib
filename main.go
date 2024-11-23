package main

import (
	"fmt"

	"github.com/ScooballyD/gsource-lib/scrapers"
	"github.com/ScooballyD/gsource-lib/server"
)

func main() {
	fmt.Println("========starting========") //test

	Agames, err := scrapers.AmzScrape()
	if err != nil {
		fmt.Println("amz error: ", err)
	}
	fmt.Println("Got AMZ games")
	Egames, err := scrapers.EpicScrape()
	if err != nil {
		fmt.Println("epic error: ", err)
	}
	fmt.Println("Got epic games")

	var games []scrapers.Game
	//fmt.Println("===AMZ Games===")
	//for _, game := range Agames {
	//fmt.Printf("Title: %v, URL: %v\n", game.Title, "https://gaming.amazon.com"+game.Href)
	games = append(games, Agames...)
	games = append(games, Egames...)
	//}
	//fmt.Println("===Epic Games===")
	//for _, game := range Egames {
	//fmt.Printf("Title: %v, URL: %v\n", game.Title, game.Href)
	//}
	deals, err := scrapers.EpicDeals()
	if err != nil {
		fmt.Println("Unable to collect deals: ", err)
	}
	fmt.Println("Epic Deals: ", deals)

	server.StartServer(games)
}
