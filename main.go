package main

import (
	"sync"

	"github.com/emidiaz3/event-driven-server/database"
	"github.com/emidiaz3/event-driven-server/server"
	"github.com/emidiaz3/event-driven-server/utils"
)

func main() {
	utils.InitEnv()
	utils.InitApiKey()
	database.InitDb()
	defer database.DB.Close()
	database.InitMysqlDb()
	defer database.DbMySQL.Close()
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		server.StartFiberServer()
	}()

	go func() {
		defer wg.Done()
		server.StartAsynqServer()
	}()

	wg.Wait()
}
