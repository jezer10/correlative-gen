// config/config.go
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/emidiaz3/event-driven-server/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabaseConnection() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	db.AutoMigrate(&models.Postulantes{})
	triggerSQL := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE OR REPLACE FUNCTION generate_correlative()
		RETURNS TRIGGER AS $$
		DECLARE
			new_id TEXT;
		BEGIN
			new_id := '1' || LPAD(NEW.id::TEXT, 8, '0');
			NEW.correlativo := new_id;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE OR REPLACE TRIGGER trigger_generate_correlative
		BEFORE INSERT ON postulantes
		FOR EACH ROW
		EXECUTE FUNCTION generate_correlative();
	`

	if err := db.Exec(triggerSQL).Error; err != nil {
		fmt.Println("Error al crear el trigger:", err)
	} else {
		fmt.Println("Trigger creado correctamente.")
	}

	db.AutoMigrate(&models.Users{})

	fmt.Println("Migración completada")
	DB = db
}
