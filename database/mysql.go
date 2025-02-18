package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/emidiaz3/event-driven-server/models"
	_ "github.com/go-sql-driver/mysql"
)

var DbMySQL *sql.DB

func InitMysqlDb() {
	var err error
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST_NAME")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")

	// Cadena de conexión
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	DbMySQL, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Error al conectar a la base de datos:", err)
	}
	if err = DbMySQL.Ping(); err != nil {
		log.Fatal("❌ No se pudo conectar a la base de datos:", err)
	}
	log.Println("✅ Conexión a la base de datos MYSQL exitosa")

	initSQL := `
	CREATE TABLE IF NOT EXISTS USUARIOS (
		ID INT AUTO_INCREMENT PRIMARY KEY,
		DNI VARCHAR(20) NULL,
		NOMBRE VARCHAR(80) NULL,
		APELLIDO VARCHAR(100) NULL,
		FECHANACIMIENTO DATE NULL,
		FECHAREGISTRO DATETIME NULL,
		CLIENTE CHAR(10) NULL,
		FLAG BOOLEAN NULL,
		NACIONALIDAD VARCHAR(40) NULL,
		STATUS CHAR(25) NULL,
		STATUSDESCRIPTION VARCHAR(40) NULL,
		STATUSNOTE VARCHAR(800) NULL,
		CORRELATIVO VARCHAR(50) NULL,
		SCORE VARCHAR(50) NULL,
		SCORE_DESCRIPTION TEXT NULL,
		SCORE_NOTE TEXT NULL,
		RESIDENCIA VARCHAR(20) NULL,
		FECHARESPUESTA DATETIME NULL
	);
	`
	_, err = DbMySQL.Exec(initSQL)
	if err != nil {
		log.Fatal("Error ejecutando comandos de inicialización:", err)
	}
	fmt.Println("✅ Base de datos MYSQL inicializada con tablas y triggers")

}

func InsertUserMysql(p models.User) error {
	var err error
	_, err = DbMySQL.Exec(`
	    INSERT INTO USUARIOS (NOMBRE, APELLIDO, DNI, FECHANACIMIENTO, NACIONALIDAD, RESIDENCIA, CORRELATIVO)
	    VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.FirstName, p.LastName, p.Identity, p.Birthday, p.NativeCountry, p.Country, p.Correlative,
	)
	if err != nil {
		return fmt.Errorf("error insertando en MYSQL: %v", err)
	}
	return nil
}

func GetUserByIdMysql(userID int) (*models.User, error) {
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
		WHERE CORRELATIVO = ?
	`
	err = DbMySQL.QueryRow(
		query, userID,
	).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Identity, &user.Birthday, &user.NativeCountry, &user.Country, &user.Correlative)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func GetUserByCorrelativeMysql(Correlative int) (*models.User, error) {
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
		WHERE CORRELATIVO = ?
	`
	err = DbMySQL.QueryRow(
		query, Correlative,
	).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Identity, &user.Birthday, &user.NativeCountry, &user.Country, &user.Correlative)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
