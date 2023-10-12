package main

import (
	"fmt"
	database "hb-crawler/rating-gain/database"
	hiking_buddies "hb-crawler/rating-gain/hiking-buddies"
	"log"
)

func main() {
	db, err := database.InitializeDatabase("./db.sqlite")

	if err != nil {
		log.Fatal("Unable to open database: %+v\n", err)
	}

	userRepository := database.CreateUserRepository(db)

	id, err := userRepository.CreateUser(&hiking_buddies.User{
		ID:       123,
		Name:     "Travis",
		LastName: "Tang",
	})

	if err != nil {
		log.Fatal("Unable to create user")
	}

	fmt.Printf("New user created with id %d\n", id)
}
