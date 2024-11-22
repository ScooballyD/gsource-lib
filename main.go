package main

import (
	"fmt"

	"github.com/ScooballyD/gsource-lib/scrapers"
)

func main() {
	fmt.Println("========starting========") //test
	//c := NewClient(5 * time.Minute)

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
	fmt.Println("===AMZ Games===")
	for _, game := range Agames {
		fmt.Printf("Title: %v, URL: %v\n", game.Title, "https://gaming.amazon.com"+game.Href)
	}
	fmt.Println("===Epic Games===")
	for _, game := range Egames {
		fmt.Printf("Title: %v, URL: %v\n", game.Title, game.Href)
	}
}
