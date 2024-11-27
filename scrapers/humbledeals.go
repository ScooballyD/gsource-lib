package scrapers

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ScooballyD/gsource-lib/internal/database"
	"github.com/andybalholm/brotli"
	"github.com/chromedp/chromedp"
)

type HResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	RequireLinkedThirdPartyAccountWhen string                 `json:"require_linked_third_party_account_when"`
	StandardCarouselImage              string                 `json:"standard_carousel_image"`
	DeliveryMethods                    []string               `json:"delivery_methods"`
	MachineName                        string                 `json:"machine_name"`
	XrayTraitsThumbnail                string                 `json:"xray_traits_thumbnail"`
	ContentTypes                       []string               `json:"content_types"`
	HumanURL                           string                 `json:"human_url"`
	Platforms                          []string               `json:"platforms"`
	IconDict                           map[string]Icon        `json:"icon_dict"`
	FeaturedImageRecommendation        string                 `json:"featured_image_recommendation"`
	LargeCapsule                       string                 `json:"large_capsule"`
	HumanName                          string                 `json:"human_name"`
	RequiredAccountLinks               []string               `json:"required_account_links"`
	Type                               string                 `json:"type"`
	Icon                               string                 `json:"icon"`
	NonRewardsCharitySplit             float64                `json:"non_rewards_charity_split"`
	RewardsSplit                       float64                `json:"rewards_split"`
	EmptyTpkds                         map[string]interface{} `json:"empty_tpkds"`
	FullPrice                          Cost                   `json:"full_price"`
	SaleEnd                            int64                  `json:"sale_end"`
	Nonrefundable                      bool                   `json:"nonrefundable"`
	CurrentPrice                       Cost                   `json:"current_price"`
	RatingForCurrentRegion             string                 `json:"rating_for_current_region"`
	SaleType                           string                 `json:"sale_type"`
}

type Icon struct {
	Available   []string `json:"available"`
	Unavailable []string `json:"unavailable"`
}

type Cost struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

// Work in progress
func HumbleHelper(url string) (HResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return HResponse{}, err
	}

	req.Header.Set("Host", "www.humblebundle.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:132.0) Gecko/20100101 Firefox/132.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.humblebundle.com/store/search?sort=bestselling&filter=onsale&_gl=1*hjfigu*_up*MQ..*_gs*MQ..&gclid=EAIaIQobChMIzt2q_Ob4iQMV5lN_AB1e0i2tEAAYASAAEgLK3_D_BwE&results=80&page=1")
	req.Header.Set("Cookie", "utmccmpn=04_1000_HumbleBundle_USA_Search_Brand; utmadid=694915429345; utmcmed=cpc; utmtrm=\"humble bundle\"; utmcntnt=694915429345; utmkywrdid=kwd-18642651705; utmcsr=google; csrf_cookie=EZbdDGY2MUkrDBbz-1-1732583091; _simpleauth_sess=eyJpZCI6InNvRVBzVGJNcDIifQ==|1732650706|a9db88379fc3349b261c65545547fa8ac301611b; OptanonConsent=isGpcEnabled=0&datestamp=Tue+Nov+26+2024+13%3A51%3A18+GMT-0600+(Central+Standard+Time)&version=202409.1.0&browserGpcFlag=0&isIABGlobal")

	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DisableCompression: false,
		},
	}

	resp, err := client.Do(req)
	fmt.Println("status: ", resp.Status)
	if err != nil {
		return HResponse{}, err
	}
	defer resp.Body.Close()
	// tester
	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return HResponse{}, err
		}
		defer resp.Body.Close()
	case "br":
		reader = brotli.NewReader(resp.Body)
	case "deflate":
		reader = flate.NewReader(resp.Body)
		defer resp.Body.Close()
	default:
		reader = resp.Body
	}

	rawBody, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return HResponse{}, err
	}

	fmt.Println(string(rawBody))
	// tester

	var response HResponse
	if err := json.NewDecoder(bytes.NewReader(rawBody)).Decode(&response); err != nil {
		fmt.Println("decode error:", err)
		return HResponse{}, err
	}

	return response, nil
}

// Work in progress
func HumbleDeals(db *database.Queries) ([]Game, error) {
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
		chromedp.Navigate("https://www.humblebundle.com/store"),
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
		//<div class="_2nuoOi5kC2aUI12z85PneA">Very Positive</div>
		chromedp.Evaluate(`
			(function() {
				const elements = document.querySelector('div.chunks-container');
				if (elements) {
					return Array.from(elements.querySelectorAll('div.entity.js-entity.on-sale'))
						.map(div => ({
							//href: div.querySelector('a').getAttribute('href'),
							//title: div.querySelector('img._2eQ4mkpf4IzUp1e9NnM2Wr').getAttribute('alt'),
							image: div.querySelector('img.entity-image').getAttribute('src'),
							//discount: div.querySelector('div.cnkoFkzVCby40gJ0jGGS4').textContent.trim(),
							//og_price: div.querySelector('div._3fFFsvII7Y2KXNLDk_krOW').textContent.trim(),
							//price: div.querySelector('div._3j4dI1yA7cRfCvK8h406OB').textContent.trim(),
							rating: "N/A",
							category: "Steam"
						}));
				}
				return [];
		})()
		`, &games),
	)
	if err != nil {
		return nil, err
	}
	fmt.Println("Games: ", games)

	return games, nil
}

//<a class="entity-link js-entity-link" href="/store/beamngdrive" aria-label="BeamNG.drive">
