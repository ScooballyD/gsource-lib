package scrapers

import (
	"context"
	"time"

	"github.com/ScooballyD/gsource-lib/internal/database"
	"github.com/chromedp/chromedp"
)

type Game struct {
	Title    string `json:"title"`
	Href     string `json:"href"`
	Image    string `json:"image"`
	Category string `json:"category"`
	Discount string `json:"discount"`
	OGprice  string `json:"og_price"`
	Price    string `json:"price"`
}

func AmzScrape(db *database.Queries) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var games []Game
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://gaming.amazon.com/intro"),
		chromedp.Sleep(1*time.Second),
		chromedp.Evaluate(`
			Array.from(document.querySelector('div[data-a-target="offer-section-offer-cards"]')
                .querySelectorAll('a'))
                .map(a => ({
                    href: "https://gaming.amazon.com" + a.getAttribute('href'),
					title: a.querySelector('img.tw-image').getAttribute('alt'),
					image: a.querySelector('img.tw-image').getAttribute('srcset').split(',')[0].trim().replace(' 1x', '')

                }))
		`, &games),
	)
	if err != nil {
		return err
	}

	for _, game := range games {
		_, err = db.AddGame(
			context.Background(),
			database.AddGameParams{
				Title:    game.Title,
				Url:      game.Href,
				Image:    game.Image,
				Category: "(prime)",
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
