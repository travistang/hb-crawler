package api

import (
	"hb-crawler/rating-gain/database"
	"hb-crawler/rating-gain/worker"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type StartServerParams struct {
	Addr        string
	Repo        *database.DatabaseRepository
	WorkerGroup *worker.WorkerGroup
	WaitGroup   *sync.WaitGroup
}

func createApi(params *StartServerParams) *gin.Engine {
	api := gin.Default()
	api.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})

	pointGainsApi := PointGainsApiHandler{
		repo: params.Repo,
	}
	pointGainsApi.Register(api)

	workerApi := WorkerApiHandler{
		workerGroup: params.WorkerGroup,
	}
	workerApi.Register(api)

	return api
}

func StartServer(params *StartServerParams) *http.Server {
	api := createApi(params)
	server := &http.Server{
		Addr:    params.Addr,
		Handler: api,
	}
	params.WaitGroup.Add(1)
	go func() {
		log.Debugf("Starting API server...\n")
		server.ListenAndServe()
		params.WaitGroup.Done()
	}()

	return server

}
