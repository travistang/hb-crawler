package hiking_buddies

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	pointsSelector = "//*/div[@class='karma']"
)

func setCookies(credential *CookieCredential) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			return network.SetCookie(loginCookieName, credential.SessionId).
				WithHTTPOnly(true).
				WithDomain(string(MainDomainWithoutProtocol)).
				WithPath("/").
				Do(ctx)
		}),
	}
}

func parsePointString(pointsString string) (*int, error) {
	re := regexp.MustCompile("[0-9]+")
	pointStrings := re.FindAllString(pointsString, 1)
	if len(pointStrings) != 1 {
		return nil, fmt.Errorf("failed to parse points string %s", pointsString)
	}
	points, err := strconv.Atoi(pointStrings[0])
	if err != nil {
		return nil, err
	}
	return &points, nil
}

func locateAndParsePoints(points *int) chromedp.Tasks {
	var pointsString string
	return chromedp.Tasks{
		chromedp.Text(pointsSelector, &pointsString),
		chromedp.ActionFunc(func(ctx context.Context) error {
			parsedPoint, err := parsePointString(pointsString)
			if err != nil {
				return err
			}
			*points = *parsedPoint
			return nil
		}),
	}
}

func fetchUserPoints(credential *CookieCredential, id int, points *int) chromedp.Tasks {
	url := fmt.Sprintf("%s%d/", string(UserDetailsEndpoint), id)

	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(createHeaders()),
		setCookies(credential),

		chromedp.Navigate(url),

		chromedp.WaitVisible(pointsSelector),
		locateAndParsePoints(points),
	}
}

func FetchUserPoints(id int, credential *CookieCredential) (*int, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	var points int
	defer cancel()

	err := chromedp.Run(ctx, fetchUserPoints(credential, id, &points))
	if err != nil {
		return nil, err
	}

	return &points, nil
}
