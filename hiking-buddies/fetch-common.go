package hiking_buddies

import (
	"net/http"

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
