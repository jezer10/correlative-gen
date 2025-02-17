package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/emidiaz3/event-driven-server/database"
	"github.com/emidiaz3/event-driven-server/models"
	"github.com/hibiken/asynq"
)

const (
	TypeInsertUser = "insert:user"
)

func SendUserToQueue(redisClient *asynq.Client, person models.User) error {
	payload, err := json.Marshal(person)
	if err != nil {
		return fmt.Errorf("error convirtiendo usuario a JSON: %w", err)
	}

	task := asynq.NewTask(TypeInsertUser, payload)

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.Timeout(30 * time.Second),
		asynq.Queue("critical"),
	}

	_, err = redisClient.Enqueue(task, opts...)
	if err != nil {
		return fmt.Errorf("error enviando tarea a Asynq: %w", err)
	}

	log.Printf("Usuario %s %s encolado correctamente", person.FirstName, person.LastName)
	return nil
}

func HandleInsertUserTask(ctx context.Context, t *asynq.Task) error {
	var user models.User

	if err := json.Unmarshal(t.Payload(), &user); err != nil {
		return fmt.Errorf("error decodificando JSON: %w", err)
	}

	if err := database.InsertUserMysql(user); err != nil {
		return fmt.Errorf("error Insertando en BD Final: %w", err)

	}

	return nil
}
