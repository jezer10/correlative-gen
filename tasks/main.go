package tasks

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/emidiaz3/event-driven-server/models"
	"github.com/emidiaz3/event-driven-server/utils"

	"github.com/hibiken/asynq"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	TypeInsertUserMysql      = "insert:user:mysql"
	TypeInsertUserSqlServer  = "insert:user:sqlserver"
	TypeInsertUserWebService = "insert:user:api"
)

func SendUserToQueue(redisClient *asynq.Client, user models.User) error {

	for _, dbConfig := range utils.Config.Databases {
		payload, _ := json.Marshal(map[string]any{
			"DSN":        dbConfig.DSN,
			"HttpMethod": dbConfig.HttpMethod,
			"Data":       user,
		})
		taskType := fmt.Sprintf("insert:user:%s", dbConfig.DbType)
		fmt.Println(taskType)
		task := asynq.NewTask(taskType, payload)

		opts := []asynq.Option{
			asynq.MaxRetry(5),
			asynq.Timeout(30 * time.Second),
			asynq.Queue("critical"),
			asynq.Retention(24 * time.Hour),
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
	if err != nil {
		return fmt.Errorf("error ejecutando query: %w", err)
	}

	return nil
}

func HandleInsertUserWebService(ctx context.Context, t *asynq.Task) error {
	var payload models.InsertTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("error decodificando JSON: %w", err)
	}
	jsonData, err := json.Marshal(payload.Data)
	if err != nil {
		return fmt.Errorf("error al serializar la data: %w", err)
	}
	fmt.Println(payload.DSN)
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(payload.HttpMethod), payload.DSN, bytes.NewBuffer(jsonData))

	if err != nil {
		return fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("error al ejecutar request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {

		return fmt.Errorf("status code inesperado: %d", resp.StatusCode)

	}

	// TODO: USAR PAYLOAD PARA REALIZAR METODO A WEB SERVICE -KEV

	return nil
}

func LoggingMiddleware(next asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		fmt.Printf("Procesando tarea: %s\n", t.Type())
		return next.ProcessTask(ctx, t)
	})
}

// logging middlware for error
func ErrorLoggingMiddleware(next asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		err := next.ProcessTask(ctx, t)
		if err != nil {
			log.Printf("Error procesando tarea: %v", err)
		}
		return err
	})
}
