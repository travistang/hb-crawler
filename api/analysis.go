package api

import (
	"hb-crawler/rating-gain/analysis"
	"hb-crawler/rating-gain/database"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AnalysisEndpointRoot = "/analysis"
	EstimateEndpoint     = "/estimate"
	OptimizeEndpoint     = "/optimize"
)

type AnalysisApiHandler struct {
	repo      *database.DatabaseRepository
	estimator *analysis.PointGainEstimator
}

func (handler *AnalysisApiHandler) estimateHandler(c *gin.Context) {

	var queryPointGain database.ReducedPointGainRecord
	if err := c.Bind(&queryPointGain); err != nil {
		reportError(c, http.StatusBadRequest, "invalid request data")
		return
	}

	estimatedPoints := handler.estimator.EstimatePointGain(
		int32(queryPointGain.UserPointsBefore),
		int32(queryPointGain.RoutePoints),
	)
	sendJSONPayload(c, http.StatusOK, gin.H{
		"pointsAfter": math.Round(estimatedPoints),
	})
}

func (handler *AnalysisApiHandler) optimizeHandler(c *gin.Context) {
	records, err := handler.repo.PointGains.GetValidPointsGainEntry(nil)
	if err != nil {
		reportError(c, http.StatusInternalServerError, "failed to get records to optimize parameters from")
		return
	}
	optimizedEstimator, err := analysis.OptimizeEstimator(*records, handler.estimator)
	if err != nil {
		reportError(c, http.StatusInternalServerError, "failed to optimize estimator")
		return
	}
	newOptimizer, err := analysis.CreatePointGainEstimator(optimizedEstimator.Params)
	if err != nil {
		reportError(c, http.StatusInternalServerError, "failed to optimize estimator")
		return
	}
	handler.estimator = newOptimizer
	sendJSONPayload(c, http.StatusOK, *optimizedEstimator)
}

func (handler *AnalysisApiHandler) Register(api *gin.Engine) {
	router := api.Group(AnalysisEndpointRoot)
	router.POST(EstimateEndpoint, handler.estimateHandler)
	router.POST(OptimizeEndpoint, handler.optimizeHandler)
}
