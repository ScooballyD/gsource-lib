package scrapers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ScooballyD/gsource-lib/internal/database"
)

type Response struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []EpicGame `json:"elements"`
			} `json:"searchStore"`
		} `json:"Catalog"`
	} `json:"data"`
}

type EpicGame struct {
	Title       string      `json:"title"`
	ID          string      `json:"id"`
	Namespace   string      `json:"namespace"`
	Description string      `json:"description"`
	KeyImages   []Image     `json:"keyImages"`
	OfferMap    []Offer     `json:"offerMappings,omitempty"`
	Price       *Price      `json:"price,omitempty"`
	Categories  []Category  `json:"categories,omitempty"`
	Promotions  *Promotions `json:"promotions,omitempty"`
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

func EpicHelper(url string) (Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("Host", "store.epicgames.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:132.0) Gecko/20100101 Firefox/132.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://store.epicgames.com/en-US/browse?sortBy=releaseDate&sortDir=DESC&priceTier=tierDiscouted&category=Game&count=40&start=0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DisableCompression: false,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&response); err != nil {
		fmt.Println("decode error:", err)
		return Response{}, err
	}

	return response, nil
}

func EpicScrape(db *database.Queries) error {
	response, err := EpicHelper("https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=en-US&country=US&allowCountries=US")
	if err != nil {
		return err
	}

	for _, game := range response.Data.Catalog.SearchStore.Elements {
		if game.Price.TotalPrice.DiscountPrice == 0 {
			_, err = db.AddGame(
				context.Background(),
				database.AddGameParams{
					Title:    game.Title,
					Url:      "https://store.epicgames.com/en-US/p/" + game.OfferMap[0].PageSlug,
					Image:    game.KeyImages[0].URL,
					Category: "(epic)",
				},
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
