package worker

import (
	"hb-crawler/rating-gain/database"
	hb "hb-crawler/rating-gain/hiking-buddies"
	"hb-crawler/rating-gain/logging"

	log "github.com/sirupsen/logrus"
)

func CreateRoutePointsWorker(config *WorkerConfig) *Worker {
	logger := logging.GetLogger(&logging.LoggerConfig{
		Prefix: "Route points worker",
		Level:  log.DebugLevel,
	})

	return &Worker{
		repository:      config.Repository,
		shouldRun:       false,
		interval:        config.Interval,
		logger:          logger,
		LastRunningTime: nil,
		ProcessFunc:     routePointsProcessFunc,
	}
}

func getNextRouteId(context *WorkerProcessContext) (*int, error) {
	worker := context.Worker
	worker.logger.Info("Computing next event to crawl")

	nextId, err := worker.repository.Route.GetNextId()
	if err != nil {
		return nil, err
	}
	return nextId, nil
}

func crawlRouteDetails(context *WorkerProcessContext, routeId int) (*database.RouteRecord, error) {
	worker := context.Worker
	worker.logger.Info("crawling route details")
	record, err := hb.FetchAllRouteDetails(routeId, context.Credential)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func routePointsProcessFunc(context *WorkerProcessContext) error {
	worker := context.Worker
	worker.logger.Info("Crawling routes details")

	crawlEventId, err := getNextRouteId(context)
	if err != nil {
		return err
	}
	if crawlEventId == nil {
		worker.logger.Info("No more routes to be crawled. Skipping")
		return nil
	}
	record, err := crawlRouteDetails(context, *crawlEventId)
	if err != nil {
		worker.logger.Warnf("problem with fetching route details, %d: %+v\n", *crawlEventId, record)
		return err
	}
	if _, err := worker.repository.Route.SaveRoute(record); err != nil {
		worker.logger.Warnf("problem saving route details for id %d: %+v\n", *crawlEventId, err)
		return err
	}

	return nil
}
