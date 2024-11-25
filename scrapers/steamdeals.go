package scrapers

import (
	"context"
	"fmt"
	"time"

	"github.com/ScooballyD/gsource-lib/internal/database"
	"github.com/chromedp/chromedp"
)

func SteamDeals(db *database.Queries) error {
	fmt.Println("Getting deals")
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
		chromedp.Navigate("https://store.steampowered.com/specials?flavor=contenthub_topsellers"),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var bodySize int
			for i := 0; i < 10; i++ {
				if err := chromedp.Evaluate(`window.scrollBy(0, window.innerHeight)`, nil).Do(ctx); err != nil {
					return err
				}
				if err := chromedp.Sleep(1 * time.Second).Do(ctx); err != nil {
					return err
				}
				if err := chromedp.Evaluate(`document.body.scrollHeight`, &bodySize).Do(ctx); err != nil {
					return err
				}

				var newHeight int
				if err := chromedp.Evaluate(`document.body.scrollHeight`, &newHeight).Do(ctx); err != nil {
					return err
				}
				if newHeight <= bodySize {
					break
				}
				bodySize = newHeight
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			const buttonSel = `button._2tkiJ4VfEdI9kq1agjZyNz`
			for i := 0; i < 3; i++ {
				var isVisible bool
				if err := chromedp.Evaluate(`document.querySelector('`+buttonSel+`') !== null`, &isVisible).Do(ctx); err != nil {
					return err
				}

				if isVisible {
					if err := chromedp.Click(buttonSel, chromedp.NodeVisible).Do(ctx); err != nil {
						return err
					}
					chromedp.Sleep(1 * time.Second).Do(ctx)
				} else {
					break
				}
			}
			return nil
		}),
		chromedp.Evaluate(`
			(function() {
				// Now retrieve data
				const elements = document.querySelector('div._3EdZTDIisUpowxwm6uJ7Iq');
				if (elements) {
					return Array.from(elements.querySelectorAll('div.ImpressionTrackedElement'))
						.map(div => ({
							href: div.querySelector('a').getAttribute('href'),
							title: div.querySelector('img._2eQ4mkpf4IzUp1e9NnM2Wr').getAttribute('alt'),
							image: div.querySelector('img._2eQ4mkpf4IzUp1e9NnM2Wr').getAttribute('src'),
							discount: div.querySelector('div.cnkoFkzVCby40gJ0jGGS4').textContent.trim(),
							og_price: div.querySelector('div._3fFFsvII7Y2KXNLDk_krOW').textContent.trim(),
							price: div.querySelector('div._3j4dI1yA7cRfCvK8h406OB').textContent.trim()
						}));
				}
				return []; // Return an empty array
		})()
		`, &games),
	)
	if err != nil {
		return err
	}

	for _, game := range games {
		_, err = db.AddDiscount(
			context.Background(),
			database.AddDiscountParams{
				Title:    game.Title,
				Url:      game.Href,
				Image:    game.Image,
				Category: "(steam)",
				Price:    game.Price,
				OgPrice:  game.OGprice,
				Discount: game.Discount,
			},
		)
		if err != nil {
			fmt.Printf("err: %v\n", game.Title) //will change
		}
	}

	return nil
}
