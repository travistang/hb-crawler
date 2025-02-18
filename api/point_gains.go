package api

import (
	"bytes"
	"encoding/csv"
	"hb-crawler/rating-gain/database"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	PointGainsEndpointRoot = "/point-gains"
	PointGainsList         = "/"
	PointsGainsSampleData  = "/sample"
	PointGainsOfEvent      = "/event/:id"
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

func returnSampleAsCSV(c *gin.Context, pointGains *[]database.ReducedPointGainRecord) {
	buffer := new(bytes.Buffer)
	w := csv.NewWriter(buffer)
	if err := w.Write([]string{"route_points", "points_before", "points_after"}); err != nil {
		reportError(c, http.StatusInternalServerError, "Failed to generate CSV")
		return
	}
	w.Flush()

	for _, gain := range *pointGains {
		if err := w.Write([]string{
			strconv.Itoa(gain.RoutePoints), strconv.Itoa(gain.UserPointsBefore), strconv.Itoa(gain.UserPointsAfter),
		}); err != nil {
			reportError(c, http.StatusInternalServerError, "Failed to generate CSV")
		}
		w.Flush()
	}
	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=hb_point_gain_sample.csv")
	_, err := c.Writer.Write(buffer.Bytes())
	if err != nil {
		reportError(c, http.StatusInternalServerError, "Failed to generate CSV")
	}
}

func (handler *PointGainsApiHandler) pointGainsSampleDataHandler(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = -1
	}
	format := strings.ToLower(c.Query("format"))
	if len(format) == 0 {
		format = "json"
	}
	records, err := handler.repo.PointGains.GetValidPointsGainEntry(&limit)
	if err != nil {
		reportError(c, http.StatusInternalServerError, "Failed to retrieve data")
		return
	}
	logrus.Infof("Records: %d", len(*records))

	switch format {
	case "json":
		sendJSONPayload(c, http.StatusOK, records)
		return
	case "csv":
		returnSampleAsCSV(c, records)
		return
	default:
		reportError(c, http.StatusBadRequest, "Unknown format")
		return
	}
}

func (handler *PointGainsApiHandler) Register(api *gin.Engine) {
	router := api.Group(PointGainsEndpointRoot)
	router.GET(PointGainsList, handler.pointGainsListHandler)
	router.GET(PointsGainsSampleData, handler.pointGainsSampleDataHandler)
	router.GET(PointGainsOfEvent, handler.pointGainsOfEventHandler)
}
