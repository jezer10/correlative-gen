package database

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/emidiaz3/event-driven-server/models"
)

var (
	DbSQLServer *sql.DB
)

func InitDBMSSQL() error {
	var err error
	DbSQLServer, err = sql.Open("sqlserver", "Server=localhost;Database=test;User Id=sa;Password=secret;TrustServerCertificate=true")
	if err != nil {
		log.Fatal("❌ Error al conectar a la base de datos:", err)
	}
	if err = DbSQLServer.Ping(); err != nil {
		log.Fatal("❌ No se pudo conectar a la base de datos:", err)
	}
	log.Println("✅ Conexión a la base de datos MSSQL exitosa")

	return nil
}

func InsertUserMSSQL(p models.User) error {
	var err error

	_, err = DbSQLServer.Exec(`
	    INSERT INTO persons (first_name, last_name, identity, birthday, native_country, country)
	    VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`,
		p.FirstName, p.LastName, p.Identity, p.Birthday, p.NativeCountry, p.Country,
	)
	if err != nil {
		log.Fatal("error insertando en MSSQL", err)
	}

	return nil
}
