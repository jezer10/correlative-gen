package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/emidiaz3/event-driven-server/models"
	"github.com/emidiaz3/event-driven-server/utils"

	"github.com/hibiken/asynq"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	TypeInsertUserMysql     = "insert:user:mysql"
	TypeInsertUserSqlServer = "insert:user:sqlserver"
)

func SendUserToQueue(redisClient *asynq.Client, user models.User) error {

	for _, dbConfig := range utils.Config.Databases {
		payload, _ := json.Marshal(map[string]any{
			"DSN":  dbConfig.DSN,
			"Data": user,
		})
		taskType := fmt.Sprintf("insert:user:%s", dbConfig.DbType)
		task := asynq.NewTask(taskType, payload)

		opts := []asynq.Option{
			asynq.MaxRetry(5),
			asynq.Timeout(30 * time.Second),
			asynq.Queue("critical"),
		}
		_, err := redisClient.Enqueue(task, opts...)
		if err != nil {
			return fmt.Errorf("error enviando tarea a Asynq: %w", err)
		}
		log.Printf("Usuario %s %s encolado correctamente a %s", user.FirstName, user.LastName, dbConfig.DbType)
	}

	return nil

}

func HandleInsertUserTaskSqlServer(ctx context.Context, t *asynq.Task) error {
	var payload models.InsertTaskPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("error decodificando JSON: %w", err)
	}
	db, err := sql.Open("sqlserver", payload.DSN)
	fmt.Println(err)

	if err != nil {
		return err
	}
	defer db.Close()
	query := `
	INSERT
	INTO 
	USUARIOS 
		(NOMBRE, APELLIDO, DNI, FECHANACIMIENTO, NACIONALIDAD, RESIDENCIA, CORRELATIVO, FECHAREGISTRO) 
	VALUES 
		(@Nombre, @Apellido, @Dni, @FechaNacimiento, @Nacionalidad, @Residencia, @Correlativo, GETDATE())
	`
	_, err = db.Exec(query,
		sql.Named("Nombre", payload.Data.FirstName),
		sql.Named("Apellido", payload.Data.LastName),
		sql.Named("Dni", payload.Data.Identity),
		sql.Named("FechaNacimiento", payload.Data.Birthday),
		sql.Named("Nacionalidad", payload.Data.NativeCountry),
		sql.Named("Residencia", payload.Data.Country),
		sql.Named("Correlativo", payload.Data.Correlative),
	)
	fmt.Println(err)

	if err != nil {
		return fmt.Errorf("error ejecutando query: %w", err)
	}

	return nil
}

func HandleInsertUserTaskMySql(ctx context.Context, t *asynq.Task) error {
	var payload models.InsertTaskPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("error decodificando JSON: %w", err)
	}
	db, err := sql.Open("mysql", payload.DSN)
	fmt.Println(err)
	if err != nil {
		return err
	}
	defer db.Close()
	query := `
	INSERT
	INTO 
	USUARIOS 
		(NOMBRE, APELLIDO, DNI, FECHANACIMIENTO, NACIONALIDAD, RESIDENCIA, CORRELATIVO, FECHAREGISTRO) 
	VALUES 
		(? , ? , ? , ? , ? , ? , ?, NOW())
	`
	_, err = db.Exec(query, payload.Data.FirstName, payload.Data.LastName, payload.Data.Identity, payload.Data.Birthday, payload.Data.NativeCountry, payload.Data.Country, payload.Data.Correlative)
	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("error ejecutando query: %w", err)
	}

	return nil
}
