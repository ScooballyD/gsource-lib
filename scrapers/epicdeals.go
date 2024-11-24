package scrapers

type Deal struct {
	Title    string  `json:"title"`
	Href     string  `json:"href"`
	Image    string  `json:"image"`
	Category string  `json:"catagory"`
	Discount float64 `json:"discount"`
}

// work in prog
func EpicDeals() ([]Deal, error) {
	response, err := EpicHelper("https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=en-US&country=US&allowCountries=US")
	if err != nil {
		return nil, err
	}

	var deals []Deal
	for _, game := range response.Data.Catalog.SearchStore.Elements {
		if game.Price.TotalPrice.DiscountPrice == 0 {
			deals = append(deals, Deal{
				Title:    game.Title,
				Href:     "https://store.epicgames.com/en-US/p/" + game.OfferMap[0].PageSlug,
				Image:    game.KeyImages[0].URL,
				Category: "(epic)",
				Discount: 1 - (game.Price.TotalPrice.DiscountPrice / game.Price.TotalPrice.OriginalPrice),
			})
		}
	}
	return deals, nil
}
