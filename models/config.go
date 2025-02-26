package models

type DBConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DbType   string `json:"dbtype"`
}

type Config struct {
	Databases []DBConfig `json:"databases"`
}
