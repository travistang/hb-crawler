package hiking_buddies

import "time"

type ExperienceLevel string

const (
	beginner          ExperienceLevel = "Beginner"
	advanced_beginner ExperienceLevel = "Advanced Beginner"
	intermediate      ExperienceLevel = "Intermediate"
	advanced          ExperienceLevel = "Advanced"
	mountain_goat     ExperienceLevel = "Mountain Goat"
)

type SACScale string

const (
	t1 SACScale = "T1"
	t2 SACScale = "T2"
	t3 SACScale = "T3"
	t4 SACScale = "T4"
	t5 SACScale = "T5"
	t6 SACScale = "T6"
)

type Organizer struct {
	ExperienceLevelCategory ExperienceLevel `json:"experience_level_category"`
	Name                    string          `json:"name"`
	ID                      int             `json:"id"`
	LastName                string          `json:"last_name"`
	Picture                 string          `json:"picture"`
}

type User = Organizer

type Route struct {
	RouteTitle    string  `json:"route_title"`
	Distance      float64 `json:"distance"`
	RouteID       int     `json:"route_id"`
	ElevationGain int     `json:"elevation_gain"`
	SacScale      string  `json:"sac_scale"`
}

type RouteData struct {
	OrganizerID   int    `json:"organizer_id"`
	Activity      string `json:"activity"`
	EventTitle    string `json:"event_title"`
	FormattedDate string `json:"formatted_date"`
	Route         Route  `json:"route"`
	EventID       int    `json:"event_id"`
	Date          string `json:"date"`
}

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

type URL string

const (
	MainDomain        URL = "https://www.hiking-buddies.com"
	LoginEndpoint     URL = "https://www.hiking-buddies.com/routes/login_user/"
	EventListEndpoint URL = "https://www.hiking-buddies.com/api/routes/upcoming_event_list/"
)
