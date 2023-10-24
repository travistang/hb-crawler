package api

import (
	"hb-crawler/rating-gain/worker"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	WorkerEndpoint       = "/worker"
	WorkerStatusEndpoint = "/status"
	WorkerStartEndpoint  = "/start"
	WorkerStopEndpoint   = "/stop"
)

type WorkerApiHandler struct {
	workerGroup *worker.WorkerGroup
}

func (handler *WorkerApiHandler) workerStatusHandler(c *gin.Context) {
	status := handler.workerGroup.GetAllWorkerStatus()
	c.JSON(http.StatusOK, status)
}

func (handler *WorkerApiHandler) workerStartHandler(c *gin.Context) {
	handler.workerGroup.Start()
	sendJSONPayload(c, http.StatusOK, handler.workerGroup.GetAllWorkerStatus())
}

func (handler *WorkerApiHandler) workerStopHandler(c *gin.Context) {
	handler.workerGroup.Stop()
	sendJSONPayload(c, http.StatusOK, handler.workerGroup.GetAllWorkerStatus())
}

func (handler *WorkerApiHandler) Register(api *gin.Engine) {
	router := api.Group(WorkerEndpoint)
	router.GET(WorkerStatusEndpoint, handler.workerStatusHandler)
	router.POST(WorkerStartEndpoint, handler.workerStartHandler)
	router.POST(WorkerStopEndpoint, handler.workerStopHandler)
}
