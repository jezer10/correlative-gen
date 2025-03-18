package models

type InsertTaskPayload struct {
	DSN        string
	HttpMethod string
	Data       User
}
