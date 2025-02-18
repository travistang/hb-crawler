package worker

import (
	"hb-crawler/rating-gain/database"
	hb "hb-crawler/rating-gain/hiking-buddies"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	PastEventWorkerID  = "past-event"
	PointsGainWorkerID = "points-gain"
)

type WorkerGroup struct {
	Repository *database.DatabaseRepository
	Credential *hb.Credential
	workers    map[string]*Worker
	waitGroup  *sync.WaitGroup
}

func (c *WorkerGroup) Stop() {
	for _, worker := range c.workers {
		worker.Stop()
	}
}

func (c *WorkerGroup) Start() {
	for _, worker := range c.workers {
		worker.StartProcessing(c.waitGroup)
	}
}

func (c *WorkerGroup) Wait() {
	c.waitGroup.Wait()
}

func CreateWorkerGroup(repo *database.DatabaseRepository, waitGroup *sync.WaitGroup) *WorkerGroup {
	pastEventWorker := CreatePastEventWorker(&WorkerConfig{
		Repository: repo,
		Interval:   12 * time.Hour,
	})

	pointsGainWorker := CreatePointsGainWorker(&WorkerConfig{
		Repository: repo,
		Interval:   time.Hour,
	})

	workers := map[string]*Worker{}
	workers[PastEventWorkerID] = pastEventWorker
	workers[PointsGainWorkerID] = pointsGainWorker

	workerGroup := WorkerGroup{
		Repository: repo,
		workers:    workers,
		waitGroup:  waitGroup,
	}

	return &workerGroup
}

func (c *WorkerGroup) GetAllWorkerStatus() map[string]*WorkerStatus {
	status := map[string]*WorkerStatus{}
	if c.workers == nil {
		logrus.Warnf("Unable to find worker")
		return status
	}
	for id, worker := range c.workers {
		status[id] = worker.Status()
	}

	return status
}
