package analysis

import (
	"fmt"
	"hb-crawler/rating-gain/database"
	"math"

	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/optimize"
)

const (
	InitialM float64 = 0.264
	InitialL float64 = 55
)

type PointGainEstimator struct {
	params []float64
	base   float64
}

type PointGainEstimatorGrad = PointGainEstimator // gradient of the params of the estimator

func CreatePointGainEstimator(params []float64) (*PointGainEstimator, error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("not enough params")
	}
	return &PointGainEstimator{
		params: params[:len(params)-1],
		base:   params[len(params)-1],
	}, nil
}

func (e *PointGainEstimator) Params() []float64 {
	return append(e.params, e.base)
}

func (e *PointGainEstimator) Base() float64 {
	return e.base
}

func (e *PointGainEstimator) rateFactor(originalPoint int32, routeRating int32) float64 {
	return 1 - 1/(1+math.Pow(10, float64(routeRating-originalPoint)/e.base))
}

func (e *PointGainEstimator) EstimatePointGain(originalPoint int32, routeRating int32) float64 {
	d := math.Abs(float64(originalPoint) - float64(routeRating))
	k := float64(0)
	for i, p := range e.params {
		k += p * math.Pow(d, float64(len(e.params)-1-i))
	}
	r := float64(originalPoint) + k*e.rateFactor(originalPoint, routeRating)
	return r
}

func (e *PointGainEstimator) GradientAt(pointGain *database.ReducedPointGainRecord) (*PointGainEstimatorGrad, error) {
	x0 := int32(pointGain.UserPointsBefore)
	xh := int32(pointGain.RoutePoints)

	rateFactor := e.rateFactor(x0, xh)
	diff := float64(pointGain.UserPointsAfter) - float64(e.EstimatePointGain(x0, xh))
	absDiff := math.Abs(float64(xh - x0))
	grads := []float64{}
	for i := range e.params {
		grads = append(grads, diff*rateFactor*math.Pow(absDiff, float64(len(e.params)-1-i)))
	}

	return &PointGainEstimatorGrad{
		params: grads,
	}, nil
}

func createOptimizerProblem(pointGains []database.ReducedPointGainRecord) *optimize.Problem {
	lossFunc := func(x []float64) float64 {
		estimator, _ := CreatePointGainEstimator(x)
		loss := float64(0)
		for _, record := range pointGains {
			estimated := estimator.EstimatePointGain(int32(record.UserPointsBefore), int32(record.RoutePoints))
			diff := float64(estimated) - float64(record.UserPointsAfter)
			loss += math.Abs(diff)
		}
		loss = loss / float64(len(pointGains))
		return loss
	}

	return &optimize.Problem{
		Func: lossFunc,
	}
}

type OptimizeResult struct {
	InitialParams []float64 `json:"initial_params"`
	InitialLoss   float64   `json:"initial_loss"`
	Params        []float64 `json:"params"`
	Loss          float64   `json:"loss"`
}

func OptimizeEstimator(
	pointGains []database.ReducedPointGainRecord,
	initialEstimator *PointGainEstimator,
) (*OptimizeResult, error) {
	problem := createOptimizerProblem(pointGains)
	initialParams := initialEstimator.Params()
	initialLoss := problem.Func(initialParams)
	result, err := optimize.Minimize(
		*problem, initialParams,
		&optimize.Settings{Concurrent: 4}, nil,
	)
	if err != nil {
		logrus.Errorf("failed to optimize estimator: %+v\n", err)
		return nil, err
	}
	return &OptimizeResult{
		InitialParams: initialParams,
		Params:        result.X,
		InitialLoss:   initialLoss,
		Loss:          result.F,
	}, nil
}
