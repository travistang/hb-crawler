package hiking_buddies

import "time"

const (
	HikingActivity = "HI"
)

type Event struct {
	Organizer       Organizer `json:"organizer"`
	Route           Route     `json:"route"`
	CoverPictureURL string    `json:"cover_picture_url"`
	DurationInDays  int       `json:"num_of_days"`
	Start           time.Time `json:"start"`
	ParticipantsId  []int     `json:"participants"`
	Activity        string    `json:"activity"`
	ID              int       `json:"id"`
	RouteData       RouteData `json:"route_data"`
	Title           string    `json:"title"`
}
