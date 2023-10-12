package hiking_buddies

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/corpix/uarand"
	"github.com/go-resty/resty/v2"
)

type EventListResponse map[string][]Event

func CrawlEventList(cookieCredential *CookieCredential) (error, *EventListResponse) {
	client := resty.New()
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

	var responseData EventListResponse
	res, err := client.R().
		EnableTrace().
		SetResult(&responseData).
		Get(string(EventListEndpoint))

	if err != nil {
		return err, nil
	}

	statusCode := res.StatusCode()
	if statusCode != 200 {
		return errors.New(fmt.Sprintf("Expected status code 200, got %d instead", statusCode)), nil
	}

	return nil, &responseData
}
