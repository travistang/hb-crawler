package hiking_buddies

import "hb-crawler/rating-gain/database"

type Route struct {
	RouteTitle    string  `json:"title"`
	Distance      float64 `json:"distance"`
	RouteID       int     `json:"id"`
	ElevationGain int     `json:"elevation_gain"`
	SacScale      string  `json:"sac_scale"`
}

func (r *Route) ToRouteRecord() *database.RouteRecord {
	return &database.RouteRecord{
		Id:        r.RouteID,
		Elevation: r.ElevationGain,
		Points:    nil,
		Distance:  float32(r.Distance),
		Name:      r.RouteTitle,
		Scale:     r.SacScale,
	}
}
