package models

type CountryData struct {
	ID          int    `json:"id"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
}

type CountryDataProfile struct {
	ID          int    `json:"id"`
	CountryName string `json:"countryName"`
}
