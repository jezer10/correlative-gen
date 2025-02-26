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
	TypeInsertUser = "insert:user"
)

func SendUserToQueue(redisClient *asynq.Client, person models.User) error {

	args := []any{person.FirstName, person.LastName, person.Identity, person.Birthday, person.NativeCountry, person.Country, person.Correlative}
	for _, dbConfig := range utils.Config.Databases {
		payload, _ := json.Marshal(map[string]any{
			"DSN":    dbConfig.DSN,
			"Query":  dbConfig.InsertQuery,
			"DBType": dbConfig.DbType,
			"Args":   args,
		})
		task := asynq.NewTask(TypeInsertUser, payload)

		opts := []asynq.Option{
			asynq.MaxRetry(5),
			asynq.Timeout(30 * time.Second),
			asynq.Queue("critical"),
		}
		_, err := redisClient.Enqueue(task, opts...)
		if err != nil {
			return fmt.Errorf("error enviando tarea a Asynq: %w", err)
		}
		log.Printf("Usuario %s %s encolado correctamente a %s", person.FirstName, person.LastName, dbConfig.DbType)
	}

	return nil

}

func HandleInsertUserTask(ctx context.Context, t *asynq.Task) error {
	var payload models.InsertTaskPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("error decodificando JSON: %w", err)
	}
	var driver string
	switch payload.DBType {
	case "sqlserver":
		driver = "sqlserver"
	case "postgres":
		driver = "postgres"
	case "mysql":
		driver = "mysql"
	default:
		return fmt.Errorf("tipo de BD desconocido: %s", payload.DBType)
	}
	fmt.Println(payload.DSN)

	db, err := sql.Open(driver, payload.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println("Payload: ", payload.Query, payload.DBType)
	_, err = db.Exec(payload.Query, payload.Args...)
	fmt.Println(err)

	return err
	// if err := database.InsertUserMysql(user); err != nil {
	// 	return fmt.Errorf("error Insertando en BD Final: %w", err)

	// }

	return nil
}
