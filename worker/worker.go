package worker

import (
	"hb-crawler/rating-gain/database"
	hb "hb-crawler/rating-gain/hiking-buddies"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Worker struct {
	repository      *database.DatabaseRepository
	shouldRun       bool
	interval        time.Duration
	LastRunningTime *time.Time
	logger          *log.Logger
	credential      *hb.Credential
	ProcessFunc     WorkerProcessFunc
}

type WorkerConfig struct {
	Repository *database.DatabaseRepository
	Credential *hb.Credential
	Interval   time.Duration
}

type WorkerStatus struct {
	Running bool       `json:"running"`
	LastRun *time.Time `json:"last_run"`
}

type WorkerProcessContext struct {
	Worker     *Worker
	Credential *hb.CookieCredential
}

type WorkerProcessFunc = func(*WorkerProcessContext) error

type ProceedSignal int

const (
	ShouldStop ProceedSignal = iota
	ShouldIgnore
	ShouldProcess
)

func (w *Worker) Stop() {
	w.logger.Info("Stop processing...")
	w.shouldRun = false
}

func (w *Worker) getProceedSignal() ProceedSignal {
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

func (w *Worker) StartProcessing(wg *sync.WaitGroup) {
	if w.shouldRun {
		w.logger.Warn("Refuse to run an already running worker")
		return
	}

	w.logger.Info("Starting worker...")
	w.shouldRun = true
	wg.Add(1)

	go func() {
		for {
			signal := w.getProceedSignal()
			if signal == ShouldStop {
				w.logger.Info("Stopping worker...")
				wg.Done()
				w.LastRunningTime = nil
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
				continue
			}

			if err := w.ProcessFunc(&WorkerProcessContext{
				Worker:     w,
				Credential: credential,
			}); err != nil {
				w.logger.Warnf("Worker encountered error %+v\n", err)
				continue
			}

			w.logger.Info("Process completed. Waiting for next run.")
		}
		w.logger.Info("Process stopped.")
	}()
}

func (w *Worker) Status() *WorkerStatus {
	return &WorkerStatus{
		LastRun: w.LastRunningTime,
		Running: w.shouldRun,
	}
}
