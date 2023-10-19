package hiking_buddies

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

type RouteData struct {
	OrganizerID   int    `json:"organizer_id"`
	Activity      string `json:"activity"`
	EventTitle    string `json:"event_title"`
	FormattedDate string `json:"formatted_date"`
	Route         Route  `json:"title"`
	EventID       int    `json:"event_id"`
	Date          string `json:"date"`
}

type URL string

const (
	AssignPointsForEventHourThreshold = 72

	MainDomainWithoutProtocol URL = "www.hiking-buddies.com"
	MainDomain                URL = "https://www.hiking-buddies.com"
	LoginEndpoint             URL = "https://www.hiking-buddies.com/routes/login_user/"
	EventListEndpoint         URL = "https://www.hiking-buddies.com/api/routes/upcoming_event_list/"
	PastEventListEndpoint     URL = "https://www.hiking-buddies.com/api/routes/recent_event_list/"
	UserDetailsEndpoint       URL = "https://www.hiking-buddies.com/routes/user/"
	RouteDetailsEndpoint      URL = "https://www.hiking-buddies.com/routes/routes_list/"
	EventDetailsEndpoint      URL = "https://www.hiking-buddies.com/routes/events/"
)
