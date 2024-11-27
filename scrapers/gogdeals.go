package scrapers

import (
	"context"
	"fmt"
	"time"

	"github.com/ScooballyD/gsource-lib/internal/database"
	"github.com/chromedp/chromedp"
)

func GogDeals(db *database.Queries) ([]Game, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 8*time.Minute)
	defer cancel()

	var games []Game
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.gog.com/en/games?discounted=true"),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 10; i++ {
				if err := chromedp.Evaluate(`
					window.scrollTo({
						top: document.body.scrollHeight,
						behavior: 'smooth'
					});
				`, nil).Do(ctx); err != nil {
					return fmt.Errorf("failed to scroll: %v", err)
				}
				if err := chromedp.Sleep(1 * time.Second).Do(ctx); err != nil {
					return fmt.Errorf("failed to sleep: %v", err)
				}
			}
			if err := chromedp.Evaluate(`window.scrollTo(0, 0);`, nil).Do(ctx); err != nil {
				return fmt.Errorf("failed to scroll to top: %v", err)
			}

			return nil
		}),
		chromedp.Evaluate(`
			Array.from(document.querySelector('div.paginated-products-grid.grid')
                .querySelectorAll('a'))
                .map(a => ({
                    href: a.getAttribute('href'),
					title: a.querySelector('img.ng-star-inserted').getAttribute('alt').trim().replace(' - cover art image', ''),
					image: a.querySelector('source.ng-star-inserted').getAttribute('srcset').split(',')[1].trim().replace(' 2x', ''),
					discount: a.querySelector('price-discount.ng-star-inserted').textContent.trim(),
					og_price: a.querySelector('span.base-value.ng-star-inserted').textContent.trim(),
					price: a.querySelector('span.final-value.ng-star-inserted').textContent.trim()
                }))
		`, &games),
	)
	if err != nil {
		return nil, err
	}

	//This section by far uses most of the scraping time.  It can be cut if ratings are not desired
	for i := range games {
		fmt.Println("Getting rating for: ", games[i].Title)
		var rating string
		err = chromedp.Run(ctx,
			chromedp.Navigate(games[i].Href),
			chromedp.Sleep(1*time.Second),
			chromedp.Evaluate(`
            const ratingElement = document.querySelector('div.rating.productcard-rating__score');
            	if (ratingElement) {
                	ratingElement.textContent.trim();
            	} else {
                	"N/A";
            	}
        	`, &rating),
		)
		if err != nil {
			return nil, err
		}
		games[i].Rating = rating
		games[i].Category = "GoG"
	}

	return games, nil
}
