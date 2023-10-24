package main

import (
	"context"
	"hb-crawler/rating-gain/api"
	"hb-crawler/rating-gain/database"
	hb "hb-crawler/rating-gain/hiking-buddies"
	"hb-crawler/rating-gain/worker"
	"os"
	"os/signal"
	"sync"

	log "github.com/sirupsen/logrus"
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

func retrieveCredential() *hb.Credential {
	cred := hb.Credential{}
	cred.Email = os.Getenv("HB_USERNAME")
	cred.Password = os.Getenv("HB_PASSWORD")

	if len(cred.Email) == 0 || len(cred.Password) == 0 {
		log.Fatalf("Missing credential to login")
	}
	return &cred
}

func main() {
	log.SetLevel(log.DebugLevel)
	credential := retrieveCredential()

	log.Info("Initializing database...\n")
	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize databse: %+v\n", err)
	}
	defer db.Close()
	repo := database.GetRepository(db)

	waitGroup := sync.WaitGroup{}
	workerGroup := worker.CreateWorkerGroup(repo, credential, &waitGroup)
	server := api.StartServer(&api.StartServerParams{
		Addr:        ":8080",
		Repo:        repo,
		WaitGroup:   &waitGroup,
		WorkerGroup: workerGroup,
	})

	registerInterrupt(func() {
		log.Info("SIGINT received, stopping...")
		workerGroup.Stop()
		if err := server.Shutdown(context.Background()); err != nil {
			log.Warnf("failed to shutdown server: %+v\n", err)
		}
	})

	waitGroup.Wait()
}
