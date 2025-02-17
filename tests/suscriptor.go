package tests

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/emidiaz3/event-driven-server/models"

	"github.com/emidiaz3/event-driven-server/database"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
)

func InitSuscriptor() {

	err := database.InitDBMSSQL()
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "TuContraseñaSegura123!",
	})
	pubsub := rdb.PSubscribe(ctx, "__keyevent@0__:set")
	defer pubsub.Close()

	fmt.Println("Suscriptor de Redis iniciado...")

	for msg := range pubsub.Channel() {
		key := msg.Payload

		data, err := rdb.Get(ctx, key).Result()
		if err != nil {
			fmt.Printf("Error obteniendo valor de Redis para la key %s: %v\n", key, err)
			continue
		}
		fmt.Print(key)
		var person models.User
		if err := json.Unmarshal([]byte(data), &person); err != nil {
			fmt.Printf("Error unmarshal JSON para la key %s: %v\n", key, err)
			continue
		}

		// Insertamos en las 3 DB
		if err := database.InsertUserMSSQL(person); err != nil {
			fmt.Printf("Error insertando Person en DB para la key %s: %v\n", key, err)
		} else {
			fmt.Printf("Se insertó la persona con key %s en las 3 DB correctamente\n", key)
		}
	}

}
