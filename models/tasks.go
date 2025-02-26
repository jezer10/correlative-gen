package models

type InsertTaskPayload struct {
	DSN    string
	Query  string
	DBType string
	Args   []any
}
