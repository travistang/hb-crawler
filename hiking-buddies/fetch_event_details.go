package hiking_buddies

import (
	"context"
	"fmt"
	"strconv"

	"hb-crawler/rating-gain/logging"
	"hb-crawler/rating-gain/utils"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

const (
	eventSummaryBoxSelector = "//*[contains(@class, 'event-dash-card')]"
)

func getLogger() *logrus.Logger {
	return logging.GetLogger(&logging.LoggerConfig{
		Prefix: "Fetch event details",
		Level:  logrus.DebugLevel,
	})
}

type FetchEventParticipantsParams struct {
	Id         int
	Credential *CookieCredential
}

func extractParticipantIds(p *FetchEventParticipantsParams, log *logrus.Logger, ids *[]string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Evaluate(`
				[...document.querySelectorAll('#id_participants_list_container a')].map((e) => e.href.split('/').pop())
			`,
			ids),
	}
}

func fetchEventParticipants(p *FetchEventParticipantsParams, log *logrus.Logger, ids *[]string) chromedp.Tasks {
	url := fmt.Sprintf("%s%d", EventDetailsEndpoint, p.Id)

	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(createHeaders()),
		setCookies(p.Credential),

		chromedp.Navigate(url),

		chromedp.ActionFunc(func(ctx context.Context) error { fmt.Printf("Navigated to %s", url); return nil }),
		chromedp.WaitVisible(eventSummaryBoxSelector),
		extractParticipantIds(p, log, ids),
	}
}

func processFetchedIds(idsString *[]string) *[]int {
	ids := []int{}
	for _, idString := range *idsString {
		id, err := strconv.Atoi(idString)
		if err != nil || utils.Find(&ids, &id) != -1 {
			continue
		}
		ids = append(ids, id)
	}
	return &ids
}

func FetchEventParticipants(p *FetchEventParticipantsParams) (*[]int, error) {
	log := getLogger()
	log.Debugf("Fetching event details for ID %d", p.Id)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var idsString []string
	err := chromedp.Run(ctx, fetchEventParticipants(p, log, &idsString))

	if err != nil {
		return nil, err
	}

	ids := processFetchedIds(&idsString)
	log.Debugf("Found %d ids in the event list", len(*ids))

	return ids, nil
}
