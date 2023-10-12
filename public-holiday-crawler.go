package main

import (
	"github.com/go-resty/resty/v2"
)

type CountryCode string

const (
	US CountryCode = "US"
	DE CountryCode = "DE"
)

type Holiday struct {
	Date        string
	Name        string
	CountryCode CountryCode
	Fixed       bool
	Global      bool
	Counties    []string
}

func GetPublicHolidays() (error, *[]Holiday) {
	client := resty.New()

	var responseData []Holiday
	_, err := client.R().
		EnableTrace().
		SetResult(&responseData).
		Get("https://date.nager.at/api/v2/publicholidays/2020/US")
	return err, &responseData
}
