package models

type User struct {
	Id            string
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Identity      string `json:"identity"`
	Birthday      string `json:"birthday"`
	NativeCountry string `json:"native_country"`
	Country       string `json:"country"`
	Correlative   string
}
