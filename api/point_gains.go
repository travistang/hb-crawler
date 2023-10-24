package api

import (
	"hb-crawler/rating-gain/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	PointGainsEndpointRoot = "/point-gains"
	PointGainsList         = "/"
	PointGainsOfEvent      = "/:id"
)

type PointGainsApiHandler struct {
	repo *database.DatabaseRepository
}

func getPointsGainQueryParams(c *gin.Context) *database.PointGainsQuery {
	params := database.PointGainsQuery{
		Limit: -1,
		Skip:  0,
	}

	if limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		params.Limit = limit
	}

	if skip, err := strconv.Atoi(c.Query("skip")); err == nil {
		params.Skip = skip
	}

	return &params
}

func (handler *PointGainsApiHandler) pointGainsListHandler(c *gin.Context) {
	params := getPointsGainQueryParams(c)
	records, err := handler.repo.PointGains.GetAllPointGains(params)

	if err != nil {
		reportError(c, http.StatusInternalServerError, "failed to retrieve data")
		return
	}

	sendJSONPayload(c, http.StatusOK, records)
}

func (handler *PointGainsApiHandler) pointGainsOfEventHandler(c *gin.Context) {
	paramIdString := c.Params.ByName("id")
	id, err := strconv.Atoi(paramIdString)
	if err != nil {
		reportError(c, http.StatusBadRequest, "id must be integer")
		return
	}
	records, err := handler.repo.PointGains.GetPointGainsByEventId(id)
	if err != nil {
		reportError(c, http.StatusInternalServerError, "failed to retrieve data")
		return
	}
	sendJSONPayload(c, http.StatusOK, records)
}

func (handler *PointGainsApiHandler) Register(api *gin.Engine) {
	router := api.Group(PointGainsEndpointRoot)
	router.GET(PointGainsList, handler.pointGainsListHandler)
	router.GET(PointGainsOfEvent, handler.pointGainsOfEventHandler)
}
