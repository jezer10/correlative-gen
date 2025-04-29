package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/emidiaz3/event-driven-server/types"
)

type CreateUserDto struct {
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	Identity      string `json:"identity" validate:"required"`
	Birthday      string `json:"birthday" validate:"required"`
	NativeCountry string `json:"native_country" validate:"required"`
	Country       string `json:"country" validate:"required"`
	Correlative   string `json:"check_id" validate:"required"`
	GlobalStatus  string `json:"global_status" validate:"required"`
}

type User struct {
	Id               int              `gorm:"column:id"`
	FirstName        string           `gorm:"column:nombre" json:"first_name"`
	LastName         string           `gorm:"column:apellido" json:"last_name"`
	Identity         string           `gorm:"column:identidad" json:"identity"`
	Birthday         string           `gorm:"column:fechanacimiento" json:"birthday"`
	NativeCountry    string           `gorm:"column:nacionalidad" json:"native_country"`
	Country          string           `gorm:"column:pais" json:"country"`
	Correlative      string           `gorm:"column:correlativo" json:"check_id"`
	GlobalStatus     types.NullString `gorm:"column:status" json:"global_status"`
	Score            types.NullString `gorm:"column:score" json:"score"`
	ScoreDescription types.NullString `gorm:"column:descripcion_score" json:"score_description"`
	ScoreNote        types.NullString `gorm:"column:nota_score" json:"score_note"`
}

type Usuarios struct {
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
type Postulantes struct {
	Id                int            `gorm:"primaryKey;autoIncrement"`
	Client            string         `gorm:"column:cliente; type:char(10);" json:"client"`
	FirstName         string         `gorm:"column:nombre; type:varchar(80);" json:"first_name"`
	LastName          string         `gorm:"column:apellido; type:varchar(100);" json:"last_name"`
	Identity          string         `gorm:"column:dni; type:varchar(20)" json:"identity"`
	Birthday          string         `gorm:"column:fecha_nacimiento; char(10)" json:"birthday"`
	NativeCountry     string         `gorm:"column:nacionalidad; varchar(40)" json:"native_country"`
	Country           string         `gorm:"column:residencia; varchar(20)" json:"country"`
	Correlative       string         `gorm:"column:correlativo; varchar(50)" json:"check_id"`
	Status            sql.NullString `gorm:"column:status; type:char(25);" json:"global_status"`
	StatusDescription *string        `gorm:"column:status_description; type:char(25);" json:"global_status_description"`
	StatusNote        *string        `gorm:"column:status_note; type:char(25);" json:"global_status_note"`
	Score             *string        `gorm:"column:score; type:varchar(50);" json:"score"`
	ScoreDescription  *string        `gorm:"column:score_description; type:text;" json:"score_description"`
	ScoreNote         *string        `gorm:"column:score_note; type:text;" json:"score_note"`
	CreatedAt         time.Time      `gorm:"column:fecha_registro; type:timestamp; default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time      `gorm:"column:fecha_respuesta; type:timestamp; default:CURRENT_TIMESTAMP"`
	Flag              bool           `gorm:"column:flag; type:boolean;"`
}

type Users struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username  string    `gorm:"column:username;type:varchar(255);uniqueIndex"`
	Password  string    `gorm:"column:password;type:varchar(255)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Users) TableName() string {
	return "users"
}

func (User) TableName() string {
	return "usuarios"
}

func (Usuarios) TableName() string {
	return "usuarios"
}

func (Postulantes) TableName() string {
	return "postulantes"
}
