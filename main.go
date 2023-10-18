package main

import (
	"hb-crawler/rating-gain/database"
	hb "hb-crawler/rating-gain/hiking-buddies"
	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	Username = "opulent_umpires0w@icloud.com"
	Password = "fygveq-5ruqJa-gusgap"
)

func registerInterrupt(handlers func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			handlers()
		}
	}()
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debugf("Starting to crawl hiking buddies...\n")
	log.Info("Initializing database...\n")

	db, err := database.InitializeDatabase("./db.sqlite")
	if err != nil {
		log.Fatalf("Failed to initialize databse: %+v\n", err)
	}

	repo := database.GetRepository(db)
	workers := sync.WaitGroup{}

	pastEventWorker := CreatePastEventWorker(&PastEventWorkerConfig{
		Repository: repo,
		Interval:   time.Hour,
		Credential: &hb.Credential{
			Email:    Username,
			Password: Password,
		},
	})

	registerInterrupt(func() {
		log.Info("SIGINT received, stopping worker...")
		pastEventWorker.Stop()
	})

	pastEventWorker.StartProcessing(&workers)
	log.Info("Launched")
	workers.Wait()
	log.Info("Exiting...")

}
