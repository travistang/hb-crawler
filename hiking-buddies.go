package main

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

type User struct {
	Name, LastName string
	Experience     ExperienceLevel
	Id             uint16
}

type Route struct {
	ElevationGain uint16   `json:"elevation_gain"`
	Id            uint16   `json:"id"`
	Distance      uint16   `json:"distance"`
	SummitHeight  uint16   `json:"summit_height"`
	Title         string   `json:"title"`
	Scale         SACScale `json:"sac_scale"`
}

type Event struct {
	Route        *Route
	Organizer    *User
	Title        string
	Id           uint16
	ActivityType string
}

type URL string

const (
	MainDomain        URL = "hiking-buddies.com"
	EventListEndpoint URL = "https://www.hiking-buddies.com/api/routes/upcoming_event_list/"
)
