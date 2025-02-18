package hiking_buddies

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"hb-crawler/rating-gain/database"
	"hb-crawler/rating-gain/logging"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

const (
	RoutePointsSelector   = `document.querySelector("table.table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(6) > td:nth-child(8)").innerText`
	ElevationGainSelector = `document.querySelector('table.table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(8)').innerText`
	ElevationLossSelector = `document.querySelector('table.table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(3) > td:nth-child(8)').innerText`
)

func getNumberFromSelector(selector string, number *int, err *error) chromedp.Tasks {
	var plainValue string
	numExtractRegexp := regexp.MustCompile("[0-9.]+")
	return chromedp.Tasks{
		chromedp.Evaluate(selector, &plainValue),
		chromedp.ActionFunc(func(ctx context.Context) error {
			matches := numExtractRegexp.FindStringSubmatch(plainValue)
			if len(matches) == 0 {
				*err = fmt.Errorf("no numbers found under selector %s", selector)
				return *err
			}
			num, convertErr := strconv.Atoi(matches[0])
			if convertErr != nil {
				*err = convertErr
				return *err
			}

			*number = num
			return nil
		}),
	}
}

func locateAndParsePointsForRoute(routeId int, points *int) chromedp.Tasks {
	var pointsString string
	return chromedp.Tasks{
		chromedp.Evaluate(RoutePointsSelector, &pointsString),
		chromedp.ActionFunc(func(ctx context.Context) error {
			parsedPoint, err := strconv.Atoi(pointsString)
			if err != nil {
				log.Warnf(
					"Failed to retrieve points for route %d, selector got '%s' instead",
					routeId, pointsString,
				)
				return err
			}
			*points = parsedPoint
			return nil
		}),
	}
}

func locateAndParseElevations(data *database.RouteRecord) chromedp.Tasks {
	var elevationGain, elevationLoss *int
	var err *error
	return chromedp.Tasks{
		getNumberFromSelector(ElevationGainSelector, elevationGain, err),
		getNumberFromSelector(ElevationLossSelector, elevationLoss, err),
	}
}

func fetchRoutePoints(id int, credential *CookieCredential, points *int) chromedp.Tasks {
	log.Infof("start crawling points for route %d", id)
	url := fmt.Sprintf("%s%d/", string(RouteDetailsEndpoint), id)
	return chromedp.Tasks{
		baseFetchFunction(url, credential),
		locateAndParsePointsForRoute(id, points),
	}
}

func fetchAllRouteDetails(id int, credential *CookieCredential, data *database.RouteRecord) chromedp.Tasks {
	log.Infof("start crawling all details for route %d", id)
	url := fmt.Sprintf("%s%d/", string(RouteDetailsEndpoint), id)
	return chromedp.Tasks{
		baseFetchFunction(url, credential),
		locateAndParsePointsForRoute(id, data.Points),
		locateAndParseElevations(data),
	}
}

func FetchRoutePoints(id int, credential *CookieCredential) (*int, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	var points int
	defer cancel()

	err := chromedp.Run(ctx, fetchRoutePoints(id, credential, &points))
	if err != nil {
		return nil, err
	}

	return &points, nil
}

func FetchAllRouteDetails(id int, credential *CookieCredential) (*database.RouteRecord, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var data database.RouteRecord
	err := chromedp.Run(ctx, fetchAllRouteDetails(id, credential, &data))

	if err != nil {
		return nil, err
	}
	return &data, nil
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
		Level:  log.DebugLevel,
	})

	log.Debugf("Getting points for route %d by cache...", p.Id)

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

	log.Debugf("Route %d's point is not cached, fetching...", p.Id)

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
