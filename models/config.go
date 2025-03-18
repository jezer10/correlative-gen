package models

type DatabaseItem struct {
	DSN         string `json:"dsn"`
	DbType      string `json:"type"`
	HttpMethod  string `json:"method"`
	InsertQuery string `json:"insertQuery"`
}

type Config struct {
	Databases []DatabaseItem `json:"databases"`
}
