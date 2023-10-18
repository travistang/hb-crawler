package hiking_buddies

import (
	"fmt"
)

type EventListResponse map[string][]Event
type PastEventListResponse struct {
	Results []Event
}

func doFetch[R interface{}](url URL, cookieCredential *CookieCredential, result *R) (*R, error) {
	client := prepareRestClient(cookieCredential)
	res, err := client.R().
		EnableTrace().
		SetResult(&result).
		Get(string(url))

	if err != nil {
		return nil, err
	}

	statusCode := res.StatusCode()
	if statusCode != 200 {
		err := fmt.Errorf("expected status code 200, got %d instead", statusCode)
		return nil, err
	}

	return result, nil
}

func FetchUpcomingEvents(cookieCredential *CookieCredential) (*EventListResponse, error) {
	return doFetch(EventListEndpoint, cookieCredential, &EventListResponse{})
}

func FetchPastEvents(cookieCredential *CookieCredential) (*PastEventListResponse, error) {
	return doFetch(PastEventListEndpoint, cookieCredential, &PastEventListResponse{})
}
