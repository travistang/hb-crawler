package hiking_buddies

import (
	"context"
	"net/http"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/corpix/uarand"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func prepareRestClient(cookieCredential *CookieCredential) *resty.Client {
	client := resty.New()
	log.Debugf("Fetching Request with credential: %+v\n", *cookieCredential)
	cookie := cookieCredential.AsCookie()
	for k, v := range cookie {
		client.SetCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	client.SetHeaders(map[string]string{
		"Accept":     "application/json",
		"User-Agent": uarand.GetRandom(),
	})

	return client
}

func baseFetchFunction(url string, cookieCredential *CookieCredential) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(createHeaders()),
		setCookies(cookieCredential),

		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Infof("--> Launching Crawler against url %s", url)
			return nil
		}),

		chromedp.Navigate(url),
	}

}
