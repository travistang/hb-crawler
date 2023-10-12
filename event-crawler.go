package main

import (
	hb "hb-crawler/rating-gain/hiking-buddies"
	"net/http"

	"github.com/corpix/uarand"
	"github.com/go-resty/resty/v2"
)

type EventListResponse map[string]*hb.Event

func CrawlEventList(sessionId *string) (error, *EventListResponse) {
	client := resty.New()
	client.SetCookie(&http.Cookie{
		Name:  "sessionid",
		Value: *sessionId,
	})
	client.SetHeaders(map[string]string{
		"Accept":     "application/json",
		"User-Agent": uarand.GetRandom(),
	})

	// TODO: Complete this
	var responseData EventListResponse
	res, err := client.R().
		EnableTrace().
		SetResult(&responseData).
		Get(string(hb.EventListEndpoint))

	if err != nil {
		return err, nil
	}
	return nil, nil
}
