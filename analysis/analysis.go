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
	params         []float64
	exponent       float64
	exponentFactor float64
	base           float64
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

func polynomial(x float64, params []float64) float64 {
	sum := float64(0)
	for i, p := range params {
		exp := len(params) - i - 1
		sum += p * math.Pow(x, float64(exp))
	}
	return sum
}

/*
*

	Given original point P_0, route rating R, point P_f after joining the hike is modelled as:

	P_f = P_0 + k(1 - (1 + 10^d/e)^-1)

	where
		d = |P_0 - R|,
		k = Ad^B + poly(d, k1, k2, ...) = k1d^n + k2d^(n - 1) + ... + k0
	in which A, B, e, k1, k2... are parameters:
		A <- e.exponentFactor
		B <- e.exponent
		e <- e.base
		k1, k2, ... kn <- e.params

*
*/
func (e *PointGainEstimator) EstimatePointGain(originalPoint int32, routeRating int32) float64 {
	d := math.Abs(float64(originalPoint) - float64(routeRating))

	k := polynomial(d, e.params)
	k += e.exponentFactor * math.Pow(d, e.exponent)

	r := float64(originalPoint) + k*e.rateFactor(originalPoint, routeRating)
	return r
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
