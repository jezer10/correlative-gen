package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/emidiaz3/event-driven-server/models"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDb() {
	var err error
	dns := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	DB, err = sql.Open("postgres", dns)
	if err != nil {
		log.Fatal("❌ Error al conectar a la base de datos:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("❌ No se pudo conectar a la base de datos:", err)
	}

	log.Println("✅ Conexión a la base de datos Postgres exitosa")
	initSQL := `
	CREATE TABLE IF NOT EXISTS USUARIOS (
		ID SERIAL PRIMARY KEY,
		DNI VARCHAR(20),
		NOMBRE VARCHAR(80),
		APELLIDO VARCHAR(100),
		FECHANACIMIENTO CHAR(10),
		FECHAREGISTRO VARCHAR(150),
		CLIENTE CHAR(10),
		FLAG BOOLEAN,
		NACIONALIDAD VARCHAR(40),
		STATUS CHAR(25),
		STATUSDESCRIPTION VARCHAR(40),
		STATUSNOTE VARCHAR(800),
		CORRELATIVO VARCHAR(50),
		SCORE VARCHAR(50),
		SCORE_DESCRIPTION TEXT,
		SCORE_NOTE TEXT,
		RESIDENCIA VARCHAR(20),
		FECHARESPUESTA TIMESTAMP
	);

	CREATE OR REPLACE FUNCTION generate_correlative()
	RETURNS TRIGGER AS $$
	DECLARE
		new_id TEXT;
	BEGIN
		new_id := '1' || LPAD(NEW.id::TEXT, 8, '0');
		NEW.CORRELATIVO := new_id;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE OR REPLACE TRIGGER trigger_generate_correlative
	BEFORE INSERT ON USUARIOS
	FOR EACH ROW
	EXECUTE FUNCTION generate_correlative();
	`
	_, err = DB.Exec(initSQL)
	if err != nil {
		log.Fatal("Error ejecutando comandos de inicialización:", err)
	}
	fmt.Println("✅ Base de datos inicializada con tablas y triggers")
}

func InsertUser(user models.User) (string, error) {
	var generatedCorrelative string

	query := `INSERT INTO USUARIOS (NOMBRE, APELLIDO, DNI, FECHANACIMIENTO, NACIONALIDAD, RESIDENCIA) VALUES ($1, $2, $3, $4, $5, $6) RETURNING CORRELATIVO`
	err := DB.QueryRow(query, user.FirstName, user.LastName, user.Identity, user.Birthday, user.NativeCountry, user.Country).Scan(&generatedCorrelative)
	if generatedCorrelative == "" {
		return "", fmt.Errorf("bad generated correlative")
	}
	return generatedCorrelative, err
}

func GetUserById(userID int) (*models.User, error) {
	var err error
	var user models.User
	query := `
		SELECT
			ID,
			NOMBRE, 
			APELLIDO, 
			DNI, 
			FECHANACIMIENTO, 
			NACIONALIDAD, 
			RESIDENCIA, 
			CORRELATIVO 
		FROM 
			USUARIOS 
		WHERE CORRELATIVO = $1
		`
	err = DB.QueryRow(
		query, userID,
	).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Identity, &user.Birthday, &user.NativeCountry, &user.Country, &user.Correlative)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
