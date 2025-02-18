package worker

import (
	"time"

	"hb-crawler/rating-gain/database"
	hb "hb-crawler/rating-gain/hiking-buddies"
	"hb-crawler/rating-gain/logging"

	log "github.com/sirupsen/logrus"
)

func CreatePastEventWorker(config *WorkerConfig) *Worker {
	logger := logging.GetLogger(&logging.LoggerConfig{
		Prefix: "past event worker",
		Level:  log.DebugLevel,
	})

	worker := Worker{
		repository:      config.Repository,
		shouldRun:       false,
		interval:        config.Interval,
		logger:          logger,
		LastRunningTime: nil,
		ProcessFunc:     pastEventProcessFunc,
	}

	return &worker
}

func processRecentPastEvent(
	context *WorkerProcessContext,
	event *hb.Event,
) (bool, error) {
	worker := context.Worker
	worker.logger.Infof("Treating event %s as a recent past event", event.Title)

	pointGains, _ := worker.repository.PointGains.GetPointGainsByEventId(event.ID)
	if pointGains != nil && len(*pointGains) > 0 {
		worker.logger.Infof("Refuse to process recent past event as %d related point gains event is found", len(*pointGains))
		return false, nil
	}

	worker.logger.Infof("Fetching route points for route '%s' under event '%s'", event.Route.RouteTitle, event.Title)
	var routeRecord database.RouteRecord
	if err := hb.GetRoutePoints(&hb.GetRoutePointsParams{
		Repo:       worker.repository.Route,
		Id:         event.Route.RouteID,
		Record:     &routeRecord,
		Credential: context.Credential,
	}); err != nil {
		return false, nil
	}

	worker.logger.Infof("Fetching participant lists of event '%s'", event.Title)
	ids, err := hb.FetchEventParticipants(&hb.FetchEventParticipantsParams{
		Id:         event.ID,
		Credential: context.Credential,
	})
	if err != nil {
		return false, nil
	}

	for _, userId := range *ids {
		worker.logger.Infof("Fetching current points for user %d", userId)
		currentUserPoints, err := hb.FetchUserPoints(userId, context.Credential)
		if err != nil {
			worker.logger.Warnf("Failed to fetch current points for user %d", userId)
			continue
		}
		if err := worker.repository.PointGains.CreatePointsGainEntry(&database.PointGainRecord{
			UserId:           userId,
			UserPointsBefore: *currentUserPoints,
			RoutePoints:      *routeRecord.Points,
			EventId:          event.ID,
			EventDate:        event.Start.Unix(),
		}); err != nil {
			worker.logger.Warnf(
				"Failed to save current points (%d) for user %d: %+v",
				*currentUserPoints, userId, err,
			)
		}
	}
	return true, nil
}

func processEvent(
	context *WorkerProcessContext,
	event *hb.Event,
) (bool, error) {
	worker := context.Worker
	if event.Activity != hb.HikingActivity {
		worker.logger.Infof("ignore non-hiking activity %s with ID %d", event.Title, event.ID)
		return false, nil
	}

	now := time.Now()
	if now.Sub(event.Start).Hours() > hb.AssignPointsForEventHourThreshold {
		worker.logger.Infof("ignore activity %s since its points has been assigned", event.Title)
		return false, nil
	}

	return processRecentPastEvent(context, event)
}

func pastEventProcessFunc(context *WorkerProcessContext) error {
	worker := context.Worker
	fetchResults, err := hb.FetchPastEvents(context.Credential)
	if err != nil {
		worker.logger.Warnf("Unable to fetch events: %+v\n", err)
		return err
	}

	for i, event := range fetchResults.Results {
		if !worker.shouldRun {
			worker.logger.Infof("Give up processing remaining %d events to stop the worker", len(fetchResults.Results)-(i+1))
			return nil
		}

		if _, err := processEvent(context, &event); err != nil {
			worker.logger.Warnf("failed to process event with id %d: %+v\n", event.ID, err)
		}
	}

	return nil
}
