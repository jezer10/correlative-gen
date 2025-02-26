package models

type DatabaseItem struct {
	DSN         string `json:"dsn"`
	DbType      string `json:"type"`
	InsertQuery string `json:"insertQuery"`
}

type Config struct {
	Databases []DatabaseItem `json:"databases"`
}
