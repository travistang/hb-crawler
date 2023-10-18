package main

import (
	"sync"
	"time"

	"hb-crawler/rating-gain/database"
	hb "hb-crawler/rating-gain/hiking-buddies"
	"hb-crawler/rating-gain/logging"

	log "github.com/sirupsen/logrus"
)

type PastEventWorker struct {
	repository      *database.DatabaseRepository
	shouldRun       bool
	interval        time.Duration
	LastRunningTime *time.Time
	logger          *log.Logger
	credential      *hb.Credential
}

type PastEventWorkerConfig struct {
	Repository *database.DatabaseRepository
	Credential *hb.Credential
	Interval   time.Duration
	Logger     *log.Logger
}

func CreatePastEventWorker(config *PastEventWorkerConfig) *PastEventWorker {
	logger := config.Logger
	if logger == nil {
		logger = logging.GetLogger(&logging.LoggerConfig{
			Prefix: "past event worker",
			Level:  log.DebugLevel,
		})
	}

	return &PastEventWorker{
		repository:      config.Repository,
		shouldRun:       false,
		interval:        config.Interval,
		credential:      config.Credential,
		logger:          logger,
		LastRunningTime: nil,
	}
}

func (w *PastEventWorker) processRecentPastEvent(
	event *hb.Event,
	credential *hb.CookieCredential,
) (bool, error) {
	participantIds := event.AllParticipantsId()
	route := event.Route.ToRouteRecord()
	w.logger.Debugf("Found route in event: %+v\n", route)
	routeRepo := w.repository.Route
	if _, err := routeRepo.SaveRoute(route); err != nil {
		w.logger.Errorf("Failed to save route: %+v\n", err)
		return false, nil
	}

	var routeRecord database.RouteRecord
	if err := hb.GetRoutePoints(&hb.GetRoutePointsParams{
		Repo:       routeRepo,
		Id:         route.Id,
		Record:     &routeRecord,
		Credential: credential,
	}); err != nil {
		return false, nil
	}

	w.logger.Debugf("Route record is now %+v\n", routeRecord)

	for _, userId := range participantIds {
		currentUserPoints, err := hb.FetchUserPoints(userId, credential)
		if err != nil {
			w.logger.Warnf("Unable to fetch current points for userId %d\n", userId)
			continue
		}

		if err := w.repository.PointGains.CreatePointsGainEntry(&database.PointGainRecord{
			EventId:          event.ID,
			RoutePoints:      *(routeRecord.Points),
			UserId:           userId,
			UserPointsBefore: *currentUserPoints,
		}); err != nil {
			w.logger.Warnf("Unable to save points gain entry for user %d in event %d", userId, event.ID)
		}
	}

	return true, nil
}

func (w *PastEventWorker) processPointsAssignedEvent(
	event *hb.Event,
	credential *hb.CookieCredential,
) (bool, error) {
	participantIds := event.AllParticipantsId()
	for _, userId := range participantIds {
		currentUserPoints, err := hb.FetchUserPoints(userId, credential)
		if err != nil {
			w.logger.Warnf("Unable to fetch current points for userId %d\n", userId)
			continue
		}

		if err := w.repository.PointGains.UpdatePointsGainEntry(&database.PointGainRecord{
			EventId:         event.ID,
			UserId:          userId,
			UserPointsAfter: currentUserPoints,
		}); err != nil {
			w.logger.Warnf("Unable to update points gain entry for user %d which currently has points %d", userId, *currentUserPoints)
		}
	}
	return true, nil
}

func (w *PastEventWorker) processEvent(
	event *hb.Event,
	credential *hb.CookieCredential,
) (bool, error) {
	if event.Activity != hb.HikingActivity {
		w.logger.Infof("ignore non-hiking activity %s with ID %d", event.Title, event.ID)
		return false, nil
	}

	now := time.Now()
	if now.Sub(event.Start).Hours() > hb.AssignPointsForEventHourThreshold {
		return w.processPointsAssignedEvent(event, credential)
	}

	return w.processRecentPastEvent(event, credential)
}

func (w *PastEventWorker) Stop() {
	w.logger.Info("Stop processing...")
	w.shouldRun = false
}

type ProceedSignal int

const (
	ShouldStop ProceedSignal = iota
	ShouldIgnore
	ShouldProcess
)

func (w *PastEventWorker) shouldProceed() ProceedSignal {
	if !w.shouldRun {
		return ShouldStop
	}

	if w.LastRunningTime == nil {
		return ShouldProcess
	}

	now := time.Now()
	timeSinceLastRun := now.Sub(*w.LastRunningTime)
	if timeSinceLastRun > w.interval {
		return ShouldProcess
	}
	return ShouldIgnore
}
func (w *PastEventWorker) StartProcessing(wg *sync.WaitGroup) {
	if w.shouldRun {
		w.logger.Warn("Refuse to run an already running worker")
		return
	}

	w.logger.Info("Starting worker...")
	w.shouldRun = true
	wg.Add(1)
	go func() {
		for {
			signal := w.shouldProceed()
			if signal == ShouldStop {
				w.logger.Info("Stopping worker...")
				wg.Done()
				break
			}
			if signal == ShouldIgnore {
				continue
			}

			now := time.Now()
			w.LastRunningTime = &now
			w.logger.Info("Start processing...")

			credential, err := hb.Login(w.repository.Login, w.credential)
			if err != nil {
				w.logger.Warnf("Unable to login as user %s", w.credential.Email)
			}

			fetchResults, err := hb.FetchPastEvents(credential)
			if err != nil {
				w.logger.Warnf("Unable to fetch events: %+v\n", err)
			}

			for _, event := range fetchResults.Results {
				if _, err := w.processEvent(&event, credential); err != nil {
					w.logger.Warnf("failed to process event with id %d: %+v\n", event.ID, err)
				}
			}
			w.logger.Info("Process finished. Waiting for next round.")

		}
		w.logger.Info("Process stopped.")
	}()
}
