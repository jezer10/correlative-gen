package server

import (
	"fmt"
	"log"
	"os"

	"github.com/emidiaz3/event-driven-server/tasks"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
)

func StartAsynqServer() {
	redisServer := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeInsertUserMysql, tasks.HandleInsertUserTaskMySql)
	mux.HandleFunc(tasks.TypeInsertUserSqlServer, tasks.HandleInsertUserTaskSqlServer)
	fmt.Println("Iniciando Worker")
	if err := redisServer.Run(mux); err != nil {
		log.Fatalf("Error iniciando el worker: %v", err)
	}
}
