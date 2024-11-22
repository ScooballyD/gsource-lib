package scrapers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []EpicGame `json:"elements"`
			} `json:"searchStore"`
		} `json:"Catalog"`
	} `json:"data"`
}

type EpicGame struct {
	Title       string     `json:"title"`
	ID          string     `json:"id"`
	Namespace   string     `json:"namespace"`
	Description string     `json:"description"`
	KeyImages   []Image    `json:"keyImages"`
	OfferMap    []Offer    `json:"offerMappings"`
	Price       Price      `json:"price"`
	Categories  []Category `json:"categories"`
	Promotions  Promotions `json:"promotions"`
}

type Offer struct {
	PageSlug string `json:"pageSlug"`
}

type Promotions struct {
	PromotionalOffers         []PromotionalOffer `json:"promotionalOffers"`
	UpcomingPromotionalOffers []PromotionDetails `json:"upcomingPromotionalOffers"`
}

type PromotionalOffer struct {
	PromotionalOffers []PromotionDetails `json:"promotionalOffers"`
}

type PromotionDetails struct {
	StartDate       string          `json:"startDate"`
	EndDate         string          `json:"endDate"`
	DiscountSetting DiscountSetting `json:"discountSetting"`
}

type DiscountSetting struct {
	DiscountType       string `json:"discountType"`
	DiscountPercentage int    `json:"discountPercentage"`
}

type Image struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Price struct {
	TotalPrice struct {
		DiscountPrice float64 `json:"discountPrice"`
		OriginalPrice float64 `json:"originalPrice"`
		CurrencyCode  string  `json:"currencyCode"`
		FmtPrice      struct {
			OriginalPrice string `json:"originalPrice"`
			DiscountPrice string `json:"discountPrice"`
		} `json:"fmtPrice"`
	} `json:"totalPrice"`
}

type Category struct {
	Path string `json:"path"`
}

func EpicScrape() ([]Game, error) {
	resp, err := http.Get("https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=en-US&country=US&allowCountries=US")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println(err)
		return nil, err
	}
	var games []Game
	for _, game := range response.Data.Catalog.SearchStore.Elements {
		if game.Price.TotalPrice.DiscountPrice == 0 {
			games = append(games, Game{
				Title: game.Title,
				Href:  "https://store.epicgames.com/en-US/p/" + game.OfferMap[0].PageSlug,
			})
			//fmt.Println("added: ", game.Title)
			//fmt.Println("offerMap: ", "https://store.epicgames.com/en-US/p/"+game.OfferMap[0].PageSlug)
		}
	}

	//fmt.Println(games)
	return games, nil
}
