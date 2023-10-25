package worker

import (
	"time"

	log "github.com/sirupsen/logrus"

	"hb-crawler/rating-gain/database"
	hb "hb-crawler/rating-gain/hiking-buddies"
	"hb-crawler/rating-gain/logging"
)

const (
	TimeAfterEventToAssignPoints = 72 * time.Hour
)

func CreatePointsGainWorker(config *WorkerConfig) *Worker {
	logger := logging.GetLogger(&logging.LoggerConfig{
		Prefix: "Points gain worker",
		Level:  log.DebugLevel,
	})

	return &Worker{
		repository:      config.Repository,
		shouldRun:       false,
		interval:        config.Interval,
		logger:          logger,
		LastRunningTime: nil,
		ProcessFunc:     pointsGainProcessFunc,
	}
}

func processDanglingPointsGainEntry(
	context *WorkerProcessContext,
	pointsGain *database.PointGainRecord,
) error {
	worker := context.Worker
	worker.logger.Infof(
		"Fetching current points for user %d who participated events %d on %s",
		pointsGain.UserId, pointsGain.EventId,
		time.Unix(pointsGain.EventDate, 0).UTC(),
	)

	currentUserPoints, err := hb.FetchUserPoints(pointsGain.UserId, context.Credential)
	if err != nil {
		return err
	}
	worker.logger.Infof("User %d now has %d points", pointsGain.UserId, *currentUserPoints)
	if err := worker.repository.PointGains.UpdatePointsGainEntry(&database.PointGainRecord{
		UserId:          pointsGain.UserId,
		EventId:         pointsGain.EventId,
		UserPointsAfter: currentUserPoints,
	}); err != nil {
		return err
	}

	return nil
}

func pointsGainProcessFunc(context *WorkerProcessContext) error {
	worker := context.Worker
	worker.logger.Info("Find and updating dangling points gain entry...")

	targetHour := time.Now().Add(-TimeAfterEventToAssignPoints)
	worker.logger.Infof("Fetching dangling points gain record in point of %s", targetHour.UTC())
	danglingRecords, err := worker.repository.PointGains.GetDanglingPointsGainEntryToday(targetHour)
	if err != nil {
		worker.logger.Warnf("Failed to fetch dangling records: %+v\n", err)
	}

	worker.logger.Infof("Found %d dangling records at around %s", len(*danglingRecords), targetHour.UTC())
	for index, pointsGain := range *danglingRecords {
		if !worker.shouldRun {
			worker.logger.Infof("Give up processing remaining %d records to stop the worker", len(*danglingRecords)-(index+1))
			break
		}

		if err := processDanglingPointsGainEntry(context, &pointsGain); err != nil {
			worker.logger.Warnf("Failed processing dangling points gain entry: %+v\n", err)
		}
	}

	return nil
}
