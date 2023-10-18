package hiking_buddies

import (
	"context"
	"fmt"
	"hb-crawler/rating-gain/database"
	"hb-crawler/rating-gain/logging"
	"strconv"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

const (
	RoutePointsSelector = "//*[@class='map-statistics']//tr[./th[contains(text(), 'Rating')]]/td[last()]"
)

func localAndParsePoints(points *int) chromedp.Tasks {
	var pointsString string
	return chromedp.Tasks{
		chromedp.Text(pointsSelector, &pointsString),
		chromedp.ActionFunc(func(ctx context.Context) error {
			parsedPoint, err := strconv.Atoi(pointsString)
			if err != nil {
				return err
			}
			*points = parsedPoint
			return nil
		}),
	}
}

func fetchRoutePoints(id int, credential *CookieCredential, points *int) chromedp.Tasks {
	url := fmt.Sprintf("%s%d/", string(RouteDetailsEndpoint), id)
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(createHeaders()),
		setCookies(credential),

		chromedp.Navigate(url),

		chromedp.WaitVisible(RoutePointsSelector),
		localAndParsePoints(points),
	}
}

func FetchRoutePoints(id int, credential *CookieCredential) (*int, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	var points int
	defer cancel()

	err := chromedp.Run(ctx, fetchUserPoints(credential, id, &points))
	if err != nil {
		return nil, err
	}

	return &points, nil
}

type GetRoutePointsParams struct {
	Repo       *database.RouteRepository
	Id         int
	Record     *database.RouteRecord
	Credential *CookieCredential
}

func GetRoutePoints(p *GetRoutePointsParams) error {
	log := logging.GetLogger(&logging.LoggerConfig{
		Prefix: "Fetch routes",
		Level:  logrus.DebugLevel,
	})

	log.Debug("Getting route points by cache...")

	var fetchedRoute database.RouteRecord
	err := p.Repo.GetRouteById(p.Id, &fetchedRoute)
	if err != nil {
		return err
	}

	if fetchedRoute.Id != 0 && fetchedRoute.Points != nil {
		log.Debugf("Route points found in cache: route %d has %d points", fetchedRoute.Id, *fetchedRoute.Points)
		*p.Record = fetchedRoute
		return nil
	}

	log.Debugf("Route %d's point is not cached, fetching...", fetchedRoute.Id)

	points, err := FetchRoutePoints(p.Id, p.Credential)
	if err != nil {
		log.Errorf("Failed to fetch points for Route %d: %+v\n", p.Id, err)
		return err
	}

	log.Infof("Route %d has %d points (found by fetching) \n", p.Id, *points)
	log.Infof("Caching route %d to database", p.Id)
	if _, err := p.Repo.SaveRoute(&database.RouteRecord{
		Id:        p.Id,
		Elevation: p.Record.Elevation,
		Points:    points,
		Distance:  p.Record.Distance,
		Name:      p.Record.Name,
		Scale:     p.Record.Scale,
	}); err != nil {
		log.Errorf("Failed to cache route %d: %+v\n", p.Id, err)
	}

	*(p.Record) = fetchedRoute
	p.Record.Points = points

	return nil
}
