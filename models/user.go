package models

import "github.com/emidiaz3/event-driven-server/types"

type User struct {
	Id               int
	FirstName        string           `json:"first_name"`
	LastName         string           `json:"last_name"`
	Identity         string           `json:"identity"`
	Birthday         string           `json:"birthday"`
	NativeCountry    string           `json:"native_country"`
	Country          string           `json:"country"`
	Correlative      string           `json:"check_id"`
	GlobalStatus     types.NullString `json:"global_status"`
	Score            types.NullString `json:"score"`
	ScoreDescription types.NullString `json:"score_description"`
	ScoreNote        types.NullString `json:"score_note"`
}
