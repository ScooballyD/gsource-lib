package scrapers

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

type Game struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

func AmzScrape() ([]Game, error) {
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
                    href: a.getAttribute('href'),
					title: a.querySelector('img.tw-image').getAttribute('alt')
                }))
		`, &games),
	)
	if err != nil {
		return nil, err
	}
	// for testing
	//fmt.Println("live scrape finishing")
	//for _, game := range games {
	//fmt.Printf("Label: %s, URL: %s\n", game.Title, game.Href)
	//}
	return games, nil
}
